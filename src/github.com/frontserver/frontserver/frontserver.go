package main

// "factored" import statment

import (
	"io"
	"log"
	"net"
)

const MAX_BUFFERSIZE = 4096

func ConnHandler(conn net.Conn) {
	recvBuf := make([]byte, MAX_BUFFERSIZE)

	for {
		n, err := conn.Read(recvBuf)
		if nil != err {
			if io.EOF == err {
				log.Printf("connection is closed from client; %v", conn.RemoteAddr().String())
				return
			}
			log.Printf("fail to receive data; err: %v", err)
			return
		}

		if 0 < n {
			// data parser
			data := recvBuf[:n]
			log.Println(string(data))
		}
	}
}

func main() {
	server_listen, err := net.Listen("tcp", ":20224")

	if nil != err {
		log.Fatalf("fail to bind address to 20224; err : %v", err)
	}

	defer server_listen.Close()

	for {
		conn, err := server_listen.Accept()
		if nil != err {
			log.Printf("fail to accept; err: %v", err)
			continue
		}

		go ConnHandler(conn)
	}
}
