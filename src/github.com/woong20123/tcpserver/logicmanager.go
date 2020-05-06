package tcpserver

import (
	"net"

	"github.com/woong20123/packet"
)

// TlogicFunc is trans request struct
type TlogicFunc func(conn *net.TCPConn, p *packet.Packet)

// LogicManager is
type LogicManager struct {
	LogicConatiner map[int32]TlogicFunc
	clientRequest  chan *SendToClientRequest
}

// SendToClientRequest is
type SendToClientRequest struct {
	conn *net.TCPConn
	p    *packet.Packet
}

// SendToServerRequest is
type SendToServerRequest struct {
	serverIndex uint32
	p           *packet.Packet
}

// Initialize is
func (lm *LogicManager) Initialize() {
	lm.LogicConatiner = make(map[int32]TlogicFunc)
	lm.clientRequest = make(chan *SendToClientRequest, 4096)
}

// RegistLogicfun regist packet processing logic
func (lm *LogicManager) RegistLogicfun(cmd int32, fun TlogicFunc) {
	lm.LogicConatiner[cmd] = fun
}

// UnregistLogicfun unregist packet processing logic
func (lm *LogicManager) UnregistLogicfun(cmd int32) {
	delete(lm.LogicConatiner, cmd)
}

// CallLogicFun is
func (lm *LogicManager) CallLogicFun(cmd int32, conn *net.TCPConn, p *packet.Packet) {
	r := SendToClientRequest{conn, p}
	lm.clientRequest <- &r
}

// RunLogicHandle is
func (lm *LogicManager) RunLogicHandle(processCount int) {
	for i := 0; i < processCount; i++ {
		go lm.handleRequest(lm.clientRequest)
	}
}

func (lm *LogicManager) handleRequest(queue chan *SendToClientRequest) {
	for cr := range queue {
		lm.packetProcess(cr)
	}
}

func (lm *LogicManager) packetProcess(cr *SendToClientRequest) {
	cmd := cr.p.Command()
	val, ok := lm.LogicConatiner[cmd]
	if ok {
		val(cr.conn, cr.p)
		packet.Pool().ReleasePacket(cr.p)
	} else {
		Instance().LoggerMgr().Logger().Println("call fail ", cmd)
	}
}
