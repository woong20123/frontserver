package examclientlogic

import (
	"example/examClient/clientuser"
	"example/share"
	"fmt"
	"net"
	"strings"

	"github.com/woong20123/packet"
	"github.com/woong20123/tcpserver"
)

// ContructLogicManager is register a logic for the command
func ContructLogicManager(lm *tcpserver.LogicManager) {

	// S2CPacketCommandLoginUserRes에 대한 처리 작업을 등록합니다.
	lm.RegistLogicfun(share.S2CPacketCommandLoginUserRes, func(conn *net.TCPConn, p *packet.Packet) {

		res := share.S2CPCLoginUserRes{}
		p.Read(&res.Result, &res.UserSn)
		res.UserID = p.ReadString()

		if res.Result == share.ResultSuccess {
			eu := GetInstance().GetObjMgr().GetUser()
			eu.SetID(res.UserID)
			eu.SetSn(res.UserSn)
			GetInstance().GetObjMgr().GetChanManager().SendChanUserState(clientuser.UserStateEnum.LoginSTATE, fmt.Sprintln("Login에 성공하였습니다.\n[유저정보] SN = ", res.UserSn, " ID = ", res.UserID))
		} else {
			GetInstance().GetObjMgr().GetChanManager().SendChanUserState(clientuser.UserStateEnum.ConnectedSTATE, "Login에 실패하였습니다.")
		}
		return
	})

	// S2CPacketCommandGolobalMsgRes에 대한 처리 작업을 등록합니다.
	lm.RegistLogicfun(share.S2CPacketCommandGolobalMsgRes, func(conn *net.TCPConn, p *packet.Packet) {
		res := share.S2CPCGolobalSendMsgRes{}
		p.Read(&res.Result, &res.Userid, &res.Msg)

		var sb strings.Builder
		sb.WriteString(res.Userid)
		sb.WriteString(" : ")
		sb.WriteString(res.Msg)
		fmt.Println(sb.String())
		return
	})

	lm.RunLogicHandle(1)
}
