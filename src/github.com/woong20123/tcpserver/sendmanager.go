package tcpserver

import (
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
func (sm *SendManager) SendToClientConn(s Session, p *packet.Packet) {
	r := SendToClientRequest{s, p}
	sm.SendtoClientRequest <- &r
}

// RunSendToClientHandle is
func (sm *SendManager) RunSendToClientHandle(processCount int) {
	for i := 0; i < processCount; i++ {
		go handleRequestProcess(sm.SendtoClientRequest, func(cr *SendToClientRequest) {
			if cr != nil && cr.s != nil {
				tcpclisession := cr.s.(*TCPClientSession)
				if nil != tcpclisession && tcpclisession.Conn() != nil {
					_, err := tcpclisession.Conn().Write(cr.p.MakeByte())
					if err != nil {
						Instance().LoggerMgr().Logger().Println("RunSendHandle p command = ", cr.p.Command(), " err = ", err)
					}
					packet.Pool().ReleasePacket(cr.p)
				}

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
	session, err := Instance().TCPServerMgr().TCPServerSession(index)
	if err == nil {
		r := SendToServerRequest{index, session, p}
		sm.SendtoServerRequest <- &r
	}
}

// RunSendToServerHandle is
func (sm *SendManager) RunSendToServerHandle(Serveridx int) {
	go func() {
		tcpserversession, err := Instance().TCPServerMgr().TCPServerSession(Serveridx)

		// 서버에 연결한 세션을 가져 오지 못했다.?
		if err != nil {
			Instance().LoggerMgr().Logger().Println(err.Error())
			return
		}

		for cr := range sm.SendtoServerRequest {
			if cr != nil && tcpserversession.Conn() != nil {
				_, err := tcpserversession.Write(cr.p.MakeByte())
				if err != nil {
					Instance().LoggerMgr().Logger().Println("RunSendHandle p command = ", cr.p.Command(), " err = ", err)
				}
				packet.Pool().ReleasePacket(cr.p)
			}
		}
	}()
}
