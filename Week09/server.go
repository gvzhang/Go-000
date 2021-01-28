package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"

	"golang.org/x/sync/errgroup"
)

type Handler interface {
	OnConnect(c net.Conn)
	OnMessage(c net.Conn, bytes []byte)
	OnClose(c net.Conn, err error)
}

func NewTcpServer(addr *net.TCPAddr, handler Handler) *TcpServer {
	ts := new(TcpServer)
	ts.addr = addr
	ts.handler = handler
	ts.ctx, ts.cancel = context.WithCancel(context.Background())
	return ts
}

type TcpServer struct {
	ctx      context.Context
	cancel   context.CancelFunc
	handler  Handler
	addr     *net.TCPAddr
	listener *net.TCPListener
}

func (ts *TcpServer) Run() error {
	var err error
	ts.listener, err = net.ListenTCP("tcp", ts.addr)
	if err != nil {
		return err
	}

	for {
		conn, err := ts.listener.Accept()
		if err != nil {
			log.Printf("accept fail, err: %v\n", err)
			continue
		}
		go func(conn net.Conn) {
			var consumerErr error
			defer func() {
				if r := recover(); r != nil {
					consumerErr = fmt.Errorf("panic %v", r)
				}
				ts.handler.OnClose(conn, consumerErr)
				err := conn.Close()
				if err != nil {
					log.Printf("conn.Close, err: %v\n", err)
				}
			}()
			consumerErr = ts.consumer(conn)
		}(conn)
	}
}

func (ts *TcpServer) consumer(conn net.Conn) error {
	ts.handler.OnConnect(conn)

	readData := make(chan []byte, 10)
	rd := bufio.NewReader(conn)
	ctx := context.Background()
	g, eCtx := errgroup.WithContext(ts.ctx)
	g.Go(func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("panic %v", r)
			}
		}()

		for {
			select {
			case <-eCtx.Done():
				return nil
			default:
			}
			data, _, rErr := rd.ReadLine()
			if rErr != nil {
				err = rErr
				return
			}
			readData <- data
		}
	})
	g.Go(func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("panic %v", r)
			}
		}()
		for {
			select {
			case <-eCtx.Done():
				return nil
			default:
			}
			data := <-readData
			ts.handler.OnMessage(conn, data)
		}
	})
	return g.Wait()
}

func (ts *TcpServer) Stop() error {
	var err error
	if ts.listener != nil {
		ts.cancel()
		err = ts.listener.Close()
	}
	return err
}
