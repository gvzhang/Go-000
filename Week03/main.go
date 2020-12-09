package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()
	g, eCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return startServer(ctx, eCtx, ":8015")
	})
	g.Go(func() error {
		return startServer(ctx, eCtx, ":81")
	})
	g.Go(func() error {
		return stopSignal(ctx, eCtx)
	})
	if err := g.Wait(); err != nil {
		log.Println("errgroup", err)
	}
	log.Println("main exit...")
}

func startServer(ctx context.Context, egCtx context.Context, addr string) error {
	errChan := make(chan error, 0)
	server := &http.Server{Addr: addr, Handler: nil}
	go func() {
		errChan <- server.ListenAndServe()
	}()
	select {
	case <-egCtx.Done():
		log.Println(addr, "http exit...")
		tCtx, _ := context.WithTimeout(ctx, 10*time.Second)
		err := server.Shutdown(tCtx)
		if err != nil {
			log.Println(addr, "shutdown err", err)
		}
		return nil
	case err := <-errChan:
		return err
	}
}

func stopSignal(ctx context.Context, egCtx context.Context) error {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case <-egCtx.Done():
		log.Println("signal exit...")
		return nil
	case q := <-quit:
		return fmt.Errorf("catch exit signal %s", q.String())
	}
}
