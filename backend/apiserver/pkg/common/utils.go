package common

//
//import (
//	"context"
//	"os"
//	"os/signal"
//	"syscall"
//	"time"
//)
//
//func ListenStopSignalAndGracefulStop(ctx context.Context, serverS ...interface{}) chan error {
//	if len(serverS) == 0 {
//		panic("want to get some server, but get nothings!")
//	}
//
//	shutdownFuncS := make([]func() error, len(serverS), len(serverS))
//	for i, s := range serverS {
//		t := s
//		var shutdownFunc func() error
//		switch s.(type) {
//		case canShutdown:
//			shutdownFunc = func() error {
//				t := t.(canShutdown)
//				log.Info(ctx, "server will shutdown on one minute")
//				timeoutCtx, cancelFunc := context.WithTimeout(context.Background(), time.Minute)
//				defer cancelFunc()
//				return t.Shutdown(timeoutCtx)
//			}
//		case grpcShutdown:
//			shutdownFunc = func() error {
//				t := t.(grpcShutdown)
//				log.Info(ctx, "grpc server will shutdown, wait for minutes")
//				t.GracefulStop()
//				return nil
//			}
//		default:
//			panic("not support server, the server should have function like \"Shutdown(ctx context.Context) error\" or \"GracefulStop()\".")
//		}
//		shutdownFuncS[i] = shutdownFunc
//	}
//
//	done := make(chan os.Signal)
//	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
//
//	ch := make(chan error)
//	go func() {
//		defer close(ch)
//		select {
//		case s := <-done:
//			log.Info(ctx, "get stop signal, shutting down server", "signal", s.String())
//		case <-ctx.Done():
//			log.Info(ctx, "context done, shutting down server")
//		}
//
//		g := errgroup.Group{}
//		for _, f := range shutdownFuncS {
//			g.Go(f)
//		}
//		ch <- g.Wait()
//	}()
//	return ch
//}
