package main

// "factored" import statment
import (
	"context"
	"example/examFrontServer/examserverlogic"
	"example/share"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"

	"github.com/woong20123/tcpserver"
)

func construct() {
	bSuccess := true
	bSuccess = bSuccess && constructConfig()
	bSuccess = bSuccess && constructTCPSession()
	bSuccess = bSuccess && constructLogic()
	bSuccess = bSuccess && constructTCPClient()
}

func constructConfig() bool {
	examserverlogic.MakeExampleConfig()
	examserverlogic.Instance().ConfigMgr().ReadConfig("\\ExampleServerConfig.json")
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
	case "main":
		examserverlogic.ChatServerModeRegistCommandLogic(lm)
	default:
		examserverlogic.ChatServerModeRegistCommandLogic(lm)
	}

	// set logic goroutines count
	lm.RunLogicHandle(runtime.NumCPU())
	return true
}

const (
	// TCPCliToSvrIdxChat is
	TCPCliToSvrIdxChat uint32 = iota + 0x10
)

func constructTCPClient() bool {
	srvConfig := examserverlogic.Instance().ConfigMgr().ServerConfig()
	if "front" == srvConfig.ServerMode {
		ChatserverIP := srvConfig.BackEndChatServerIP
		ChatserverPort := srvConfig.BackEndChatServerPort

		err := tcpserver.Instance().TCPClientMgr().AddTCPClient(TCPCliToSvrIdxChat, ChatserverIP, ChatserverPort)
		if err != nil {
			examserverlogic.Logger().Println(err.Error())
			return false
		} else {
			tcpserver.Instance().SendManager().RunSendToServerHandle(TCPCliToSvrIdxChat)
		}
	}
	return true
}

func constructTCPServer(ctxServer context.Context, cancel context.CancelFunc, chClosed chan struct{}) (wg sync.WaitGroup) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// server consturct
	tcpserver.Consturct(share.ExamplePacketSerialkey, runtime.NumCPU())

	// start server handler
	port := examserverlogic.Instance().ConfigMgr().ServerConfig().ServerPort
	address := ":" + strconv.Itoa(port)
	go tcpserver.HandleListener(ctxServer, address, &wg, chClosed)
	println("[Server 정보] ", address)
	examserverlogic.Logger().Println("[Server 정보] ", address)
	return
}

func main() {

	sigChan := make(chan os.Signal, 1)
	signal.Ignore()
	signal.Notify(sigChan, syscall.SIGINT)

	chClosed := make(chan struct{})
	examserverlogic.Logger().Println(fmt.Sprint("[Server Ver ", share.ExamVer, "]"))
	println(fmt.Sprint("[Server Ver ", share.ExamVer, "]"))

	serverCtx, shutdown := context.WithCancel(context.Background())

	construct()
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
