package main

// "factored" import statment
import (
	"context"
	"example/examFrontServer/examserverlogic"
	"example/examshare"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/woong20123/packet"
	"github.com/woong20123/tcpserver"
)

func construct() bool {
	bSuccess := true
	bSuccess = bSuccess && constructConfig()
	bSuccess = bSuccess && constructTCPSession()
	bSuccess = bSuccess && constructLogic()
	bSuccess = bSuccess && constructFrontMode()
	return bSuccess
}

func constructConfig() bool {
	examserverlogic.MakeExampleConfig()
	examserverlogic.Instance().ConfigMgr().ReadConfig("\\ExampleServerConfig.json")
	serverConfig := examserverlogic.Instance().ConfigMgr().ServerConfig()

	// 인자값이 있다면 우선합니다.
	if len(os.Args) >= 3 {
		serverPort, err := strconv.Atoi(os.Args[1])
		if err != nil {
			examserverlogic.Logger().Println("Args 변환 오류 ", os.Args[1])
		}
		serverMode := os.Args[2]

		serverConfig.ServerPort = serverPort
		serverConfig.ServerMode = serverMode
	}

	println("[Server 정보] Port : ", serverConfig.ServerPort, ", Mode : ", serverConfig.ServerMode)
	examserverlogic.Logger().Println("[Server 정보] Port : ", serverConfig.ServerPort, ", Mode : ", serverConfig.ServerMode)

	return true
}

func constructTCPSession() bool {
	// client 세션에 대한 처리 로직 등록하는 곳
	clientsessionhandler := tcpserver.Instance().ClientSessionHandler()
	examserverlogic.RegistClientSessionLogic(clientsessionhandler)

	// server proxy 세션에 대한 처리 로직 등록하는 곳
	serverproxysessionhandler := tcpserver.Instance().ServerProxySessionHandler()
	examserverlogic.RegistServerProxySessionLogic(serverproxysessionhandler)

	return true
}

func constructLogic() bool {
	// regist exam Logic
	lm := tcpserver.Instance().ClientLogicManager()

	switch examserverlogic.Instance().ConfigMgr().ServerConfig().ServerMode {
	case "front":
		examserverlogic.FrontModeRegistCommandLogic(lm)
		slm := tcpserver.Instance().ServerLogicManager()
		examserverlogic.RegistServerLogic(slm)
		slm.RunLogicHandler(1)
	default:
		examserverlogic.ChatServerModeRegistCommandLogic(lm)
	}

	// set logic goroutines count
	lm.RunLogicHandler(1)
	return true
}

func constructFrontMode() bool {
	srvConfig := examserverlogic.Instance().ConfigMgr().ServerConfig()

	// front 모드인 경우에는 상위 서버로의 연결 작업을 진행합니다.
	if "front" == srvConfig.ServerMode {

		// Chat Server Proxy 생성 및 연결  작업
		ChatserverIP := srvConfig.BackEndChatServerIP
		ChatserverPort := srvConfig.BackEndChatServerPort

		ChatSeverProxy := tcpserver.NewTCPClientSession()
		ChatSeverProxy.SetIndex(examshare.TCPCliToSvrIdxChat)
		err := ChatSeverProxy.Connect(ChatserverIP, ChatserverPort)
		if err != nil {
			println("[ChatServer] Connect fail ", ChatserverIP, ":", ChatserverPort)
			examserverlogic.Logger().Println("[ChatServer] Connect fail ", ChatserverIP, ":", ChatserverPort)
			examserverlogic.Logger().Println(err.Error())
			return false
		}
		tcpserver.Instance().ServerProxySessionHandler().RunConnectFunc(tcpserver.SessionStateEnum.OnConnected, ChatSeverProxy)
		err = tcpserver.Instance().TCPClientMgr().AddTCPClientSession(ChatSeverProxy)

		if err != nil {
			println("[ChatServer] Connect fail ", ChatserverIP, ":", ChatserverPort)
			examserverlogic.Logger().Println("[ChatServer] Connect fail ", ChatserverIP, ":", ChatserverPort)
			examserverlogic.Logger().Println(err.Error())
			return false
		}

		println("[ChatServer] Connect success ", ChatserverIP, ":", ChatserverPort)

		// 서버 등록 패킷 전송
		p := packet.Pool().AcquirePacket()
		p.SetHeaderByDefaultKey(0, int32(examshare.Cmd_F2CSServerRegistReq))
		req := examshare.F2CS_ServerRegistReq{}
		req.Ip = srvConfig.ServerIP
		req.Port = int32(srvConfig.ServerPort)
		req.Servertype = examshare.SrvType_FrontSrvMode
		err = p.MarshalFromProto(&req)
		if err == nil {
			tcpserver.Instance().SendManager().SendToServerConn(examshare.TCPCliToSvrIdxChat, p)
		} else {
			packet.Pool().ReleasePacket(p)
		}
	}
	return true
}

func constructTCPServer(ctxServer context.Context, cancel context.CancelFunc, chClosed chan struct{}) (wg sync.WaitGroup) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// server consturct
	packet.RegistSerialKey(uint32(examshare.Etc_ExamplePacketSerialkey))
	tcpserver.Consturct(uint32(examshare.Etc_ExamplePacketSerialkey), runtime.NumCPU())

	// start server handler
	serverConfig := examserverlogic.Instance().ConfigMgr().ServerConfig()
	port := serverConfig.ServerPort
	address := ":" + strconv.Itoa(port)
	go tcpserver.HandleListener(ctxServer, address, &wg, chClosed)
	return
}

func main() {

	sigChan := make(chan os.Signal, 1)
	signal.Ignore()
	signal.Notify(sigChan, syscall.SIGINT)

	chClosed := make(chan struct{})
	examserverlogic.Logger().Println(fmt.Sprint("[Server build Ver ", uint32(examshare.Etc_BuildVer), "]"))
	println(fmt.Sprint("[Server build Ver ", uint32(examshare.Etc_BuildVer), "]"))

	serverCtx, shutdown := context.WithCancel(context.Background())

	bSuccess := construct()
	if bSuccess == false {
		println("construct fail")
		time.Sleep(time.Second)
		return
	}

	wg := constructTCPServer(serverCtx, shutdown, chClosed)

	s := <-sigChan

	switch s {
	case syscall.SIGINT:
		log.Println("Server shutdown...")
		shutdown()
		wg.Wait()
		<-chClosed
		log.Println("Server shutdown completed")
	default:
		panic("unexpected signal has been received")
	}
}
