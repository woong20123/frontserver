package tcpserver

import (
	"context"
	"net"
	"strings"
	"sync"
)

const (
	listenerCloseMatcher = "use of closed network connection"
	maxBufferSize        = 4096
)

// Consturct is
func Consturct(serialKey uint32, SendProcessCount int) {
	Instance().SetSerialkey(serialKey)
	Instance().SendManager().RunSendHandle(SendProcessCount)
}

// HandleRead handles packet read operations for connected sessions
// kor : HandleRead 연결된 세션에 대한 패킷 Read 작업을 처리합니다.
func HandleRead(conn *net.TCPConn, errRead context.CancelFunc) {
	defer errRead()

	// sessesion을 통해서 전달받기 위한 버퍼 생성
	recvBuf := make([]byte, maxBufferSize)

	// session으로부터 전달받은 버퍼를 packet형태로 변환처리하기 위한 Packet
	// TCP의 데이터 전달이 패킷단위로 전달되지 않기 때문에 조립 작업을 합니다.
	AssemblyBuf := make([]byte, maxBufferSize+128)
	var AssemPos uint32 = 0

	for {
		if conn == nil {
			Instance().GetLoggerMgr().GetLogger().Println("conn == nil")
			return
		}
		n, err := conn.Read(recvBuf)
		if err != nil {
			if ne, ok := err.(net.Error); ok {
				switch {
				case ne.Temporary():
					continue
				}
			}

			Instance().GetLoggerMgr().GetLogger().Println("Read", err)
			return
		}

		if 0 < n {
			copylength := copy(AssemblyBuf[AssemPos:], recvBuf[:n])
			AssemPos += uint32(copylength)

			AssemPos = Instance().SessionMgr().RunRecvFunc(conn, AssemblyBuf, AssemPos)
		}
	}
}

// HandleConnection register job for connected session
// kor : HandleConnection은 연결된 세션에 대한 작업을 등록합니다.
func HandleConnection(serverCtx context.Context, conn *net.TCPConn, wg *sync.WaitGroup) {
	defer func() {
		Instance().SessionMgr().RunConnectFunc(SessionStateEnum.OnClosed, conn)
		conn.Close()
		wg.Done()
	}()

	readCtx, errRead := context.WithCancel(context.Background())

	go HandleRead(conn, errRead)

	select {
	case <-readCtx.Done():
	case <-serverCtx.Done():
	}
}

func listenerCloseError(err error) bool {
	return strings.Contains(err.Error(), listenerCloseMatcher)
}

// HandleListener register the task to listen to the socket
// kor : HandleListener은 전달된 server address로 소켓을 Listen하는 작업 등록합니다.
func HandleListener(ctxServer context.Context, address string, wg *sync.WaitGroup, chClosed chan struct{}) {

	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		Instance().GetLoggerMgr().GetLogger().Println("ResolveTCPAddr", err)
		return
	}
	tcpListen, err := net.ListenTCP("tcp", tcpAddr)

	if nil != err {
		Instance().GetLoggerMgr().GetLogger().Println("ListenTCP", err)
		return
	}

	// if HandleListener close, it process
	defer func() {
		tcpListen.Close()
		close(chClosed)
	}()

	for {
		conn, err := tcpListen.AcceptTCP()

		// if occur error
		if err != nil {
			if ne, ok := err.(net.Error); ok {
				if ne.Temporary() {
					Instance().GetLoggerMgr().GetLogger().Println("AcceptTCP", err)
					continue
				}
			}
			if listenerCloseError(err) {
				select {
				case <-ctxServer.Done():
					return
				default:
					// fallthrough
				}
			}

			Instance().GetLoggerMgr().GetLogger().Println("AcceptTcp", err)
			return
		}
		Instance().SessionMgr().RunConnectFunc(SessionStateEnum.OnConnected, conn)
		wg.Add(1)
		go HandleConnection(ctxServer, conn, wg)
	}
}
