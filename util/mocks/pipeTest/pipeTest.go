package pipeTest

import (
	"bufio"
	"github.com/stretchr/testify/require"
	"net"
	"sync"
	"testing"
)

// ReplyMsgFuncType 回应函数的型态，测试时，当客户端或服务端接收到讯息时，可以利用此函数去建立回传讯息
type ReplyMsgFuncType func([]uint8) []uint8

// TestReplyMsgFunc　，目前是用于验证测试流程是否正确，在这里会处理常接收到什么讯息，要将下来跟着回应什么讯息
// 每次的回应讯息为接收讯息加 1
//     比如 当接收值为 1，就会回传值为 2 给对方
//     比如 当接收值为 2，就会回传值为 3 给对方
func TestReplyMsgFunc(data []uint8) []uint8 {
	return []uint8{data[0] + 1} // 回应讯息为接收讯息加 1
}

// DcMocker 用来模拟数据库服务器的读取和回应的类
type DcMocker struct {
	t         *testing.T       // 单元测试的类
	bufReader *bufio.Reader    // 有缓存的读取 (接收端)
	bufWriter *bufio.Writer    // 有缓存的写入 (传送端)
	connRead  net.Conn         // pipe 的读取连线 (接收端)
	connWrite net.Conn         // pipe 的写入连线 (传送端)
	wg        *sync.WaitGroup  // 在测试流程的操作边界等待
	replyFunc ReplyMsgFuncType // 设定相对应的回应函数
	err       error            // 错误
}

// NewDcServerClient 产生直连 DC 模拟双方对象，包含 客户端对象 和 服务端对象
func NewDcServerClient(t *testing.T, reply ReplyMsgFuncType) (mockClient *DcMocker, mockServer *DcMocker) {
	// 先产生两组 Pipe
	read0, write0 := net.Pipe() // 第一组 Pipe
	read1, write1 := net.Pipe() // 第二组 Pipe

	// 产生客户端和服务端直连 DC 模拟双方对象，分别为 mockClient 和 mockServer
	mockClient = NewDcMocker(t, read0, write1, reply) // 客户端
	mockServer = NewDcMocker(t, read1, write0, reply) // 服务端

	// 结束
	return
}

// NewDcMocker 产生新的直连 dc 模拟对象
func NewDcMocker(t *testing.T, connRead, connWrite net.Conn, reply ReplyMsgFuncType) *DcMocker {
	return &DcMocker{
		t:         t,                          // 单元测试的对象
		bufReader: bufio.NewReader(connRead),  // 服务器的读取 (实现缓存)
		bufWriter: bufio.NewWriter(connWrite), // 服务器的写入 (实现缓存)
		connRead:  connRead,                   // pipe 的读取连线 (接收端)
		connWrite: connWrite,                  // pipe 的写入连线 (传送端)
		wg:        &sync.WaitGroup{},          // 在测试流程的操作边界等待
		replyFunc: reply,                      // 设定相对应的回应函数
	}
}

// GetConnRead 为获得直连 dc 模拟对象的读取连线
func (dcM *DcMocker) GetConnRead() net.Conn {
	return dcM.connRead
}

// GetConnWrite 为获得直连 dc 模拟对象的写入连线
func (dcM *DcMocker) GetConnWrite() net.Conn {
	return dcM.connWrite
}

// GetBufReader 为获得直连 dc 模拟对象的缓存读取
func (dcM *DcMocker) GetBufReader() *bufio.Reader {
	return dcM.bufReader
}

// GetBufWriter 为获得直连 dc 模拟对象的缓存写入
func (dcM *DcMocker) GetBufWriter() *bufio.Writer {
	return dcM.bufWriter
}

// OverwriteConnBufRead 为临时覆写取代直连 dc 模拟对象的 读取连线 connRead 或者是 缓存读取 bufReader
func (dcM *DcMocker) OverwriteConnBufRead(connRead net.Conn, bufReader *bufio.Reader) error {
	// 先进行修改
	if connRead != nil {
		dcM.connRead = connRead // 修改读取连线
	}
	if bufReader != nil {
		dcM.bufReader = bufReader // 修改读取缓存
	}

	// 正确回传
	return nil
}

// OverwriteConnBufWrite 为临时覆写取代直连 dc 模拟对象的 写入连线 connWrite 或者是 缓存写入 bufWriter
func (dcM *DcMocker) OverwriteConnBufWrite(connWrite net.Conn, bufWriter *bufio.Writer) error {
	// 如果 写入连线 connWrite 参数传入为空值时，则进行修改
	if connWrite != nil {
		dcM.connWrite = connWrite
	}
	// 如果 缓存写入 bufWriter 参数传入为空值时，则进行修改
	if bufWriter != nil {
		dcM.bufWriter = bufWriter
	}

	// 正确回传
	return nil
}

