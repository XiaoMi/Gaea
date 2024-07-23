package zap

import (
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"time"
)

var zapBufferPool = buffer.NewPool()

type ZapEncoder struct {
}

func (e *ZapEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	// 将文本写入到输出流中
	buf := zapBufferPool.Get()
	buf.AppendString("[" + entry.Time.Format("2006-01-02 15:04:05.000") + "] ")
	buf.AppendString("[" + entry.Level.CapitalString() + "] ")
	buf.AppendString(entry.Message + "\n")
	return buf, nil
}

func (e *ZapEncoder) Clone() zapcore.Encoder {
	return e
}

func (e *ZapEncoder) AddArray(key string, arr zapcore.ArrayMarshaler) error {
	return nil
}

func (e *ZapEncoder) AddObject(key string, obj zapcore.ObjectMarshaler) error {
	return nil
}

func (e *ZapEncoder) AddBinary(key string, val []byte) {
}

func (e *ZapEncoder) AddByteString(key string, val []byte) {
}

func (e *ZapEncoder) AddBool(key string, val bool) {
}

func (e *ZapEncoder) AddComplex128(key string, val complex128) {
}

func (e *ZapEncoder) AddDuration(key string, val time.Duration) {
}

func (e *ZapEncoder) AddFloat64(key string, val float64) {
}

func (e *ZapEncoder) AddInt64(key string, val int64) {
}

func (e *ZapEncoder) AddReflected(key string, obj interface{}) error {
	return nil
}

func (e *ZapEncoder) OpenNamespace(key string) {
}

func (e *ZapEncoder) AddString(key, val string) {
}

func (e *ZapEncoder) AddTime(key string, val time.Time) {
}

func (e *ZapEncoder) AddUint64(key string, val uint64) {
}

func (e *ZapEncoder) AddComplex64(k string, v complex64) { e.AddComplex128(k, complex128(v)) }
func (e *ZapEncoder) AddFloat32(k string, v float32)     { e.AddFloat64(k, float64(v)) }
func (e *ZapEncoder) AddInt(k string, v int)             { e.AddInt64(k, int64(v)) }
func (e *ZapEncoder) AddInt32(k string, v int32)         { e.AddInt64(k, int64(v)) }
func (e *ZapEncoder) AddInt16(k string, v int16)         { e.AddInt64(k, int64(v)) }
func (e *ZapEncoder) AddInt8(k string, v int8)           { e.AddInt64(k, int64(v)) }
func (e *ZapEncoder) AddUint(k string, v uint)           { e.AddUint64(k, uint64(v)) }
func (e *ZapEncoder) AddUint32(k string, v uint32)       { e.AddUint64(k, uint64(v)) }
func (e *ZapEncoder) AddUint16(k string, v uint16)       { e.AddUint64(k, uint64(v)) }
func (e *ZapEncoder) AddUint8(k string, v uint8)         { e.AddUint64(k, uint64(v)) }
func (e *ZapEncoder) AddUintptr(k string, v uintptr)     { e.AddUint64(k, uint64(v)) }
