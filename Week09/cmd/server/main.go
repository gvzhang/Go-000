package main

import (
	"Go-000/Week09/pkg"
	"bufio"
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
	writer := bufio.NewWriter(c)
	_, err := writer.Write([]byte("receive "))
	if err != nil {
		log.Printf("writer write error %+v\n", err)
		return
	}
	_, err = writer.Write(bytes)
	if err != nil {
		log.Printf("writer write error %+v\n", err)
		return
	}
	err = writer.Flush()
	if err != nil {
		log.Printf("writer flush error %+v\n", err)
	}
}

func (bh *BizHandler) OnClose(c net.Conn, err error) {
	log.Printf("OnClose %+v %+v\n", c, err)
}

func main() {
	serAddr := "127.0.0.1:50000"
	tcpAddr, err := net.ResolveTCPAddr("tcp", serAddr)
	if err != nil {
		panic("ResolveTCPAddr failed:" + err.Error())
	}
	bizHandler := new(BizHandler)
	server := pkg.NewTcpServer(tcpAddr, bizHandler)

	errChan := make(chan error, 0)
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
	err = server.Stop()
	if err != nil {
		log.Println("server.Stop error", err)
	}
	log.Println("server exit...")
}
