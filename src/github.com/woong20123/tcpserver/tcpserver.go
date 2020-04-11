package tcpserver

import (
	"context"
	"log"
	"net"
	"strings"
	"sync"
)

const (
	listenerCloseMatcher = "use of closed network connection"
	maxBufferSize        = 4096
)

// HandleRead 연결된 세션에 대한 패킷 Read 작업을 처리합니다.
func HandleRead(conn *net.TCPConn, errRead context.CancelFunc) {
	defer errRead()

	recvBuf := make([]byte, maxBufferSize)

	for {
		n, err := conn.Read(recvBuf)
		if err != nil {
			if ne, ok := err.(net.Error); ok {
				switch {
				case ne.Temporary():
					continue
				}
			}

			log.Println("Read", err)
			return
		}

		n, err = conn.Write(recvBuf[:n])
		if err != nil {
			log.Println("Write", err)
			return
		}

		if 0 < n {
			data := recvBuf[:n]
			log.Println(string(data))
		}
	}
}

// HandleConnection은 연결된 세션에 대한 작업을 등록합니다.
func HandleConnection(conn *net.TCPConn, serverCtx context.Context, wg *sync.WaitGroup) {
	defer func() {
		conn.Close()
		wg.Done()
	}()

	readCtx, errRead := context.WithCancel(context.Background())

	go HandleRead(conn, errRead)

	select {
	case <-readCtx.Done():
		log.Println("readCtx.Done()")
	case <-serverCtx.Done():
		log.Println("serverCtx.Done()")
	}
}

func listenerCloseError(err error) bool {
	return strings.Contains(err.Error(), listenerCloseMatcher)
}

// HandleListener register the task to listen to the socket
// HandleListener은 전달된 server address로 소켓을 Listen하는 작업 등록합니다.
func HandleListener(ctxServer context.Context, address string, wg *sync.WaitGroup, chClosed chan struct{}) {

	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		log.Println("ResolveTCPAddr", err)
		return
	}
	tcpListen, err := net.ListenTCP("tcp", tcpAddr)

	if nil != err {
		log.Println("ListenTCP", err)
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
					log.Println("AcceptTCP", err)
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

			log.Println("AcceptTcp", err)
			return
		}
		wg.Add(1)
		go HandleConnection(conn, ctxServer, wg)
	}
}