// ResetDcMockers 为重置单一连线方向的直连 dc 模拟对象
func (dcM *DcMocker) ResetDcMockers(otherSide *DcMocker) error {
	// 重新建立全新两组 Pipe
	newRead, newWrite := net.Pipe() // 第一组 Pipe

	// 单方向的状况为 dcM 写入 Pipe，otherSide 读取 Pipe

	// 先重置 发送讯息的那一方 部份
	dcM.bufWriter = bufio.NewWriter(newWrite) // 服务器的回应 (实现缓存)
	dcM.connWrite = newWrite                  // pipe 的写入连线
	dcM.wg = &sync.WaitGroup{}                // 流程的操作边界

	// 先重置 mockServer 部份
	otherSide.bufReader = bufio.NewReader(newRead) // 服务器的读取 (实现缓存)
	otherSide.connRead = newRead                   // pipe 的读取连线

	// 正常回传
	return nil
}

// SendOrReceiveMsg 为直连 dc 用来模拟接收或传入讯息
// 比如客户端 "传送 Send" 讯息到服务端、客户端再 "接收 Receive" 服务端的回传讯息
func (dcM *DcMocker) SendOrReceiveMsg(data []uint8) *DcMocker {
	// dc 模拟开始
	dcM.wg.Add(1) // 只要等待直到确认资料有写入 pipe

	// 在这里执行 1传送讯息 或者是 2接收讯息
	go func() {
		// 执行写入工作
		_, err := dcM.bufWriter.Write(data) // 写入资料到 pipe
		err = dcM.bufWriter.Flush()         // 把缓存资料写进 pipe
		require.Equal(dcM.t, err, nil)
		err = dcM.connWrite.Close() // 资料写入完成，终结连线
		require.Equal(dcM.t, err, nil)

		// 写入工作完成
		dcM.wg.Done()
	}()

	// 重复使用对象
	return dcM
}

// UseAnonymousFuncSendMsg 使用匿名函式去傳送訊息
func (dcM *DcMocker) UseAnonymousFuncSendMsg(customFunc func()) *DcMocker {
	// dc 模拟开始
	dcM.wg.Add(1) // 只要等待直到确认资料有写入 pipe

	// 在这里执行 1传送讯息 或者是 2接收讯息
	go func() {
		customFunc()

		// 写入工作完成
		dcM.wg.Done()
	}()

	// 重复使用对象
	return dcM
}

// ReplyMsg 为直连 dc 用来模拟回应数据，大部份接连 SendOrReceiveMsg 函数后执行
func (dcM *DcMocker) ReplyMsg(otherSide *DcMocker) (msg []uint8) {
	// 读取传送过来的讯息
	b, _, err := otherSide.bufReader.ReadLine() // 由另一方接收传来的讯息
	require.Equal(dcM.t, err, nil)

	// 等待和确认资料已经写入 pipe
	dcM.wg.Wait()

	// 重置模拟对象
	err = dcM.ResetDcMockers(otherSide)
	require.Equal(dcM.t, err, nil)

	// 回传回应讯息
	if otherSide.replyFunc != nil {
		msg = otherSide.replyFunc(b)
	}

	// 结束
	return
}

// CheckArrivedMsg 是用于当模拟时，直连 dc 讯息传送到对方时，立刻对传送到的讯息进行查询，以方便后续的除错和检查
func (dcM *DcMocker) CheckArrivedMsg(otherSide *DcMocker) (msg []uint8) {
	// 读取传送过来的讯息
	b, _, err := otherSide.bufReader.ReadLine() // 由另一方接收传来的讯息
	require.Equal(dcM.t, err, nil)

	// 等待和确认资料已经写入 pipe
	dcM.wg.Wait()

	// 重置模拟对象
	err = dcM.ResetDcMockers(otherSide)
	require.Equal(dcM.t, err, nil)

	// 回传回应讯息
	msg = b

	// 结束
	return
}

// WaitAndReset 为直连 dc 用来等待在 Pipe 的整个数据读写操作完成
func (dcM *DcMocker) WaitAndReset(otherSide *DcMocker) error {
	// 先等待整个数据读写操作完成
	dcM.wg.Wait()

	// 单方向完成 Pipe 的连线重置
	err := dcM.ResetDcMockers(otherSide)
	require.Equal(dcM.t, err, nil)

	// 正确回传
	return nil
}
