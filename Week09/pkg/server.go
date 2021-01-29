package pkg

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"time"

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
		select {
		case <-ts.ctx.Done():
			log.Printf("listen done\n")
			return nil
		default:
		}
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
	g, eCtx := errgroup.WithContext(ts.ctx)
	g.Go(func() error {
		errChan := make(chan error)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					errChan <- fmt.Errorf("panic %v", r)
				}
			}()
			for {
				data, _, err := rd.ReadLine()
				if err != nil {
					errChan <- err
					return
				}
				readData <- data
			}
		}()
		select {
		case <-eCtx.Done():
			log.Printf("Read done\n")
			return nil
		case err := <-errChan:
			log.Printf("Read err %+v\n", err)
			return err
		}
	})

	g.Go(func() error {
		errChan := make(chan error)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					errChan <- fmt.Errorf("panic %v", r)
				}
			}()
			for {
				data := <-readData
				ts.handler.OnMessage(conn, data)
			}
		}()
		select {
		case <-eCtx.Done():
			log.Printf("onMessage done\n")
			return nil
		case err := <-errChan:
			log.Printf("onMessage err %+v\n", err)
			return err
		}
	})
	err := g.Wait()
	if err != nil {
		log.Printf("consumer err %+v\n", err)
	}
	return err
}

func (ts *TcpServer) Stop() error {
	var err error
	if ts.listener != nil {
		ts.cancel()
		time.Sleep(3 * time.Second)
		err = ts.listener.Close()
	}
	return err
}
