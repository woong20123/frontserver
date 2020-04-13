package logicmanager

import (
	"log"

	"github.com/woong20123/packet"
)

// NewLogicManager is
func NewLogicManager() *LogicManager {
	lm := LogicManager{}
	lm.LogicConatiner = make(map[uint32]func(p *packet.Packet))
	return &lm
}

// LogicManager is
type LogicManager struct {
	LogicConatiner map[uint32]func(p *packet.Packet)
}

// RegistLogicfun regist packet processing logic
func (lm *LogicManager) RegistLogicfun(cmd uint32, fun func(p *packet.Packet)) {
	lm.LogicConatiner[cmd] = fun
}

// UnregistLogicfun unregist packet processing logic
func (lm *LogicManager) UnregistLogicfun(cmd uint32, fun func(p *packet.Packet)) {
	delete(lm.LogicConatiner, cmd)
}

// CallLogicFun is
func (lm *LogicManager) CallLogicFun(cmd uint32, p *packet.Packet) {
	val, exist := lm.LogicConatiner[cmd]
	if exist {
		val(p)
	} else {
		log.Fatalln("call fail ", cmd)
	}
}
