package logicmanager

import (
	"log"
	"net"

	"github.com/woong20123/packet"
)

// TlogicFunc is trans request struct
type TlogicFunc func(conn *net.TCPConn, p *packet.Packet)

// NewLogicManager is
func NewLogicManager() *LogicManager {
	lm := LogicManager{}
	lm.LogicConatiner = make(map[uint32]TlogicFunc)
	lm.clientRequest = make(chan *Request, 2048)
	return &lm
}

// LogicManager is
type LogicManager struct {
	LogicConatiner map[uint32]TlogicFunc
	clientRequest  chan *Request
}

// Request is
type Request struct {
	conn *net.TCPConn
	p    *packet.Packet
}

// RegistLogicfun regist packet processing logic
func (lm *LogicManager) RegistLogicfun(cmd uint32, fun TlogicFunc) {
	lm.LogicConatiner[cmd] = fun
}

// UnregistLogicfun unregist packet processing logic
func (lm *LogicManager) UnregistLogicfun(cmd uint32) {
	delete(lm.LogicConatiner, cmd)
}

// CallLogicFun is
func (lm *LogicManager) CallLogicFun(cmd uint32, conn *net.TCPConn, p *packet.Packet) {
	r := Request{conn, p}
	lm.clientRequest <- &r
}

// RunLogicHandle is
func (lm *LogicManager) RunLogicHandle(processCount int) {
	for i := 0; i < processCount; i++ {
		go lm.handleRequest(lm.clientRequest)
	}
}

func (lm *LogicManager) handleRequest(queue chan *Request) {
	for cr := range queue {
		lm.packetProcess(cr)
	}
}

func (lm *LogicManager) packetProcess(cr *Request) {
	cmd := cr.p.GetCommand()
	val, exist := lm.LogicConatiner[cmd]
	if exist {
		val(cr.conn, cr.p)
	} else {
		log.Println("call fail ", cmd)
	}
}