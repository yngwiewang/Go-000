package main

// 基于 errgroup 实现一个 http server 的启动和关闭 ，
// 以及 linux signal 信号的注册和处理，要保证能够 一个退出，全部注销退出。

// 发送 INT 或者 TERM 信号后，新的请求立刻会返回连接拒绝（Connection refused），
// 但是已经发出的请求会继续处理，10秒内处理完成的话会正常返回
// 10秒内未完成处理会返回连接失败（Connection aborted）。

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "", log.Lshortfile|log.Lmicroseconds)
	rand.Seed(time.Now().UnixNano())
}

func getExceptionSvr(ctx context.Context, wg *sync.WaitGroup, cancel context.CancelFunc) func() error {
	return func() error {
		var err error
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("throw exception"))
			err = errors.New("exception server triggered")
			cancel()
		})
		server := &http.Server{Addr: ":8001", Handler: mux}

		go func() {
			<-ctx.Done()
			logger.Println("the exception server is closing...")
			server.Shutdown(context.Background())
			logger.Println("the exception server is closed")
			wg.Done()
		}()
		logger.Println("the exception server is starting...")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			return fmt.Errorf("error starting the exception server: %s", err)
		}
		wg.Wait()
		return err
	}
}

func getEchoSvr(ctx context.Context, wg *sync.WaitGroup) func() error {
	return func() error {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			io.Copy(w, r.Body)
		})
		server := &http.Server{Addr: ":8002", Handler: mux}

		errChan := make(chan error, 1)
		go func() {
			<-ctx.Done()
			shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			logger.Println("the echo server is closing...")
			if err := server.Shutdown(shutCtx); err != nil {
				errChan <- fmt.Errorf("error shutting down the echo server: %s", err)
			}
			logger.Println("the echo server is closed")
			close(errChan)
			wg.Done()
		}()

		logger.Println("the echo server is starting...")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			return fmt.Errorf("error starting the echo server: %s", err)
		}
		err := <-errChan
		wg.Wait()
		return err
	}
}

func getSleepSvr(ctx context.Context, wg *sync.WaitGroup) func() error {
	return func() error {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			t := rand.Intn(20)
			time.Sleep(time.Duration(t) * time.Second)
			w.WriteHeader(http.StatusOK)
			s := fmt.Sprintf("I have slept %d seconds\n", t)
			w.Write([]byte(s))
		})
		server := &http.Server{Addr: ":8000", Handler: mux}
		// 如果是超时强制停止服务器，也就是有未完成的在途请求，就向上返回这个错误信息
		errChan := make(chan error, 1)
		go func() {
			<-ctx.Done()
			shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if err := server.Shutdown(shutCtx); err != nil {
				// 超时关闭
				errChan <- fmt.Errorf("error shutting down the sleep server: %s", err)
			}
			logger.Println("the sleep server is closed")
			// 正常关闭
			close(errChan)
			wg.Done()
		}()

		logger.Println("the sleep server is starting...")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			return fmt.Errorf("error starting the sleep server: %s", err)
		}
		logger.Println("the sleep server is closing...")
		err := <-errChan
		wg.Wait()
		return err
	}
}

func main() {

	var wg sync.WaitGroup
	wg.Add(3)

	// WithContext 返回一个Group和一个由父ctx派生的新context
	// 派生的context有两种情况会被触发cancel：
	// 1. 传给Go的函数第一次返回non-nil的error
	// 2. 第一次Wait方法返回
	// 以上两种情况哪个先发生，哪个触发派生context的cancel
	eg, egCtx := errgroup.WithContext(context.Background())

	ctx, cancel := context.WithCancel(context.Background())

	// Go在一个新的goroutine中调用作为参数的函数。
	// 第一个返回非空error的调用cancel掉整个errgroup；它的error
	// 将会被Wait方法返回
	eg.Go(getSleepSvr(ctx, &wg))
	eg.Go(getEchoSvr(ctx, &wg))
	eg.Go(getExceptionSvr(ctx, &wg, cancel))

	go func() {
		<-egCtx.Done()
		cancel()
	}()

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		<-signals
		cancel()
	}()

	// Wait方法阻塞直到所有Go方法的函数都返回为止
	// 然后Wait方法返回这些函数中的第一个non-nil error(如果有的话)
	if err := eg.Wait(); err != nil {
		logger.Printf("error in the server goroutines: %s\n", err)
		os.Exit(1)
	}

	logger.Println("servers closed successfully")
}
