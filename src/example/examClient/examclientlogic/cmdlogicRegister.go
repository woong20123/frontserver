package examclientlogic

import (
	"example/examClient/clientuser"
	"example/share"
	"fmt"
	"net"
	"strings"

	"github.com/nsf/termbox-go"
	"github.com/woong20123/packet"
	"github.com/woong20123/tcpserver"
)

// ContructLogicManager is register a logic for the command
func ContructLogicManager(lm *tcpserver.LogicManager) {

	// S2CPacketCommandLoginUserRes에 대한 처리 작업을 등록합니다.
	lm.RegistLogicfun(share.S2CPacketCommandLoginUserRes, func(conn *net.TCPConn, p *packet.Packet) {

		res := share.S2CPCLoginUserRes{}
		p.Read(&res.Result, &res.UserSn, &res.UserID)

		if res.Result == share.ResultSuccess {
			eu := GetInstance().GetObjMgr().GetUser()
			eu.SetID(res.UserID)
			eu.SetSn(res.UserSn)
			GetInstance().GetObjMgr().GetChanManager().SendChanUserState(clientuser.UserStateEnum.LobbySTATE, "[Login 성공]")
		} else {
			if res.Result == share.ResultExistUserID {
				GetInstance().GetObjMgr().GetChanManager().SendChanUserState(clientuser.UserStateEnum.ConnectedSTATE, fmt.Sprint("[Login 실패] ", res.UserID, "  유저가 이미 존재합니다."))
			} else {
				GetInstance().GetObjMgr().GetChanManager().SendChanUserState(clientuser.UserStateEnum.ConnectedSTATE, "[Login 실패] ")
			}
		}
		return
	})

	// S2CPacketCommandGolobalMsgRes에 대한 처리 작업을 등록합니다.
	lm.RegistLogicfun(share.S2CPacketCommandLobbyMsgRes, func(conn *net.TCPConn, p *packet.Packet) {
		res := share.S2CPCLobbySendMsgRes{}
		p.Read(&res.Result, &res.Userid, &res.Msg)

		var sb strings.Builder
		sb.WriteString(res.Userid)
		sb.WriteString(" : ")
		sb.WriteString(res.Msg)
		GetInstance().GetObjMgr().GetChanManager().SendchanRequestToGui(ToGUIEnum.TYPEMsgPrint, sb.String(), termbox.ColorDefault)
		return
	})

	lm.RunLogicHandle(1)
}
