package main

import (
	"context"
	"example/examClient/examclientlogic"
	"example/examshare"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nsf/termbox-go"
	"github.com/woong20123/packet"
	"github.com/woong20123/tcpserver"
)

const (
	maxBufferSize = 4096
)

func handleRead(client *tcpserver.TCPClientSession, errRead context.CancelFunc) {
	defer errRead()
	recvBuf := make([]byte, maxBufferSize)
	AssemblyBuf := make([]byte, maxBufferSize+128)
	var AssemPos uint32 = 0
	var onPacket *packet.Packet = nil

	// session으로부터 전달받은 버퍼를 packet형태로 변환처리하기 위한 Packet
	// TCP의 데이터 전달이 패킷단위로 전달되지 않기 때문에 조립 작업을 합니다.

	for {
		n, err := client.Read(recvBuf)

		if err != nil {
			return
		}

		if 0 < n {
			copylength := copy(AssemblyBuf[AssemPos:], recvBuf[:n])
			AssemPos += uint32(copylength)

			for {
				AssemPos, onPacket = packet.AssemblyFromBuffer(AssemblyBuf, AssemPos, uint32(examshare.Etc_ExamplePacketSerialkey))
				if onPacket == nil {
					break
				}
				tcpserver.Instance().ClientLogicManager().CallLogicFun(onPacket.Command(), client, onPacket)
			}
		}
	}
}

func handleSend(client *tcpserver.TCPClientSession, errSend context.CancelFunc, sendPacketChan <-chan *packet.Packet) {
	defer errSend()

	for {
		// 패킷이 전달되면 패킷을 서버에 전송합니다
		p := <-sendPacketChan
		if client != nil && p != nil {
			n, err := client.Write(p.MakeByte())

			if err != nil {
				return
			}

			// 패킷이 모두 전송되지 않았습니다.
			if p.PacketTotalSize() != uint16(n) {
				log.Println("Packet Write Not PacketSize != WriteReturn ", p.PacketSize(), ":", n)
			}

			packet.Pool().ReleasePacket(p)
		}
	}
}

// HandleNetwork is
func HandleNetwork(errProc context.CancelFunc, sendPacketChan <-chan *packet.Packet) {
	defer errProc()

	chanConnectSrvInfo := examclientlogic.Instance().ObjMgr().ChanManager().ChanSrvInfo()
	srvInfo := <-chanConnectSrvInfo

	chatClient := examclientlogic.Instance().ObjMgr().ChatClient()
	err := chatClient.Connect(srvInfo.Ip, srvInfo.Port)

	if err != nil {
		log.Println("DialTCP", err)
		os.Exit(1)
	}

	sendUserSceneChan(examclientlogic.UserStateEnum.ConnectedSTATE, []string{fmt.Sprint("==========================", "[ 접 속 화 면 ]", "=========================="), fmt.Sprint("서버 접속 성공 [접속정보]", srvInfo.Ip, ":", srvInfo.Port, " ")})

	defer chatClient.Close()

	readCtx, errRead := context.WithCancel(context.Background())
	sendCtx, errSend := context.WithCancel(context.Background())

	go handleRead(chatClient, errRead)
	go handleSend(chatClient, errSend, sendPacketChan)

	select {
	case <-readCtx.Done():
	case <-sendCtx.Done():
	}
}

func readsceneClear() {
	examclientlogic.Instance().ObjMgr().ChanManager().SendchanRequestToGui(examclientlogic.ToGUIEnum.TYPEWindowClear, "", termbox.ColorDefault)
}

func readSceneSystemWrite(msg string) {
	examclientlogic.Instance().ObjMgr().ChanManager().SendchanRequestToGui(examclientlogic.ToGUIEnum.TYPEMsgPrint, msg, termbox.ColorYellow)
}

