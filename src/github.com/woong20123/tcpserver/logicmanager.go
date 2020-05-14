package tcpserver

import (
	"net"

	"github.com/woong20123/packet"
)

// TlogicFunc is trans request struct
type TlogicFunc func(conn *net.TCPConn, p *packet.Packet)

// LogicManager is
type ClientLogicManager struct {
	LogicConatiner map[int32]TlogicFunc
	clientRequest  chan *SendToClientRequest
}

// SendToClientRequest is
type SendToClientRequest struct {
	conn *net.TCPConn
	p    *packet.Packet
}

// Initialize is
func (lm *ClientLogicManager) Initialize() {
	lm.LogicConatiner = make(map[int32]TlogicFunc)
	lm.clientRequest = make(chan *SendToClientRequest, 4096)
}

// RegistLogicfun regist packet processing logic
func (lm *ClientLogicManager) RegistLogicfun(cmd int32, fun TlogicFunc) {
	lm.LogicConatiner[cmd] = fun
}

// UnregistLogicfun unregist packet processing logic
func (lm *ClientLogicManager) UnregistLogicfun(cmd int32) {
	delete(lm.LogicConatiner, cmd)
}

// CallLogicFun is
func (lm *ClientLogicManager) CallLogicFun(cmd int32, conn *net.TCPConn, p *packet.Packet) {
	r := SendToClientRequest{conn, p}
	lm.clientRequest <- &r
}

// RunLogicHandle is
func (lm *ClientLogicManager) RunLogicHandler(processCount int) {
	for i := 0; i < processCount; i++ {
		go lm.handleRequest(lm.clientRequest)
	}
}

func (lm *ClientLogicManager) handleRequest(queue chan *SendToClientRequest) {
	for cr := range queue {
		lm.packetProcess(cr)
	}
}

func (lm *ClientLogicManager) packetProcess(cr *SendToClientRequest) {
	cmd := cr.p.Command()
	val, ok := lm.LogicConatiner[cmd]
	if ok {
		val(cr.conn, cr.p)
		packet.Pool().ReleasePacket(cr.p)
	} else {
		Instance().LoggerMgr().Logger().Println("call fail ", cmd)
	}
}

// ======== [ServerLogicManager] =======

// LogicManager is
type ServerLogicManager struct {
	LogicConatiner map[int32]TlogicFunc
	serverRequest  chan *SendToServerRequest
}

// SendToServerRequest is
type SendToServerRequest struct {
	serverIndex int
	conn        *net.TCPConn
	p           *packet.Packet
}

// Initialize is
func (lm *ServerLogicManager) Initialize() {
	lm.LogicConatiner = make(map[int32]TlogicFunc)
	lm.serverRequest = make(chan *SendToServerRequest, 4096)
}

// RegistLogicfun regist packet processing logic
func (lm *ServerLogicManager) RegistLogicfun(cmd int32, fun TlogicFunc) {
	lm.LogicConatiner[cmd] = fun
}

// UnregistLogicfun unregist packet processing logic
func (lm *ServerLogicManager) UnregistLogicfun(cmd int32) {
	delete(lm.LogicConatiner, cmd)
}

// CallLogicFun is
func (lm *ServerLogicManager) CallLogicFun(serverindex int, conn *net.TCPConn, p *packet.Packet) {
	r := SendToServerRequest{serverindex, conn, p}
	lm.serverRequest <- &r
}

// RunLogicHandler is
func (lm *ServerLogicManager) RunLogicHandler(processCount int) {
	for i := 0; i < processCount; i++ {
		go lm.handleRequest(lm.serverRequest)
	}
}

func (lm *ServerLogicManager) handleRequest(queue chan *SendToServerRequest) {
	for cr := range queue {
		lm.packetProcess(cr)
	}
}

func (lm *ServerLogicManager) packetProcess(cr *SendToServerRequest) {
	cmd := cr.p.Command()
	val, ok := lm.LogicConatiner[cmd]
	if ok {
		val(cr.conn, cr.p)
		packet.Pool().ReleasePacket(cr.p)
	} else {
		Instance().LoggerMgr().Logger().Println("call fail ", cmd)
	}
}
