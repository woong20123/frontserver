package main

// "factored" import statment

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

const (
	listenerCloseMatcher = "use of closed network connection"
	maxBufferSize        = 4096
)

// // ConnHandler 연결된 객체를 처리합니다.
// func ConnHandler(conn net.Conn) {
// 	recvBuf := make([]byte, maxBufferSize)

// 	for {
// 		n, err := conn.Read(recvBuf)
// 		if nil != err {
// 			if io.EOF == err {
// 				log.Printf("connection is closed from client; %v", conn.RemoteAddr().String())
// 				return
// 			}
// 			log.Printf("fail to receive data; err: %v", err)
// 			return
// 		}

// 		if 0 < n {
// 			// data parser
// 			data := recvBuf[:n]
// 			log.Println(string(data))
// 		}
// 	}
// }

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

func HandleConnection(conn *net.TCPConn, serverCtx context.Context, wg *sync.WaitGroup) {
	defer func() {
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

func HandleListener(tcpListen *net.TCPListener, serverCtx context.Context, wg *sync.WaitGroup, chClosed chan struct{}) {

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
				case <-serverCtx.Done():
					return
				default:
					// fallthrough
				}
			}

			log.Println("AcceptTcp", err)
			return
		}
		wg.Add(1)
		go HandleConnection(conn, serverCtx, wg)
	}
}

func listenerCloseError(err error) bool {
	return strings.Contains(err.Error(), listenerCloseMatcher)
}

func main() {

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":20224")

	if err != nil {
		log.Println("ResolveTCPAddr", err)
		return
	}

	tcpListen, err := net.ListenTCP("tcp", tcpAddr)

	if nil != err {
		log.Println("ListenTCP", err)
		return
	}
	defer tcpListen.Close()

	sigChan := make(chan os.Signal, 1)
	signal.Ignore()
	signal.Notify(sigChan, syscall.SIGINT)

	var wg sync.WaitGroup
	chClosed := make(chan struct{})

	serverCtx, shutdown := context.WithCancel(context.Background())

	go HandleListener(tcpListen, serverCtx, &wg, chClosed)

	log.Println("On Server ", tcpAddr)

	s := <-sigChan

	switch s {
	case syscall.SIGINT:
		log.Println("Server shutdown...")
		shutdown()
		tcpListen.Close()

		wg.Wait()
		<-chClosed
		log.Println("Server shutdown completed")
	default:
		panic("unexpected signal has been received")
	}

}
