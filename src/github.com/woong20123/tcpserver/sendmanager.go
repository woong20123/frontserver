package tcpserver

import (
	"net"

	"github.com/woong20123/packet"
)

// SendManager is
type SendManager struct {
	serverRequest chan *Request
}

// Initialize is
func (sm *SendManager) Initialize() {
	sm.serverRequest = make(chan *Request, 4096)
}

// SendToConn is
func (sm *SendManager) SendToConn(conn *net.TCPConn, p *packet.Packet) {
	r := Request{conn, p}
	sm.serverRequest <- &r
}

// RunSendHandle is
func (sm *SendManager) RunSendHandle(processCount int) {
	for i := 0; i < processCount; i++ {
		go handleRequestProcess(sm.serverRequest, func(cr *Request) {
			if cr != nil && cr.conn != nil {
				_, err := cr.conn.Write(cr.p.Byte())
				if err != nil {
					Instance().GetLoggerMgr().GetLogger().Println("RunSendHandle p command = ", cr.p.Command(), " err = ", err)
				}
				packet.Pool().ReleasePacket(cr.p)
			}
		})
	}
}

func handleRequestProcess(queue chan *Request, process func(cr *Request)) {
	for cr := range queue {
		process(cr)
	}
}
