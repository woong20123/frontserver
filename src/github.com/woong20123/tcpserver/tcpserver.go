package tcpserver

import (
	"context"
	"encoding/binary"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/woong20123/logicmanager"
	"github.com/woong20123/packet"
)

const (
	listenerCloseMatcher = "use of closed network connection"
	maxBufferSize        = 4096
)

var sd StaticData

// Consturct is
func Consturct(serialKey uint32, lm *logicmanager.LogicManager) {
	sd = StaticData{}
	sd.SetSerialkey(serialKey)
	sd.lm = lm
}

// Assembly assembles a read buffer to make a packet.
// kor : Assembly는 패킷을 만들기 위해서 read buffer를 조립합니다.
func Assembly(conn *net.TCPConn, buffer []byte, bufferpos uint32) (resultpos uint32, onPacket bool) {
	resultpos = bufferpos
	var headerFind bool = false
	onPacket = false

	if resultpos < packet.PacketHeaderSize {
		return
	}

	findIndex := resultpos - packet.PacketHeaderSize

	// 전달받은 버퍼를 순회하면서 패킷이 정상적으로 전달되었는지 확인합니다.
	for index := range buffer {
		// must receive At least as much as the header
		// 최소한 헤더만큼 수신받아야 처리 가능,
		// 이곳에 들어왔다면 검사한 버퍼를 모두 버립니다.
		if uint32(index) > findIndex {
			copy(buffer, buffer[index:resultpos])
			resultpos = resultpos - uint32(index)
			break
		}

		if true == packet.HeaderChack(buffer[index:], sd.GetSerialkey()) {
			headerFind = true

			// 만약 패킷 헤더가 처음이 아니라면 나머지 버퍼를 버립니다.
			if index != 0 {
				copy(buffer, buffer[index:resultpos])
				resultpos = resultpos - uint32(index)
			}
			break
		}
	}

	// 패킷 시작지점을 찾았다면
	if true == headerFind {
		PacketSize := uint32(binary.LittleEndian.Uint16(buffer[4:]))
		PacketCommand := uint32(binary.LittleEndian.Uint16(buffer[6:]))
		PacketHeaderSize := packet.PacketHeaderSize
		TotalPacketSize := PacketSize + PacketHeaderSize

		// 패킷을 만들 수 있을 만큼 패킷을 전달 받았다면 패킷을 만들고
		// Logic 처리 goroutine에 전달합니다.
		if TotalPacketSize <= resultpos {
			packet := packet.NewPacket(PacketSize)
			packet.SetHeader(0, uint16(PacketSize), PacketCommand)
			packet.CopyByte(buffer[PacketHeaderSize:TotalPacketSize])
			log.Println("Make packet logic")

			copy(buffer, buffer[TotalPacketSize:resultpos])
			resultpos = resultpos - TotalPacketSize

			sd.lm.CallLogicFun(PacketCommand, conn, packet)
			onPacket = true
		}
	}
	return
}

// HandleRead handles packet read operations for connected sessions
// kor : HandleRead 연결된 세션에 대한 패킷 Read 작업을 처리합니다.
func HandleRead(conn *net.TCPConn, errRead context.CancelFunc) {
	defer errRead()

	// sessesion을 통해서 전달받기 위한 버퍼 생성
	recvBuf := make([]byte, maxBufferSize)

	// session으로부터 전달받은 버퍼를 packet형태로 변환처리하기 위한 Packet
	// TCP의 데이터 전달이 패킷단위로 전달되지 않기 때문에 조립 작업을 합니다.
	AssemblyBuf := make([]byte, maxBufferSize+128)
	var AssemPos uint32 = 0
	onPacket := false

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

		n, err = conn.Write(recvBuf[:n])
		if err != nil {
			log.Println("Write", err)
			return
		}

		if 0 < n {
			copylength := copy(AssemblyBuf[AssemPos:], recvBuf[:n])
			AssemPos += uint32(copylength)

			// 남은 버퍼에서 패킷을 조립할 수 있을 수도 있기 때문에 재호출
			AssemPos, onPacket = Assembly(conn, AssemblyBuf, AssemPos)
			for onPacket == true {
				AssemPos, onPacket = Assembly(conn, AssemblyBuf, AssemPos)
			}

		}
	}
}

// HandleConnection register job for connected session
// kor : HandleConnection은 연결된 세션에 대한 작업을 등록합니다.
func HandleConnection(serverCtx context.Context, conn *net.TCPConn, wg *sync.WaitGroup) {
	defer func() {
		conn.Close()
		wg.Done()
	}()

	readCtx, errRead := context.WithCancel(context.Background())

	go HandleRead(conn, errRead)

	select {
	case <-readCtx.Done():
	case <-serverCtx.Done():
	}
}

func listenerCloseError(err error) bool {
	return strings.Contains(err.Error(), listenerCloseMatcher)
}

// HandleListener register the task to listen to the socket
// kor : HandleListener은 전달된 server address로 소켓을 Listen하는 작업 등록합니다.
func HandleListener(ctxServer context.Context, address string, wg *sync.WaitGroup, chClosed chan struct{}) {

	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		log.Println("ResolveTCPAddr", err)
		return
	}
	tcpListen, err := net.ListenTCP("tcp", tcpAddr)

	if nil != err {
		log.Println("ListenTCP", err)
		return
	}

	// if HandleListener close, it process
	defer func() {
		tcpListen.Close()
		close(chClosed)
	}()

	for {
		conn, err := tcpListen.AcceptTCP()

		// if occur error
		if err != nil {
			if ne, ok := err.(net.Error); ok {
				if ne.Temporary() {
					log.Println("AcceptTCP", err)
					continue
				}
			}
			if listenerCloseError(err) {
				select {
				case <-ctxServer.Done():
					return
				default:
					// fallthrough
				}
			}

			log.Println("AcceptTcp", err)
			return
		}
		wg.Add(1)
		go HandleConnection(ctxServer, conn, wg)
	}
}
