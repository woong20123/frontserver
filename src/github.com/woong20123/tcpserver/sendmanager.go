package tcpserver

import (
	"net"

	"github.com/woong20123/packet"
)

// SendManager is
type SendManager struct {
	SendtoClientRequest chan *SendToClientRequest
	SendtoServerRequest chan *SendToServerRequest
}

// Initialize is
func (sm *SendManager) Initialize() {
	sm.SendtoClientRequest = make(chan *SendToClientRequest, 4096)
	sm.SendtoServerRequest = make(chan *SendToServerRequest, 4096)
}

// SendToClientConn is
func (sm *SendManager) SendToClientConn(conn *net.TCPConn, p *packet.Packet) {
	r := SendToClientRequest{conn, p}
	sm.SendtoClientRequest <- &r
}

// RunSendToClientHandle is
func (sm *SendManager) RunSendToClientHandle(processCount int) {
	for i := 0; i < processCount; i++ {
		go handleRequestProcess(sm.SendtoClientRequest, func(cr *SendToClientRequest) {
			if cr != nil && cr.conn != nil {
				_, err := cr.conn.Write(cr.p.Byte())
				if err != nil {
					Instance().LoggerMgr().Logger().Println("RunSendHandle p command = ", cr.p.Command(), " err = ", err)
				}
				packet.Pool().ReleasePacket(cr.p)
			}
		})
	}
}

func handleRequestProcess(queue chan *SendToClientRequest, process func(cr *SendToClientRequest)) {
	for cr := range queue {
		process(cr)
	}
}

// SendToServerConn is
func (sm *SendManager) SendToServerConn(index int, p *packet.Packet) {
	r := SendToServerRequest{index, p}
	sm.SendtoServerRequest <- &r
}

// RunSendToServerHandle is
func (sm *SendManager) RunSendToServerHandle(Serveridx uint32) {
	go func() {
		tcpclient, err := Instance().TCPClientMgr().TCPClient(Serveridx)

		// 서버에 연결한 세션을 가져 오지 못했다.?
		if err != nil {
			Instance().LoggerMgr().Logger().Println(err.Error())
		}

		for cr := range sm.SendtoServerRequest {
			if cr != nil && tcpclient.Conn() != nil {
				_, err := tcpclient.Write(cr.p.Byte())
				if err != nil {
					Instance().LoggerMgr().Logger().Println("RunSendHandle p command = ", cr.p.Command(), " err = ", err)
				}
				packet.Pool().ReleasePacket(cr.p)
			}
		}
	}()
}
