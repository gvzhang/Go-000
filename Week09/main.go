package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type BizHandler struct {
}

func (bh *BizHandler) OnConnect(c net.Conn) {
	log.Printf("OnConnect %+v\n", c)
}

func (bh *BizHandler) OnMessage(c net.Conn, bytes []byte) {
	log.Printf("OnMessage %+v %s\n", c, string(bytes))
}

func (bh *BizHandler) OnClose(c net.Conn, err error) {
	log.Printf("OnClose %+v %+v\n", c, err)
}

func main() {
	errChan := make(chan error, 0)
	tcpAddr := new(net.TCPAddr)
	tcpAddr.IP = net.IPv4(127, 0, 0, 1)
	tcpAddr.Port = 50000

	bizHandler := new(BizHandler)
	server := NewTcpServer(tcpAddr, bizHandler)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				errChan <- fmt.Errorf("panic %v", r)
			}
		}()
		errChan <- server.Run()
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case err := <-errChan:
		if err != nil {
			log.Println("server.Run error", err)
		}
	case q := <-quit:
		log.Printf("catch exit signal %s\n", q.String())
	}
	err := server.Stop()
	if err != nil {
		log.Println("server.Stop error", err)
	}
	log.Println("server exit...")
}
