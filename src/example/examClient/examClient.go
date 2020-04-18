package main

import (
	"context"
	"example/examClient/clientuser"
	"example/examClient/examclientgui"
	"example/examClient/examclientlogic"
	"example/share"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

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
				tcpserver.GetObjInstance().GetLogicManager().CallLogicFun(onPacket.GetCommand(), conn, onPacket)
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
func SocketClient(errProc context.CancelFunc, ip string, port int, sendPacketChan <-chan *packet.Packet) {
	defer errProc()

	var remoteaddr net.TCPAddr
	remoteaddr.IP = net.ParseIP(ip)
	remoteaddr.Port = port
	conn, err := net.DialTCP("tcp", nil, &remoteaddr)

	if err != nil {
		log.Println("DialTCP", err)
		os.Exit(1)
	}

	sendUserSceneChan(clientuser.UserStateEnum.ConnectedSTATE, fmt.Sprintln("서버에 접속을 성공하였습니다.\n[접속정보]", ip, ":", port, " "))

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
	examclientlogic.GetInstance().GetObjMgr().GetChanManager().SendchanRequestToGui(examclientlogic.ToGUIEnum.TYPEWindowClear, "")
}

func readSceneWrite(msg string) {
	examclientlogic.GetInstance().GetObjMgr().GetChanManager().SendchanRequestToGui(examclientlogic.ToGUIEnum.TYPEMsgPrint, msg)
}

func loginSceneCommand() {
	readSceneWrite("전체 채팅을 사용하시려면 메시지를 입력하고 enter키를 누르세요")
	readSceneWrite("[명령어]")
	readSceneWrite("\"/RoomEnter <방이름>\" : 방입장")
	readSceneWrite("\"/Close\" : 종료")
}

func sendUserSceneChan(state int, msg string) {
	examclientlogic.GetInstance().GetObjMgr().GetChanManager().SendChanUserState(state, msg)
}

func handleUserScene(errProc context.CancelFunc) {
	defer errProc()

	cobjmgr := examclientlogic.GetInstance().GetObjMgr()
	chanState := cobjmgr.GetChanManager().GetChanUserState()
	user := cobjmgr.GetUser()
	readSceneWrite("Input Start")

	// User가 NoneSTATE일때 Scene을 정의합니다.
	user.RegistScene(clientuser.UserStateEnum.NoneSTATE, func(closechan chan int) {
		timer := time.NewTimer(time.Millisecond * 500)
		readSceneWrite("connecting server")

		for {
			select {
			case <-timer.C:
				readSceneWrite(".")
				timer.Reset(time.Millisecond * 200)
			case <-closechan:
				return
			}
		}
	})

	// User가 ConnectedSTATE 일 때 Scene을 정의합니다.
	user.RegistScene(clientuser.UserStateEnum.ConnectedSTATE, func(closechan chan int) {
		readSceneWrite("접속하려는 ID를 입력해주세요 : ")

		select {
		case <-closechan:
			return
		}
	})

	// User가 LoginSTATE 일 때 Scene을 정의합니다.
	user.RegistScene(clientuser.UserStateEnum.LoginSTATE, func(closechan chan int) {
		loginSceneCommand()
		for {
			select {
			case <-closechan:
				return
			}
		}
	})

	// User가 RoomEnterSTATE 일 때 Scene을 정의합니다.
	user.RegistScene(clientuser.UserStateEnum.RoomEnterSTATE, func(closechan chan int) {
		readSceneWrite("[RoomEnterSTATE]")
		for {
			select {
			case <-closechan:
				return
			}
		}
	})

	// User가 CloseSTATE 일 때 Scene을 정의합니다.
	user.RegistScene(clientuser.UserStateEnum.CloseSTATE, func(closechan chan int) {
		readSceneWrite("클라이언트를 종료합니다. 종료 처리 등록")
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
			readSceneWrite(nextstate.Msg)
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
			case clientuser.UserStateEnum.ConnectedSTATE:
				userid := requestFromGui.Msg
				userid = strings.Trim(userid, "\n")

				// User Login 패킷 전송
				p := packet.NewPacket()
				p.SetHeader(share.ExamplePacketSerialkey, 0, share.C2SPacketCommandLoginUserReq)
				p.WriteString(userid)
				sendPacketChan <- p
			case clientuser.UserStateEnum.LoginSTATE:
				msg := requestFromGui.Msg
				msg = strings.Trim(msg, "\n")

				if strings.Index(msg, "/") == 0 {
					if strings.Contains(msg, "/RoomEnter") {
						// strs := strings.Fields(msg)
						// if 2 == len(strs) {
						// 	readSceneWrite("방입장 패킷 => ", strs[1])
						// } else {
						// 	readSceneWrite("정상적인 명령을 입력해주세요")
						//}
					} else if strings.Contains(msg, "/Close") {
						sendUserSceneChan(clientuser.UserStateEnum.CloseSTATE, "클라이언트를 종료합니다. Bye")
					} else if strings.Contains(msg, "/?") {
						readsceneClear()
						loginSceneCommand()
					} else {
						readSceneWrite("정상적인 명령을 입력해주세요")
					}
				} else {
					// global msg 패킷 전송
					p := packet.NewPacket()
					p.SetHeader(share.ExamplePacketSerialkey, 0, share.C2SPacketCommandGolobalMsgReq)
					p.WriteString(msg)
					sendPacketChan <- p
				}
			}
		}
	}
}

func main() {
	var (
		ip   = "127.0.0.1"
		port = 20224
	)

	GuiInitchan := make(chan int)
	go examclientgui.RunGui(GuiInitchan)
	<-GuiInitchan

	ProcCtx, shutdown := context.WithCancel(context.Background())
	sendPacketChan := make(chan *packet.Packet, 1024)

	// set LogicManager
	lm := tcpserver.GetObjInstance().GetLogicManager()
	examclientlogic.ContructLogicManager(lm)

	go SocketClient(shutdown, ip, port, sendPacketChan)
	go handleScene(shutdown, sendPacketChan)
	go handleUserScene(shutdown)

	select {
	case <-ProcCtx.Done():
		shutdown()
	}
}
