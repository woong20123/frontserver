package tcpserver

import (
	"errors"
	"log"
	"net"
)

// Session is
type Session interface {
	SetConn(conn *net.TCPConn)
	Conn() *net.TCPConn
	Index() int
	SetIndex(idx int)
	Connect(ip string, port int) error
	Close() (err error)
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
}

// TCPServerSession is
type TCPServerSession struct {
	conn *net.TCPConn
	idx  int
}

// NewTCPServerSession is make TCPClient
func NewTCPServerSession() *TCPServerSession {
	tss := TCPServerSession{}
	tss.conn = nil
	tss.idx = -1

	return &tss
}

// Conn is
func (tc *TCPServerSession) Conn() *net.TCPConn {
	return tc.conn
}

// SetConn is
func (tc *TCPServerSession) SetConn(conn *net.TCPConn) {
	tc.conn = conn
}

// Index is return TCPClent index
func (tc *TCPServerSession) Index() int {
	return tc.idx
}

// SetIndex is
func (tc *TCPServerSession) SetIndex(idx int) {
	tc.idx = idx
}

// Connect is TcpClient connet to target server
func (tc *TCPServerSession) Connect(ip string, port int) error {
	var err error
	err = errors.New("Server Session Not Connect")
	return err
}

// Close is TcpClient close logic
func (tc *TCPServerSession) Close() (err error) {

	if tc.conn != nil {
		Instance().ServerProxySessionHandler().RunConnectFunc(SessionStateEnum.OnClosed, tc)
		tc.conn.Close()
		tc.conn = nil
	}
	return
}

// Read is
func (tc *TCPServerSession) Read(b []byte) (n int, err error) {
	n, err = tc.conn.Read(b)
	if err != nil {
		if ne, ok := err.(net.Error); ok {
			switch {
			case ne.Temporary():
				return
			}
			log.Println("Read", err)
			return
		}
	}
	return
}

// HandleRead is
func (tc *TCPServerSession) HandleRead() {

	// sessesion을 통해서 전달받기 위한 버퍼 생성
	recvBuf := make([]byte, maxBufferSize)

	// session으로부터 전달받은 버퍼를 packet형태로 변환처리하기 위한 Packet
	// TCP의 데이터 전달이 패킷단위로 전달되지 않기 때문에 조립 작업을 합니다.
	AssemblyBuf := make([]byte, maxBufferSize+128)
	var AssemPos uint32 = 0

	for {
		if tc.conn == nil {
			Instance().LoggerMgr().Logger().Println("conn == nil")
			break
		}

		n, err := tc.Read(recvBuf)
		if err != nil {
			break
		}

		if 0 < n {
			copylength := copy(AssemblyBuf[AssemPos:], recvBuf[:n])
			AssemPos += uint32(copylength)

			AssemPos = Instance().ServerProxySessionHandler().RunRecvFunc(tc, AssemblyBuf, AssemPos)
		}
	}
	return
}

// Write is
func (tc *TCPServerSession) Write(b []byte) (n int, err error) {
	n, err = tc.conn.Write(b)

	// 전송중에 에러가 발생했습니다.
	if err != nil {
		log.Println("Write", err)
	}
	return
}
