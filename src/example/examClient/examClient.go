package main

import (
	"context"
	"example/examClient/clientuser"
	"example/examClient/examclientlogic"
	"example/share"
	"fmt"
	"log"
	"net"
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

func handleRead(conn *net.TCPConn, errRead context.CancelFunc) {
	defer errRead()
	recvBuf := make([]byte, maxBufferSize)
	AssemblyBuf := make([]byte, maxBufferSize+128)
	var AssemPos uint32 = 0
	var onPacket *packet.Packet = nil

	// session으로부터 전달받은 버퍼를 packet형태로 변환처리하기 위한 Packet
	// TCP의 데이터 전달이 패킷단위로 전달되지 않기 때문에 조립 작업을 합니다.

	for {
		n, err := conn.Read(recvBuf)
		if err != nil {
			if ne, ok := err.(net.Error); ok {
				switch {
				case ne.Temporary():
					continue
				}
			}

			log.Println("Read", err)
			return
		}
		if err != nil {
			log.Println("Write", err)
			return
		}

		if 0 < n {
			copylength := copy(AssemblyBuf[AssemPos:], recvBuf[:n])
			AssemPos += uint32(copylength)

			for {
				AssemPos, onPacket = packet.AssemblyFromBuffer(AssemblyBuf, AssemPos, share.ExamplePacketSerialkey)
				if onPacket == nil {
					break
				}
				tcpserver.GetInstance().GetLogicManager().CallLogicFun(onPacket.GetCommand(), conn, onPacket)
			}
		}
	}
}

func handleSend(conn *net.TCPConn, errSend context.CancelFunc, sendPacketChan <-chan *packet.Packet) {
	defer errSend()

	for {
		// 패킷이 전달되면 패킷을 서버에 전송합니다
		p := <-sendPacketChan
		conn.Write(p.GetByte())
	}
}

// SocketClient is
func HandleNetwork(errProc context.CancelFunc, sendPacketChan <-chan *packet.Packet) {
	defer errProc()

	chanConnectSrvInfo := examclientlogic.GetInstance().GetObjMgr().GetChanManager().GetChanSrvInfo()
	srvInfo := <-chanConnectSrvInfo

	var remoteaddr net.TCPAddr
	remoteaddr.IP = net.ParseIP(srvInfo.Ip)
	remoteaddr.Port = srvInfo.Port
	conn, err := net.DialTCP("tcp", nil, &remoteaddr)

	if err != nil {
		log.Println("DialTCP", err)
		os.Exit(1)
	}

	sendUserSceneChan(clientuser.UserStateEnum.ConnectedSTATE, fmt.Sprint("서버 접속 성공 [접속정보]", srvInfo.Ip, ":", srvInfo.Port, " "))

	defer conn.Close()

	readCtx, errRead := context.WithCancel(context.Background())
	sendCtx, errSend := context.WithCancel(context.Background())

	go handleRead(conn, errRead)
	go handleSend(conn, errSend, sendPacketChan)

	select {
	case <-readCtx.Done():
	case <-sendCtx.Done():
	}
}

func readsceneClear() {
	examclientlogic.GetInstance().GetObjMgr().GetChanManager().SendchanRequestToGui(examclientlogic.ToGUIEnum.TYPEWindowClear, "", termbox.ColorDefault)
}

func readSceneSystemWrite(msg string) {
	examclientlogic.GetInstance().GetObjMgr().GetChanManager().SendchanRequestToGui(examclientlogic.ToGUIEnum.TYPEMsgPrint, msg, termbox.ColorYellow)
}

func readSceneErrorWrite(msg string) {
	examclientlogic.GetInstance().GetObjMgr().GetChanManager().SendchanRequestToGui(examclientlogic.ToGUIEnum.TYPEMsgPrint, msg, termbox.ColorRed)
}

func readSceneWrite(msg string, colword termbox.Attribute) {
	examclientlogic.GetInstance().GetObjMgr().GetChanManager().SendchanRequestToGui(examclientlogic.ToGUIEnum.TYPEMsgPrint, msg, colword)
}

func noneStateSceneCommand() {
	readSceneSystemWrite("[서버주소:포트번호] 형태로 키를 입력하시면 서버에 접속합니다.")
	readSceneSystemWrite("입력없이 엔터키를 누르시면 기본 서버주소로 접속합니다.")
	readSceneSystemWrite("[기본주소 정보] 127.0.0.1:20224")
	readSceneSystemWrite("-----------------------------------------------------------")
}

func lobbySceneCommand() {
	readSceneSystemWrite(fmt.Sprint("[유저정보] ID = ", examclientlogic.GetInstance().GetObjMgr().GetUser().GetID()))
	readSceneSystemWrite("전체 채팅을 사용하시려면 메시지를 입력하고 enter키를 누르세요")
	readSceneSystemWrite("[명령어]")
	readSceneSystemWrite("\"/?\" : 화면 클리어 및 명령어 재출력")
	readSceneSystemWrite("\"/RoomEnter <방이름>\" : 방입장")
	readSceneSystemWrite("\"/RoomList\" : 현재 생성된 방 목록")
	readSceneSystemWrite("\"/Close\" : 종료")
	readSceneSystemWrite("-----------------------------------------------------------")
}

func sendUserSceneChan(state int, msg string) {
	examclientlogic.GetInstance().GetObjMgr().GetChanManager().SendChanUserState(state, msg)
}

func handleUserScene(errProc context.CancelFunc) {
	defer errProc()

	cobjmgr := examclientlogic.GetInstance().GetObjMgr()
	chanState := cobjmgr.GetChanManager().GetChanUserState()
	user := cobjmgr.GetUser()

	// User가 NoneSTATE일때 Scene을 정의합니다.
	user.RegistScene(clientuser.UserStateEnum.NoneSTATE, func(closechan chan int) {
		noneStateSceneCommand()

		select {
		case <-closechan:
			return
		}
	})

	// User가 ConnectedSTATE 일 때 Scene을 정의합니다.
	user.RegistScene(clientuser.UserStateEnum.ConnectedSTATE, func(closechan chan int) {
		readSceneSystemWrite("접속시 사용할 ID를 입력해주세요")

		select {
		case <-closechan:
			return
		}
	})

	// User가 LoginSTATE 일 때 Scene을 정의합니다.
	user.RegistScene(clientuser.UserStateEnum.LobbySTATE, func(closechan chan int) {
		lobbySceneCommand()
		for {
			select {
			case <-closechan:
				return
			}
		}
	})

	// User가 RoomEnterSTATE 일 때 Scene을 정의합니다.
	user.RegistScene(clientuser.UserStateEnum.RoomEnterSTATE, func(closechan chan int) {
		for {
			select {
			case <-closechan:
				return
			}
		}
	})

	// User가 CloseSTATE 일 때 Scene을 정의합니다.
	user.RegistScene(clientuser.UserStateEnum.CloseSTATE, func(closechan chan int) {
		readSceneErrorWrite("클라이언트를 종료합니다. 종료 처리 등록")
		user.GetConn().Close()
		time.Sleep(time.Millisecond * 500)
		errProc()
	})

	for {
		curState := user.GetState()
		user.RunScene(curState)

		select {
		case nextstate := <-chanState:
			user.SetState(nextstate.State)
			readsceneClear()
			user.CloseScene()
			readSceneSystemWrite(nextstate.Msg)
		}
	}
}

func handleScene(errProc context.CancelFunc, sendPacketChan chan<- *packet.Packet) {
	defer errProc()

	cobjmgr := examclientlogic.GetInstance().GetObjMgr()
	user := cobjmgr.GetUser()
	chanRequestFromGui := cobjmgr.GetChanManager().GetchanRequestFromGui()

	for {
		select {
		case requestFromGui := <-chanRequestFromGui:
			switch user.GetState() {
			case clientuser.UserStateEnum.NoneSTATE:
				serverinfo := requestFromGui.Msg
				serverinfo = strings.Trim(serverinfo, "\n")
				if serverinfo == "" {
					// 기본서버로 접속합니다.
					cobjmgr.GetChanManager().SendChanSrvInfo(baseIp, basePort)
					readSceneSystemWrite(fmt.Sprint(baseIp, ":", basePort, "서버에 접속을 요청합니다."))
				} else {
					infos := strings.Split(serverinfo, ":")
					// 유저가 입력한 서버로 접속합니다.
					if len(infos) == 2 {
						port, err := strconv.Atoi(infos[1])
						if err == nil {
							cobjmgr.GetChanManager().SendChanSrvInfo(infos[0], port)
							readSceneSystemWrite(fmt.Sprint(infos[0], ":", port, "서버에 접속을 요청합니다."))
						}
					} else {
						readSceneErrorWrite("입력한 서버 정보가 비정상적입니다.")
					}
				}

			case clientuser.UserStateEnum.ConnectedSTATE:
				userid := requestFromGui.Msg
				userid = strings.Trim(userid, "\n")

				// User Login 패킷 전송
				p := packet.NewPacket()
				p.SetHeader(share.ExamplePacketSerialkey, 0, share.C2SPacketCommandLoginUserReq)
				p.Write(&userid)
				sendPacketChan <- p
			case clientuser.UserStateEnum.LobbySTATE:
				msg := requestFromGui.Msg
				msg = strings.Trim(msg, "\n")

				if strings.Index(msg, "/") == 0 {
					if strings.Contains(msg, "/RoomEnter") {
						readsceneClear()
						lobbySceneCommand()
					} else if strings.Contains(msg, "/RoomList") {
						readsceneClear()
						lobbySceneCommand()
					} else if strings.Contains(msg, "/Close") {
						sendUserSceneChan(clientuser.UserStateEnum.CloseSTATE, "클라이언트를 종료합니다. Bye")
					} else if strings.Contains(msg, "/?") {
						readsceneClear()
						lobbySceneCommand()
					} else {
						readSceneErrorWrite("정상적인 명령을 입력해주세요")
					}
				} else {
					// global msg 패킷 전송
					p := packet.NewPacket()
					p.SetHeader(share.ExamplePacketSerialkey, 0, share.C2SPacketCommandGolobalMsgReq)
					p.Write(&msg)
					sendPacketChan <- p
				}
			}
		}
	}
}

const (
	baseIp   = "127.0.0.1"
	basePort = 20224
)

func main() {
	GuiInitchan := make(chan int)
	go examclientlogic.RunGui(GuiInitchan)
	<-GuiInitchan

	ProcCtx, shutdown := context.WithCancel(context.Background())
	chanSendPacket := make(chan *packet.Packet, 1024)

	// set LogicManager
	lm := tcpserver.GetInstance().GetLogicManager()
	examclientlogic.ContructLogicManager(lm)

	go HandleNetwork(shutdown, chanSendPacket)
	go handleScene(shutdown, chanSendPacket)
	go handleUserScene(shutdown)

	select {
	case <-ProcCtx.Done():
		shutdown()
	}
}
