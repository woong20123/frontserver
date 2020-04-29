package tcpclient

import (
	"log"
	"net"
)

type TcpClient struct {
	conn *net.TCPConn
}

// NewTcpClient is make TcpClient
func NewTcpClient() *TcpClient {
	tc := TcpClient{}
	tc.conn = nil

	return &tc
}

// Conn is
func (tc *TcpClient) Conn() *net.TCPConn {
	return tc.conn
}

func (tc *TcpClient) connect(addr net.TCPAddr) error {
	var err error
	tc.conn, err = net.DialTCP("tcp", nil, &addr)
	return err
}

// Connect is TcpClient connet to target server
func (tc *TcpClient) Connect(ip string, port int) error {
	var remoteaddr net.TCPAddr
	remoteaddr.IP = net.ParseIP(ip)
	remoteaddr.Port = port
	return tc.connect(remoteaddr)
}

// Close is TcpClient close logic
func (tc *TcpClient) Close() (err error) {

	if tc.conn != nil {
		tc.conn.Close()
		tc.conn = nil
	}
	return
}

// Read is
func (tc *TcpClient) Read(b []byte) (n int, err error) {
	n, err = tc.conn.Read(b)
	if err != nil {
		if ne, ok := err.(net.Error); ok {
			switch {
			case ne.Temporary():
				return
			}
		}
		log.Println("Read", err)
		return
	}
	return
}

// Write is
func (tc *TcpClient) Write(b []byte) (n int, err error) {
	n, err = tc.conn.Write(b)

	// 전송중에 에러가 발생했습니다.
	if err != nil {
		log.Println("Write", err)
	}
	return
}