func readSceneErrorWrite(msg string) {
	examclientlogic.Instance().ObjMgr().ChanManager().SendchanRequestToGui(examclientlogic.ToGUIEnum.TYPEMsgPrint, msg, termbox.ColorRed)
}

func readSceneWrite(msg string, colword termbox.Attribute) {
	examclientlogic.Instance().ObjMgr().ChanManager().SendchanRequestToGui(examclientlogic.ToGUIEnum.TYPEMsgPrint, msg, colword)
}

func noneStateSceneCommand() {
	readSceneSystemWrite(fmt.Sprint("[프로세스정보] PID = ", os.Getpid()))
	readSceneSystemWrite("[서버주소:포트번호] 형태로 키를 입력하시면 서버에 접속합니다.")
	readSceneSystemWrite("입력없이 엔터키를 누르시면 기본 서버주소로 접속합니다.")
	readSceneSystemWrite("[기본주소 정보] 127.0.0.1:20224")
	readSceneSystemWrite("-----------------------------------------------------------")
}

func lobbySceneCommand() {
	readSceneSystemWrite(fmt.Sprint("[유저정보] ID = ", examclientlogic.Instance().ObjMgr().User().ID()))
	readSceneSystemWrite("전체 채팅을 사용하시려면 메시지를 입력하고 enter키를 누르세요")
	readSceneSystemWrite("[명령어]")
	readSceneSystemWrite("\"/?\" : 화면 클리어 및 명령어 재출력")
	readSceneSystemWrite("\"/RoomCreate <방이름>\" : 방 생성 요청")
	readSceneSystemWrite("\"/RoomEnter <방이름>\" : 방 입장 요청")
	//readSceneSystemWrite("\"/RoomList\" : 현재 생성된 방 목록")
	readSceneSystemWrite("\"/Close\" : 종료")
	readSceneSystemWrite("-----------------------------------------------------------")
}

func roomSceneCommand() {
	readSceneSystemWrite(fmt.Sprint("[유저정보] ID = ", examclientlogic.Instance().ObjMgr().User().ID()))
	readSceneSystemWrite(fmt.Sprint("[Room이름] = ", examclientlogic.Instance().ObjMgr().User().RoomName()))
	readSceneSystemWrite("채팅을 사용하시려면 메시지를 입력하고 enter키를 누르세요")
	readSceneSystemWrite("[명령어]")
	readSceneSystemWrite("\"/?\" : 화면 클리어 및 명령어 재출력")
	readSceneSystemWrite("\"/RoomLeave\" : 방을 나갑니다.")
	readSceneSystemWrite("\"/Close\" : 종료")
	readSceneSystemWrite("-----------------------------------------------------------")
}

func sendUserSceneChan(state int, msgs []string) {
	examclientlogic.Instance().ObjMgr().ChanManager().SendChanUserState(state, msgs)
}

func handleUserScene(errProc context.CancelFunc) {
	defer errProc()

	cobjmgr := examclientlogic.Instance().ObjMgr()
	chanState := cobjmgr.ChanManager().ChanUserState()
	user := cobjmgr.User()

	// User가 NoneSTATE일때 Scene을 정의합니다.
	user.RegistScene(examclientlogic.UserStateEnum.NoneSTATE, func(closechan chan int) {
		noneStateSceneCommand()

		select {
		case <-closechan:
			return
		}
	})

	// User가 ConnectedSTATE 일 때 Scene을 정의합니다.
	user.RegistScene(examclientlogic.UserStateEnum.ConnectedSTATE, func(closechan chan int) {
		readSceneSystemWrite("접속시 사용할 ID를 입력해주세요")

		select {
		case <-closechan:
			return
		}
	})

	// User가 LoginSTATE 일 때 Scene을 정의합니다.
	user.RegistScene(examclientlogic.UserStateEnum.LobbySTATE, func(closechan chan int) {
		lobbySceneCommand()
		for {
			select {
			case <-closechan:
				return
			}
		}
	})

	// User가 RoomEnterSTATE 일 때 Scene을 정의합니다.
	user.RegistScene(examclientlogic.UserStateEnum.RoomSTATE, func(closechan chan int) {
		roomSceneCommand()
		for {
			select {
			case <-closechan:
				return
			}
		}
	})

	// User가 CloseSTATE 일 때 Scene을 정의합니다.
	user.RegistScene(examclientlogic.UserStateEnum.CloseSTATE, func(closechan chan int) {
		readSceneErrorWrite("클라이언트를 종료합니다. 종료 처리 등록")
		user.Conn().Close()
		time.Sleep(time.Millisecond * 500)
		errProc()
	})

	for {
		curState := user.State()
		user.RunScene(curState)

		select {
		case nextstate := <-chanState:
			user.SetState(nextstate.State)
			readsceneClear()
			user.CloseScene()
			for _, msg := range nextstate.Msgs {
				readSceneSystemWrite(msg)
			}

		}
	}
}

