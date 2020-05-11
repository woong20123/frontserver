package tcpserver

import (
	"log"
	"net"
)

// TCPClientSession is
type TCPClientSession struct {
	conn *net.TCPConn
	idx  int
}

// NewTCPClientSession is make TCPClient
func NewTCPClientSession() *TCPClientSession {
	tc := TCPClientSession{}
	tc.conn = nil
	tc.idx = -1

	return &tc
}

// Conn is
func (tc *TCPClientSession) Conn() *net.TCPConn {
	return tc.conn
}

// SetConn is
func (tc *TCPClientSession) SetConn(conn *net.TCPConn) {
	tc.conn = conn
}

// Index is return TCPClent index
func (tc *TCPClientSession) Index() int {
	return tc.idx
}

// SetIndex is
func (tc *TCPClientSession) SetIndex(idx int) {
	tc.idx = idx
}

func (tc *TCPClientSession) connect(addr net.TCPAddr) error {
	var err error
	tc.conn, err = net.DialTCP("tcp", nil, &addr)
	return err
}

// Connect is TcpClient connet to target server
func (tc *TCPClientSession) Connect(ip string, port int) error {
	var remoteaddr net.TCPAddr
	remoteaddr.IP = net.ParseIP(ip)
	remoteaddr.Port = port
	err := tc.connect(remoteaddr)

	if err == nil {
	}

	return err
}

// Close is TcpClient close logic
func (tc *TCPClientSession) Close() (err error) {

	if tc.conn != nil {
		Instance().ServerProxySessionHandler().RunConnectFunc(SessionStateEnum.OnClosed, tc)
		tc.conn.Close()
		tc.conn = nil
	}
	return
}

// Read is
func (tc *TCPClientSession) Read(b []byte) (n int, err error) {
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
func (tc *TCPClientSession) HandleRead() {

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
func (tc *TCPClientSession) Write(b []byte) (n int, err error) {
	n, err = tc.conn.Write(b)

	// 전송중에 에러가 발생했습니다.
	if err != nil {
		log.Println("Write", err)
	}
	return
}
