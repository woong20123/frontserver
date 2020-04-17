package main

import (
	"bufio"
	"context"
	"example/examClient/clientuser"
	"example/examClient/examclientlogic"
	"example/share"
	"log"
	"net"
	"os"
	"os/exec"
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

	sendUserSceneChan(clientuser.UserStateEnum.ConnectedSTATE)

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

func sceneClear() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func loginSceneCommand() {
	println("전체 채팅을 사용하시려면 메시지를 입력하고 enter키를 누르세요")
	println("[명령어]")
	println("\"/RoomEnter <방이름>\" : 방입장")
	println("\"/Close\" : 종료")
}

func sendUserSceneChan(state int) {
	examclientlogic.GetInstance().GetObjMgr().GetChanManager().SendChanUserState(state, 200)
}

func handleUserScene(errProc context.CancelFunc, user *clientuser.ExamUser) {
	defer errProc()

	cobjmgr := examclientlogic.GetInstance().GetObjMgr()
	chanState := cobjmgr.GetChanManager().GetChanUserState()
	log.Println("Input Start")

	// User가 NoneSTATE일때 Scene을 정의합니다.
	user.RegistScene(clientuser.UserStateEnum.NoneSTATE, func(closechan chan int) {
		timer := time.NewTimer(time.Millisecond * 500)
		print("connecting server")

		for {
			select {
			case <-timer.C:
				print(".")
				timer.Reset(time.Millisecond * 200)
			case <-closechan:
				return
			}
		}
	})

	// User가 ConnectedSTATE 일 때 Scene을 정의합니다.
	user.RegistScene(clientuser.UserStateEnum.ConnectedSTATE, func(closechan chan int) {
		println("127.0.0.1:20224", " 서버에 접속 하였습니다.")
		print("접속하려는 ID를 입력해주세요 : ")

		select {
		case <-closechan:
			return
		}
	})

	// User가 LoginSTATE 일 때 Scene을 정의합니다.
	user.RegistScene(clientuser.UserStateEnum.LoginSTATE, func(closechan chan int) {
		println("Login을 성공하였습니다.")
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
		println("[RoomEnterSTATE]")
		for {
			select {
			case <-closechan:
				return
			}
		}
	})

	// User가 CloseSTATE 일 때 Scene을 정의합니다.
	user.RegistScene(clientuser.UserStateEnum.CloseSTATE, func(closechan chan int) {
		println("클라이언트를 종료합니다. 종료 처리 등록")
		user.GetConn().Close()
		time.Sleep(time.Millisecond * 500)
		errProc()
	})

	for {
		sceneClear()
		curState := user.GetState()
		user.RunScene(curState)

		select {
		case nextstate := <-chanState:
			user.SetState(nextstate.State)
			user.CloseScene()
		}
	}
}

func handleInputIO(errProc context.CancelFunc, sendPacketChan chan<- *packet.Packet, user *clientuser.ExamUser) {
	defer errProc()
	reader := bufio.NewReader(os.Stdin)

	for {
		switch user.GetState() {
		case clientuser.UserStateEnum.ConnectedSTATE:
			userid, _ := reader.ReadString('\n')

			// User Login 패킷 전송
			p := packet.NewPacket()
			p.SetHeader(share.ExamplePacketSerialkey, 0, share.C2SPacketCommandLoginUserReq)
			p.WriteString(userid)
			sendPacketChan <- p

		case clientuser.UserStateEnum.LoginSTATE:
			msg, _ := reader.ReadString('\n')

			if strings.Index(msg, "/") == 0 {
				if strings.Contains(msg, "/RoomEnter") {
					strs := strings.Fields(msg)
					if 2 == len(strs) {
						log.Println("방입장 패킷 => ", strs[1])
					} else {
						log.Println("정상적인 명령을 입력해주세요")
					}
				} else if strings.Contains(msg, "/Close") {
					sendUserSceneChan(clientuser.UserStateEnum.CloseSTATE)
				} else if strings.Contains(msg, "/?") {
					sceneClear()
					loginSceneCommand()
				} else {
					log.Println("정상적인 명령을 입력해주세요")
				}
			} else {
				// global msg 패킷 전송
				p := packet.NewPacket()
				p.SetHeader(share.ExamplePacketSerialkey, 0, share.C2SPacketCommandGolobalMsgReq)
				p.WriteString(msg)
				sendPacketChan <- p
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func main() {
	var (
		ip   = "127.0.0.1"
		port = 20224
	)

	ProcCtx, shutdown := context.WithCancel(context.Background())
	sendPacketChan := make(chan *packet.Packet, 1024)

	// set LogicManager
	lm := tcpserver.GetObjInstance().GetLogicManager()
	examclientlogic.ContructLogicManager(lm)

	// make ExamUser
	eu := clientuser.NewExamUser()

	go SocketClient(shutdown, ip, port, sendPacketChan)
	go handleInputIO(shutdown, sendPacketChan, eu)
	go handleUserScene(shutdown, eu)

	select {
	case <-ProcCtx.Done():
		shutdown()
	}
}
