package main

import (
	"bufio"
	"context"
	"encoding/binary"
	"example/examClient/clientobjmanager"
	"example/examClient/clientuser"
	"example/share"
	"log"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/woong20123/logicmanager"

	"github.com/woong20123/packet"
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
			log.Println(recvBuf[:n])
			for {
				AssemPos, onPacket = packet.AssemblyFromBuffer(AssemblyBuf, AssemPos, share.ExamplePacketSerialkey)
				if onPacket == nil {
					break
				}
				//lm.CallLogicFun(onPacket.GetCommand(), conn, onPacket)
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
		log.Println("On packet send")
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

func sendUserSceneChan(state int) {
	clientobjmanager.GetInstance().GetChanManager().GetChanUserState() <- state
}

func handleUserScene(errProc context.CancelFunc, user *clientuser.ExamUser) {
	defer errProc()

	cobjmgr := clientobjmanager.GetInstance()
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
		return
	})

	// User가 LoginSTATE 일 때 Scene을 정의합니다.
	user.RegistScene(clientuser.UserStateEnum.LoginSTATE, func(closechan chan int) {
		println("[LoginSTATE]")
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

	for {
		sceneClear()
		curState := user.GetState()
		user.RunScene(curState)

		select {
		case nextstate := <-chanState:
			user.SetState(nextstate)
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
			p.WriteString(binary.LittleEndian, userid)
			sendPacketChan <- p

		case clientuser.UserStateEnum.LoginSTATE:
			msg, _ := reader.ReadString('\n')

			// global msg 패킷 전송
			p := packet.NewPacket()
			p.SetHeader(share.ExamplePacketSerialkey, 0, share.C2SPacketCommandGolobalSendMsgReq)
			p.Write(binary.LittleEndian, share.C2SPCGolobalSendMsgReq{msg})
			sendPacketChan <- p

		}
		time.Sleep(1 * time.Millisecond)
	}
}

// ContructLogicManager is
func ContructLogicManager(lm *logicmanager.LogicManager) {
	lm.RegistLogicfun(share.S2CPacketCommandLoginUserRes, func(conn *net.TCPConn, p *packet.Packet) {
		log.Println("S2CPacketCommandLoginUserRes")
		// 패킷을 확인하고 체크합니다.
		sendUserSceneChan(clientuser.UserStateEnum.LoginSTATE)
		return
	})

	lm.RegistLogicfun(share.S2CPacketCommandGolobalSendMsgRes, func(conn *net.TCPConn, p *packet.Packet) {
		log.Println("S2CPacketCommandGolobalSendMsgRes")
		return
	})
}

func main() {
	var (
		ip   = "127.0.0.1"
		port = 20224
	)

	clientobjmanager.GetInstance()

	ProcCtx, shutdown := context.WithCancel(context.Background())
	sendPacketChan := make(chan *packet.Packet, 1024)

	// set LogicManager
	lm := logicmanager.NewLogicManager()
	ContructLogicManager(lm)

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
