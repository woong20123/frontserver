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

func constructTCPSession() {
	sessionm := tcpserver.Instance().GetSessionMgr()
	examserverlogic.RegistSessionLogic(sessionm)
}

func constructLogic() {
	// regist exam Logic
	lm := tcpserver.Instance().LogicManager()
	examserverlogic.RegistCommandLogic(lm)
	// set logic goroutines count
	lm.RunLogicHandle(runtime.NumCPU())
}

func constructTCPServer(port int, chClosed chan struct{}) (wg sync.WaitGroup, cancel context.CancelFunc) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	serverCtx, shutdown := context.WithCancel(context.Background())
	// server consturct
	tcpserver.Consturct(share.ExamplePacketSerialkey, runtime.NumCPU())

	// start server handler
	address := ":" + strconv.Itoa(port)
	go tcpserver.HandleListener(serverCtx, address, &wg, chClosed)
	println("[Server 정보] ", address)
	examserverlogic.GetLogger().Println("[Server 정보] ", address)
	cancel = shutdown
	return
}

func main() {
	chClosed := make(chan struct{})
	examserverlogic.GetLogger().Println(fmt.Sprint("[Server Ver ", share.ExamVer, "]"))
	println(fmt.Sprint("[Server Ver ", share.ExamVer, "]"))

	constructTCPSession()
	constructLogic()
	wg, shutdown := constructTCPServer(20224, chClosed)

	sigChan := make(chan os.Signal, 1)
	signal.Ignore()
	signal.Notify(sigChan, syscall.SIGINT)

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
