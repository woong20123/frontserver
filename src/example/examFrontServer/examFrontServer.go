package main

// "factored" import statment
import (
	"context"
	"example/examFrontServer/examserverlogic"
	"example/share"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"

	"github.com/woong20123/logicmanager"
	"github.com/woong20123/tcpserver"
)

func constructTCPServer(port int) {

	sigChan := make(chan os.Signal, 1)
	signal.Ignore()
	signal.Notify(sigChan, syscall.SIGINT)

	var wg sync.WaitGroup
	chClosed := make(chan struct{})
	serverCtx, shutdown := context.WithCancel(context.Background())

	// set LogicManager
	lm := logicmanager.GetInstance()
	runtime.GOMAXPROCS(runtime.NumCPU())
	lm.RunLogicHandle(runtime.NumCPU())

	// regist exam Logic
	examserverlogic.RegistCommandLogic(lm)
	tcpserver.Consturct(share.ExamplePacketSerialkey, lm)

	// start server handler
	address := ":" + strconv.Itoa(port)
	go tcpserver.HandleListener(serverCtx, address, &wg, chClosed)
	log.Println("On Server ", address)

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

func main() {
	constructTCPServer(20224)
}
