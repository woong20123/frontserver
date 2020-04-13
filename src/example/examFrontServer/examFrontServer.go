package main

// "factored" import statment
import (
	"context"
	"example/examlogic"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/woong20123/logicmanager"
	"github.com/woong20123/tcpserver"
)

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Ignore()
	signal.Notify(sigChan, syscall.SIGINT)

	var wg sync.WaitGroup
	chClosed := make(chan struct{})

	serverCtx, shutdown := context.WithCancel(context.Background())

	// regist serial key
	lm := logicmanager.NewLogicManager()
	//examserverlogic.RegistCommandLogic(lm)
	tcpserver.Consturct(examlogic.ExamplePacketSerialkey, lm)

	address := ":20224"
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
