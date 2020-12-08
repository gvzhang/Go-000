package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	g, eCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		var errChan chan error
		go func() {
			errChan <- http.ListenAndServe(":8011", nil)
		}()
		select {
		case <-eCtx.Done():
			return nil
		case err := <-errChan:
			return err
		}
	})
	g.Go(func() error {
		select {
		case <-eCtx.Done():
			return nil
		case <-quit:
			return errors.New("catch exit signal")
		}
	})

	if err := g.Wait(); err != nil {
		fmt.Println("errgroup ", err)
	}
}
