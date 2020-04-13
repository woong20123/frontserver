package tcpserver

import "github.com/woong20123/logicmanager"

type StaticData struct {
	packetSerialkey uint32
	lm              *logicmanager.LogicManager
}

func (sd *StaticData) SetSerialkey(key uint32) {
	sd.packetSerialkey = key
}

// GetSerialkey is
func (sd *StaticData) GetSerialkey() uint32 {
	return sd.packetSerialkey
}
