/*
Copyright 2017 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreedto in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Copyright 2019 The Gaea Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mysql

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"

	"github.com/XiaoMi/Gaea/util/bucketpool"
	"github.com/XiaoMi/Gaea/util/sync2"
)

const (
	// connBufferSize is how much we buffer for reading and
	// writing. It is also how much we allocate for ephemeral buffers.
	connBufferSize = 128

	// MaxPacketSize is the maximum payload length of a packet(16MB)
	// the server supports.
	MaxPacketSize = (1 << 24) - 1
)

// Constants for how ephemeral buffers were used for reading / writing.
const (
	// ephemeralUnused means the ephemeral buffer is not in use at this
	// moment. This is the default value, and is checked so we don't
	// read or write a packet while one is already used.
	ephemeralUnused = iota

	// ephemeralWrite means we currently in process of writing from  currentEphemeralBuffer
	ephemeralWrite

	// ephemeralRead means we currently in process of reading into currentEphemeralBuffer
	ephemeralRead
)

// Conn is a connection between a client and a server, using the MySQL
// binary protocol. It is built on top of an existing net.Conn, that
// has already been established.
//
// Use Connect on the client side to create a connection.
// Use NewListener to create a server side and listen for connections.
type Conn struct {
	// conn is the underlying network connection.
	// Calling Close() on the Conn will close this connection.
	// If there are any ongoing reads or writes, they may get interrupted.
	conn net.Conn

	// ConnectionID is set:
	// - at Connect() time for clients, with the value returned by
	// the server.
	// - at accept time for the server.
	ConnectionID uint32

	// closed is set to true when Close() is called on the connection.
	closed sync2.AtomicBool

	// Packet encoding variables.
	bufferedReader *bufio.Reader
	bufferedWriter *bufio.Writer
	sequence       uint8

	// Keep track of how and of the buffer we allocated for an
	// ephemeral packet on the read and write sides.
	// These fields are used by:
	// - StartEphemeralPacket / writeEphemeralPacket methods for writes.
	// - ReadEphemeralPacket / RecycleReadPacket methods for reads.
	currentEphemeralPolicy int
	// currentEphemeralBuffer for tracking allocated temporary buffer for writes and reads respectively.
	// It can be allocated from bufPool or heap and should be recycled in the same manner.
	currentEphemeralBuffer *[]byte
}

// bufPool is used to allocate and free buffers in an efficient way.
var bufPool = bucketpool.New(connBufferSize, MaxPacketSize)

// writersPool is used for pooling bufio.Writer objects.
var writersPool = sync.Pool{New: func() interface{} { return bufio.NewWriterSize(nil, connBufferSize) }}

// NewConn is an internal method to create a Conn. Used by client and server
// side for common creation code.
func NewConn(conn net.Conn) *Conn {
	return &Conn{
		conn:           conn,
		closed:         sync2.NewAtomicBool(false),
		bufferedReader: bufio.NewReaderSize(conn, connBufferSize),
	}
}

// StartWriterBuffering starts using buffered writes. This should
// be terminated by a call to flush.
func (c *Conn) StartWriterBuffering() {
	c.bufferedWriter = writersPool.Get().(*bufio.Writer)
	c.bufferedWriter.Reset(c.conn)
}

// Flush flushes the written data to the socket.
// This must be called to terminate startBuffering.
func (c *Conn) Flush() error {
	if c.bufferedWriter == nil {
		return nil
	}

	defer func() {
		c.bufferedWriter.Reset(nil)
		writersPool.Put(c.bufferedWriter)
		c.bufferedWriter = nil
	}()

	return c.bufferedWriter.Flush()
}

// getWriter returns the current writer. It may be either
// the original connection or a wrapper.
func (c *Conn) getWriter() io.Writer {
	if c.bufferedWriter != nil {
		return c.bufferedWriter
	}
	return c.conn
}

// getReader returns reader for connection. It can be *bufio.Reader or net.Conn
// depending on which buffer size was passed to newServerConn.
func (c *Conn) getReader() io.Reader {
	if c.bufferedReader != nil {
		return c.bufferedReader
	}
	return c.conn
}

func (c *Conn) readHeaderFrom(r io.Reader) (int, error) {
	var header [4]byte
	// Note io.ReadFull will return two different types of errors:
	// 1. if the socket is already closed, and the go runtime knows it,
	//   then ReadFull will return an error (different than EOF),
	//   someting like 'read: connection reset by peer'.
	// 2. if the socket is not closed while we start the read,
	//   but gets closed after the read is started, we'll get io.EOF.
	if _, err := io.ReadFull(r, header[:]); err != nil {
		// The special casing of propagating io.EOF up
		// is used by the server side only, to suppress an error
		// message if a client just disconnects.
		if err == io.EOF {
			return 0, err
		}
		if strings.HasSuffix(err.Error(), "read: connection reset by peer") {
			return 0, io.EOF
		}
		return 0, fmt.Errorf("io.ReadFull(header size) failed: %v", err)
	}

	sequence := uint8(header[3])
	if sequence != c.sequence {
		return 0, fmt.Errorf("invalid sequence, expected %v got %v", c.sequence, sequence)
	}

	c.sequence++

	return int(uint32(header[0]) | uint32(header[1])<<8 | uint32(header[2])<<16), nil
}

// ReadEphemeralPacket attempts to read a packet into buffer from sync.Pool.  Do
// not use this method if the contents of the packet needs to be kept
// after the next ReadEphemeralPacket.
//
// Note if the connection is closed already, an error will be
// returned, and it may not be io.EOF. If the connection closes while
// we are stuck waiting for data, an error will also be returned, and
// it most likely will be io.EOF.
func (c *Conn) ReadEphemeralPacket() ([]byte, error) {
	if c.currentEphemeralPolicy != ephemeralUnused {
		panic(fmt.Errorf("ReadEphemeralPacket: unexpected currentEphemeralPolicy: %v", c.currentEphemeralPolicy))
	}
	c.currentEphemeralPolicy = ephemeralRead

	r := c.getReader()
	length, err := c.readHeaderFrom(r)
	if err != nil {
		return nil, err
	}

	if length == 0 {
		// This can be caused by the packet after a packet of
		// exactly size MaxPacketSize.
		return nil, nil
	}

	// Use the bufPool.
	if length < MaxPacketSize {
		c.currentEphemeralBuffer = bufPool.Get(length)
		if _, err := io.ReadFull(r, *c.currentEphemeralBuffer); err != nil {
			return nil, fmt.Errorf("io.ReadFull(packet body of length %v) failed: %v", length, err)
		}
		return *c.currentEphemeralBuffer, nil
	}

	// Much slower path, revert to allocating everything from scratch.
	// We're going to concatenate a lot of data anyway, can't really
	// optimize this code path easily.
	data := make([]byte, length)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, fmt.Errorf("io.ReadFull(packet body of length %v) failed: %v", length, err)
	}
	for {
		next, err := c.readOnePacket()
		if err != nil {
			return nil, err
		}

		if len(next) == 0 {
			// Again, the packet after a packet of exactly size MaxPacketSize.
			break
		}

		data = append(data, next...)
		if len(next) < MaxPacketSize {
			break
		}
	}

	return data, nil
}

// ReadEphemeralPacketDirect attempts to read a packet from the socket directly.
// It needs to be used for the first handshake packet the server receives,
// so we do't buffer the SSL negotiation packet. As a shortcut, only
// packets smaller than MaxPacketSize can be read here.
// This function usually shouldn't be used - use ReadEphemeralPacket.
func (c *Conn) ReadEphemeralPacketDirect() ([]byte, error) {
	if c.currentEphemeralPolicy != ephemeralUnused {
		panic(fmt.Errorf("ReadEphemeralPacketDirect: unexpected currentEphemeralPolicy: %v", c.currentEphemeralPolicy))
	}
	c.currentEphemeralPolicy = ephemeralRead

	var r io.Reader = c.conn
	length, err := c.readHeaderFrom(r)
	if err != nil {
		return nil, err
	}

	if length == 0 {
		// This can be caused by the packet after a packet of
		// exactly size MaxPacketSize.
		return nil, nil
	}

	if length < MaxPacketSize {
		c.currentEphemeralBuffer = bufPool.Get(length)
		if _, err := io.ReadFull(r, *c.currentEphemeralBuffer); err != nil {
			return nil, fmt.Errorf("io.ReadFull(packet body of length %v) failed: %v", length, err)
		}
		return *c.currentEphemeralBuffer, nil
	}

	return nil, fmt.Errorf("ReadEphemeralPacketDirect doesn't support more than one packet")
}

// RecycleReadPacket recycles the read packet. It needs to be called
// after ReadEphemeralPacket was called.
func (c *Conn) RecycleReadPacket() {
	if c.currentEphemeralPolicy != ephemeralRead {
		// Programming error.
		panic(fmt.Errorf("trying to call RecycleReadPacket while currentEphemeralPolicy is %d", c.currentEphemeralPolicy))
	}
	if c.currentEphemeralBuffer != nil {
		// We are using the pool, put the buffer back in.
		bufPool.Put(c.currentEphemeralBuffer)
		c.currentEphemeralBuffer = nil
	}
	c.currentEphemeralPolicy = ephemeralUnused
}

// readOnePacket reads a single packet into a newly allocated buffer.
func (c *Conn) readOnePacket() ([]byte, error) {
	r := c.getReader()
	length, err := c.readHeaderFrom(r)
	if err != nil {
		return nil, err
	}
	if length == 0 {
		// This can be caused by the packet after a packet of
		// exactly size MaxPacketSize.
		return nil, nil
	}

	data := make([]byte, length)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, fmt.Errorf("io.ReadFull(packet body of length %v) failed: %v", length, err)
	}
	return data, nil
}

// readPacket reads a packet from the underlying connection.
// It re-assembles packets that span more than one message.
// This method returns a generic error, not a SQLError.
func (c *Conn) readPacket() ([]byte, error) {
	// Optimize for a single packet case.
	data, err := c.readOnePacket()
	if err != nil {
		return nil, err
	}

	// This is a single packet.
	if len(data) < MaxPacketSize {
		return data, nil
	}

	// There is more than one packet, read them all.
	for {
		next, err := c.readOnePacket()
		if err != nil {
			return nil, err
		}

		if len(next) == 0 {
			// Again, the packet after a packet of exactly size MaxPacketSize.
			break
		}

		data = append(data, next...)
		if len(next) < MaxPacketSize {
			break
		}
	}

	return data, nil
}

// ReadPacket reads a packet from the underlying connection.
// it is the public API version, that returns a SQLError.
// The memory for the packet is always allocated, and it is owned by the caller
// after this function returns.
func (c *Conn) ReadPacket() ([]byte, error) {
	result, err := c.readPacket()
	if err != nil {
		return nil, err
	}
	return result, err
}

// WritePacket writes a packet, possibly cutting it into multiple
// chunks.  Note this is not very efficient, as the client probably
// has to build the []byte and that makes a memory copy.
// Try to use StartEphemeralPacket/writeEphemeralPacket instead.
//
// This method returns a generic error, not a SQLError.
func (c *Conn) WritePacket(data []byte) error {
	index := 0
	length := len(data)

	w := c.getWriter()

	for {
		// Packet length is capped to MaxPacketSize.
		packetLength := length
		if packetLength > MaxPacketSize {
			packetLength = MaxPacketSize
		}

		// Compute and write the header.
		var header [4]byte
		header[0] = byte(packetLength)
		header[1] = byte(packetLength >> 8)
		header[2] = byte(packetLength >> 16)
		header[3] = c.sequence
		if n, err := w.Write(header[:]); err != nil {
			return fmt.Errorf("Write(header) failed: %v", err)
		} else if n != 4 {
			return fmt.Errorf("Write(header) returned a short write: %v < 4", n)
		}

		// Write the body.
		if n, err := w.Write(data[index : index+packetLength]); err != nil {
			return fmt.Errorf("Write(packet) failed: %v", err)
		} else if n != packetLength {
			return fmt.Errorf("Write(packet) returned a short write: %v < %v", n, packetLength)
		}

		// Update our state.
		c.sequence++
		length -= packetLength
		if length == 0 {
			if packetLength == MaxPacketSize {
				// The packet we just sent had exactly
				// MaxPacketSize size, we need to
				// sent a zero-size packet too.
				header[0] = 0
				header[1] = 0
				header[2] = 0
				header[3] = c.sequence
				if n, err := w.Write(header[:]); err != nil {
					return fmt.Errorf("Write(empty header) failed: %v", err)
				} else if n != 4 {
					return fmt.Errorf("Write(empty header) returned a short write: %v < 4", n)
				}
				c.sequence++
			}
			return nil
		}
		index += packetLength
	}
}

// StartEphemeralPacket get []byte from pool
func (c *Conn) StartEphemeralPacket(length int) []byte {
	if c.currentEphemeralPolicy != ephemeralUnused {
		panic("StartEphemeralPacket cannot be used while a packet is already started.")
	}

	c.currentEphemeralPolicy = ephemeralWrite
	// get buffer from pool or it'll be allocated if length is too big
	c.currentEphemeralBuffer = bufPool.Get(length)
	return *c.currentEphemeralBuffer
}

// WriteEphemeralPacket writes the packet that was allocated by
// StartEphemeralPacket.
func (c *Conn) WriteEphemeralPacket() error {
	defer c.recycleWritePacket()

	switch c.currentEphemeralPolicy {
	case ephemeralWrite:
		if err := c.WritePacket(*c.currentEphemeralBuffer); err != nil {
			return fmt.Errorf("Conn %v: %v", c.GetConnectionID(), err)
		}
	case ephemeralUnused, ephemeralRead:
		// Programming error.
		panic(fmt.Errorf("Conn %v: trying to call writeEphemeralPacket while currentEphemeralPolicy is %v", c.GetConnectionID(), c.currentEphemeralPolicy))
	}

	return nil
}

// recycleWritePacket recycles the write packet. It needs to be called
// after writeEphemeralPacket was called.
func (c *Conn) recycleWritePacket() {
	if c.currentEphemeralPolicy != ephemeralWrite {
		// Programming error.
		panic(fmt.Errorf("trying to call recycleWritePacket while currentEphemeralPolicy is %d", c.currentEphemeralPolicy))
	}
	// Release our reference so the buffer can be gced
	bufPool.Put(c.currentEphemeralBuffer)
	c.currentEphemeralBuffer = nil
	c.currentEphemeralPolicy = ephemeralUnused
}

// writeComQuit writes a Quit message for the server, to indicate we
// want to close the connection.
// Client -> Server.
// Returns SQLError(CRServerGone) if it can't.
func (c *Conn) writeComQuit() error {
	// This is a new command, need to reset the sequence.
	c.sequence = 0

	data := c.StartEphemeralPacket(1)
	data[0] = ComQuit
	if err := c.WriteEphemeralPacket(); err != nil {
		return err
	}
	return nil
}

// RemoteAddr returns the underlying socket RemoteAddr().
func (c *Conn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// GetConnectionID returns the MySQL connection ID for this connection.
func (c *Conn) GetConnectionID() uint32 {
	return c.ConnectionID
}

// SetConnectionID set connection id of conn.
func (c *Conn) SetConnectionID(connectionID uint32) {
	c.ConnectionID = connectionID
}

// SetSequence set sequence of conn
func (c *Conn) SetSequence(sequence uint8) {
	c.sequence = sequence
}

// GetSequence return sequence of conn
func (c *Conn) GetSequence() uint8 {
	return c.sequence
}

// Ident returns a useful identification string for error logging
func (c *Conn) String() string {
	return fmt.Sprintf("client %v (%s)", c.ConnectionID, c.RemoteAddr().String())
}

// Close closes the connection. It can be called from a different go
// routine to interrupt the current connection.
func (c *Conn) Close() {
	if c.closed.CompareAndSwap(false, true) {
		c.conn.Close()
	}
}

// IsClosed returns true if this connection was ever closed by the
// Close() method.  Note if the other side closes the connection, but
// Close() wasn't called, this will return false.
func (c *Conn) IsClosed() bool {
	return c.closed.Get()
}

//
// Packet writing methods, for generic packets.
//

// WriteOKPacket writes an OK packet.
// Server -> Client.
// This method returns a generic error, not a SQLError.
func (c *Conn) WriteOKPacket(affectedRows, lastInsertID uint64, flags uint16, warnings uint16) error {
	length := 1 + // OKHeader
		LenEncIntSize(affectedRows) +
		LenEncIntSize(lastInsertID) +
		2 + // flags
		2 // warnings
	data := c.StartEphemeralPacket(length)
	pos := 0
	pos = WriteByte(data, pos, OKHeader)
	pos = WriteLenEncInt(data, pos, affectedRows)
	pos = WriteLenEncInt(data, pos, lastInsertID)
	pos = WriteUint16(data, pos, flags)
	pos = WriteUint16(data, pos, warnings)

	return c.WriteEphemeralPacket()
}

// WriteOKPacketWithEOFHeader writes an OK packet with an EOF header.
// This is used at the end of a result set if
// CapabilityClientDeprecateEOF is set.
// Server -> Client.
// This method returns a generic error, not a SQLError.
func (c *Conn) WriteOKPacketWithEOFHeader(affectedRows, lastInsertID uint64, flags uint16, warnings uint16) error {
	length := 1 + // EOFHeader
		LenEncIntSize(affectedRows) +
		LenEncIntSize(lastInsertID) +
		2 + // flags
		2 // warnings
	data := c.StartEphemeralPacket(length)
	pos := 0
	pos = WriteByte(data, pos, EOFHeader)
	pos = WriteLenEncInt(data, pos, affectedRows)
	pos = WriteLenEncInt(data, pos, lastInsertID)
	pos = WriteUint16(data, pos, flags)
	pos = WriteUint16(data, pos, warnings)

	return c.WriteEphemeralPacket()
}

// WriteErrorPacket writes an error packet.
// Server -> Client.
// This method returns a generic error, not a SQLError.
func (c *Conn) WriteErrorPacket(errorCode uint16, sqlState string, format string, args ...interface{}) error {
	errorMessage := fmt.Sprintf(format, args...)
	length := 1 + 2 + 1 + 5 + len(errorMessage)
	data := c.StartEphemeralPacket(length)
	pos := 0
	pos = WriteByte(data, pos, ErrHeader)
	pos = WriteUint16(data, pos, errorCode)
	pos = WriteByte(data, pos, '#')
	if sqlState == "" {
		sqlState = DefaultMySQLState
	}
	if len(sqlState) != 5 {
		panic("sqlState has to be 5 characters long")
	}
	pos = writeEOFString(data, pos, sqlState)
	pos = writeEOFString(data, pos, errorMessage)

	return c.WriteEphemeralPacket()
}

// WriteErrorPacketFromError writes an error packet, from a regular error.
// See writeErrorPacket for other info.
func (c *Conn) WriteErrorPacketFromError(err error) error {
	if se, ok := err.(*SQLError); ok {
		return c.WriteErrorPacket(se.SQLCode(), se.SQLState(), "%v", se.Message)
	}

	return c.WriteErrorPacket(ErrUnknown, DefaultMySQLState, "unknown error: %v", err)
}

// WriteEOFPacket writes an EOF packet, through the buffer, and
// doesn't flush (as it is used as part of a query result).
func (c *Conn) WriteEOFPacket(flags uint16, warnings uint16) error {
	length := 5
	data := c.StartEphemeralPacket(length)
	pos := 0
	pos = WriteByte(data, pos, EOFHeader)
	pos = WriteUint16(data, pos, warnings)
	pos = WriteUint16(data, pos, flags)

	return c.WriteEphemeralPacket()
}

//
// Packet parsing methods, for generic packets.
//

// isEOFHeader determines whether or not a data packet is a "true" EOF. DO NOT blindly compare the
// first byte of a packet to EOFHeader as you might do for other packet types, as 0xfe is overloaded
// as a first byte.
//
// Per https://dev.mysql.com/doc/internals/en/packet-EOF_Packet.html, a packet starting with 0xfe
// but having length >= 9 (on top of 4 byte header) is not a true EOF but a LengthEncodedInteger
// (typically preceding a LengthEncodedString). Thus, all EOF checks must validate the payload size
// before exiting.
//
// More specifically, an EOF packet can have 3 different lengths (1, 5, 7) depending on the client
// flags that are set. 7 comes from server versions of 5.7.5 or greater where ClientDeprecateEOF is
// set (i.e. uses an OK packet starting with 0xfe instead of 0x00 to signal EOF). Regardless, 8 is
// an upper bound otherwise it would be ambiguous w.r.t. LengthEncodedIntegers.
//
// More docs here:
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_basic_response_packets.html
func isEOFHeader(data []byte) bool {
	return data[0] == EOFHeader && len(data) < 9
}

// parseEOFHeader returns the warning count and a boolean to indicate if there
// are more results to receive.
//
// Note: This is only valid on actual EOF packets and not on OK packets with the EOF
// type code set, i.e. should not be used if ClientDeprecateEOF is set.
func parseEOFHeader(data []byte) (warnings uint16, more bool, err error) {
	// The warning count is in position 2 & 3
	warnings, _, ok := ReadUint16(data, 1)

	// The status flag is in position 4 & 5
	statusFlags, _, ok := ReadUint16(data, 3)
	if !ok {
		return 0, false, fmt.Errorf("invalid EOF packet statusFlags: %v", data)
	}
	return warnings, (statusFlags & ServerMoreResultsExists) != 0, nil
}

func parseOKHeader(data []byte) (uint64, uint64, uint16, uint16, error) {
	// We already read the type.
	pos := 1

	// Affected rows.
	affectedRows, pos, _, ok := ReadLenEncInt(data, pos)
	if !ok {
		return 0, 0, 0, 0, fmt.Errorf("invalid OK packet affectedRows: %v", data)
	}

	// Last Insert ID.
	lastInsertID, pos, _, ok := ReadLenEncInt(data, pos)
	if !ok {
		return 0, 0, 0, 0, fmt.Errorf("invalid OK packet lastInsertID: %v", data)
	}

	// Status flags.
	statusFlags, pos, ok := ReadUint16(data, pos)
	if !ok {
		return 0, 0, 0, 0, fmt.Errorf("invalid OK packet statusFlags: %v", data)
	}

	// Warnings.
	warnings, pos, ok := ReadUint16(data, pos)
	if !ok {
		return 0, 0, 0, 0, fmt.Errorf("invalid OK packet warnings: %v", data)
	}

	return affectedRows, lastInsertID, statusFlags, warnings, nil
}

// IsErrorPacket determines whether or not the packet is an error packet. Mostly here for
// consistency with isEOFHeader
func IsErrorPacket(data []byte) bool {
	return data[0] == ErrHeader
}

// IsOKPacket determines whether or not the packet is an ok packet.
func IsOKPacket(data []byte) bool {
	return data[0] == OKHeader
}

// ParseErrorPacket parses the error packet and returns a SQLError.
func ParseErrorPacket(data []byte) error {
	// We already read the type.
	pos := 1

	// Error code is 2 bytes.
	code, pos, ok := ReadUint16(data, pos)
	if !ok {
		return errors.New("invalid error packet code")
	}

	// '#' marker of the SQL state is 1 byte. Ignored.
	pos++

	// SQL state can be calculated
	_, pos, ok = ReadBytes(data, pos, 5)
	if !ok {
		return errors.New("invalid error packet sqlState")
	}

	// Human readable error message is the rest.
	msg := string(data[pos:])

	return NewError(code, msg)
}