func handleScene(errProc context.CancelFunc, sendPacketChan chan<- *packet.Packet) {
	defer errProc()

	cobjmgr := examclientlogic.Instance().ObjMgr()
	user := cobjmgr.User()
	chanRequestFromGui := cobjmgr.ChanManager().ChanRequestFromGui()

	for {
		select {
		case requestFromGui := <-chanRequestFromGui:
			msg := requestFromGui.Msg
			msg = strings.Trim(msg, "\n")

			switch user.State() {
			case examclientlogic.UserStateEnum.NoneSTATE:
				serverinfo := msg
				if serverinfo == "" {
					// 기본서버로 접속합니다.
					cobjmgr.ChanManager().SendChanSrvInfo(baseIP, basePort)
					readSceneSystemWrite(fmt.Sprint(baseIP, ":", basePort, "서버에 접속을 요청합니다."))
				} else {
					infos := strings.Split(serverinfo, ":")
					// 유저가 입력한 서버로 접속합니다.
					if len(infos) == 2 {
						port, err := strconv.Atoi(infos[1])
						if err == nil {
							cobjmgr.ChanManager().SendChanSrvInfo(infos[0], port)
							readSceneSystemWrite(fmt.Sprint(infos[0], ":", port, "서버에 접속을 요청합니다."))
						}
					} else {
						readSceneErrorWrite("입력한 서버 정보가 비정상적입니다.")
					}
				}

			case examclientlogic.UserStateEnum.ConnectedSTATE:
				userid := msg
				if userid != "" {
					// User Login 패킷 전송
					p := packet.Pool().AcquirePacket()
					p.SetHeader(uint32(examshare.Etc_ExamplePacketSerialkey), 0, int32(examshare.Cmd_C2SLoginUserReq))
					req := examshare.C2CS_LoginUserReq{}
					req.UserID = userid
					err := p.MarshalFromProto(&req)
					if err == nil {
						sendPacketChan <- p
					} else {
						readSceneErrorWrite("유저의 ID가 비정상적입니다.")
						packet.Pool().ReleasePacket(p)
					}

				} else {
					readSceneErrorWrite("유저의 ID가 빈문자열입니다.")
				}

			case examclientlogic.UserStateEnum.LobbySTATE:
				if strings.Index(msg, "/") == 0 {
					if strings.Contains(msg, "/RoomEnter ") {
						readSceneErrorWrite("/RoomEnter 입력")
						fileds := strings.Fields(msg)
						if len(fileds) != 2 {
							readSceneErrorWrite(fmt.Sprint("len(fileds) = ", len(fileds)))
							continue
						}

						roomName := fileds[1]

						p := packet.Pool().AcquirePacket()
						p.SetHeaderByDefaultKey(0, int32(examshare.Cmd_C2SRoomEnterReq))
						req := examshare.C2CS_RoomEnterReq{}
						req.RoomName = roomName
						err := p.MarshalFromProto(&req)
						if err == nil {
							sendPacketChan <- p
							readSceneErrorWrite(fmt.Sprint("Send Packet Cmd_C2SRoomEnterReq name = ", roomName))
						} else {
							packet.Pool().ReleasePacket(p)
						}

					} else if strings.Contains(msg, "/RoomCreate") {
						fileds := strings.Fields(msg)
						if len(fileds) != 2 {
							readSceneErrorWrite(fmt.Sprint("len(fileds) = ", len(fileds)))
							continue
						}
						roomName := fileds[1]
						p := packet.Pool().AcquirePacket()
						p.SetHeaderByDefaultKey(0, int32(examshare.Cmd_C2SRoomCreateReq))
						req := examshare.C2CS_RoomEnterReq{}
						req.RoomName = roomName
						err := p.MarshalFromProto(&req)
						if err == nil {
							sendPacketChan <- p
							readSceneErrorWrite(fmt.Sprint("Send Packet Cmd_C2SRoomCreateReq name = ", roomName))
						} else {
							packet.Pool().ReleasePacket(p)
						}

					} else if strings.Contains(msg, "/RoomList") {
						readsceneClear()
						lobbySceneCommand()
					} else if strings.Contains(msg, "/Close") {

						sendUserSceneChan(examclientlogic.UserStateEnum.CloseSTATE, []string{})
					} else if strings.Contains(msg, "/?") {
						readsceneClear()
						lobbySceneCommand()
					} else {
						readSceneErrorWrite("정상적인 명령을 입력해주세요")
					}
				} else {
					if msg != "" {
						// lobby msg 패킷 전송
						p := packet.Pool().AcquirePacket()
						p.SetHeaderByDefaultKey(0, int32(examshare.Cmd_C2SLobbyMsgReq))
						req := examshare.C2CS_LobbySendMsgReq{}
						req.Msg = msg
						err := p.MarshalFromProto(&req)
						if err == nil {
							sendPacketChan <- p
						} else {
							packet.Pool().ReleasePacket(p)
						}
					}
				}
			case examclientlogic.UserStateEnum.RoomSTATE:
				if strings.Index(msg, "/") == 0 {
					if strings.Contains(msg, "/RoomLeave") {
						// global msg 패킷 전송
						p := packet.Pool().AcquirePacket()
						p.SetHeaderByDefaultKey(0, int32(examshare.Cmd_C2SRoomLeaveReq))
						sendPacketChan <- p
					} else if strings.Contains(msg, "/?") {
						readsceneClear()
						roomSceneCommand()
					} else if strings.Contains(msg, "/Close") {
						sendUserSceneChan(examclientlogic.UserStateEnum.CloseSTATE, []string{})
					}
				} else {
					if msg != "" {
						// room msg 패킷 전송
						p := packet.Pool().AcquirePacket()
						p.SetHeaderByDefaultKey(0, int32(examshare.Cmd_C2SRoomMsgReq))
						req := examshare.C2CS_RoomSendMsgReq{}
						req.RoomIdx = user.RoomIdx()
						req.Msg = msg
						err := p.MarshalFromProto(&req)
						if err == nil {
							sendPacketChan <- p
						}
					}
				}
			}
		}
	}
}

const (
	baseIP   = "127.0.0.1"
	basePort = 20224
)

func main() {
	GuiInitchan := make(chan int)
	go examclientlogic.RunGui(GuiInitchan)
	<-GuiInitchan

	packet.RegistSerialKey(uint32(examshare.Etc_ExamplePacketSerialkey))

	ProcCtx, shutdown := context.WithCancel(context.Background())
	chanSendPacket := make(chan *packet.Packet, 1024)

	// set LogicManager
	examclientlogic.ContructLogicManager()

	go HandleNetwork(shutdown, chanSendPacket)
	go handleScene(shutdown, chanSendPacket)
	go handleUserScene(shutdown)

	select {
	case <-ProcCtx.Done():
		shutdown()
	}
}
