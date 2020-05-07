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
	bSuccess = bSuccess && constructTCPClient()
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
	sessionm := tcpserver.Instance().SessionMgr()
	examserverlogic.RegistSessionLogic(sessionm)
	return true
}

func constructLogic() bool {
	// regist exam Logic
	lm := tcpserver.Instance().LogicManager()

	switch examserverlogic.Instance().ConfigMgr().ServerConfig().ServerMode {
	case "front":
		examserverlogic.FrontModeRegistCommandLogic(lm)
	default:
		examserverlogic.ChatServerModeRegistCommandLogic(lm)
	}

	// set logic goroutines count
	lm.RunLogicHandle(runtime.NumCPU())
	return true
}

func constructTCPClient() bool {
	srvConfig := examserverlogic.Instance().ConfigMgr().ServerConfig()
	if "front" == srvConfig.ServerMode {
		ChatserverIP := srvConfig.BackEndChatServerIP
		ChatserverPort := srvConfig.BackEndChatServerPort

		err := tcpserver.Instance().TCPClientMgr().AddTCPClient(examshare.TCPCliToSvrIdxChat, ChatserverIP, ChatserverPort)
		if err != nil {
			println("Not Connect to BackEndChatServer ", ChatserverIP, ":", ChatserverPort)
			examserverlogic.Logger().Println("Not Connect to BackEndChatServer ", ChatserverIP, ":", ChatserverPort)
			examserverlogic.Logger().Println(err.Error())
			return false
		} else {
			tcpserver.Instance().SendManager().RunSendToServerHandle(examshare.TCPCliToSvrIdxChat)

			p := packet.Pool().AcquirePacket()
			p.SetHeaderByDefaultKey(0, int32(examshare.Cmd_F2CSServerRegistReq))
			req := examshare.F2CS_ServerRegistReq{}
			req.Ip = srvConfig.ServerIP
			req.Port = int32(srvConfig.ServerPort)
			req.ServerMode = examshare.SrvMode_FrontSrvMode
			err := p.MarshalFromProto(&req)
			if err == nil {
				tcpserver.Instance().SendManager().SendToServerConn(examshare.TCPCliToSvrIdxChat, p)
			} else {
				packet.Pool().ReleasePacket(p)
			}
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
