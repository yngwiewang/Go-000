// Week09 实现一个tcp server ，用两个 goroutine 读写 conn，
// 两个 goroutine 通过 chan 可以传递 message，能够正确退出。
package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Server struct {
	listener net.Listener
	quit     chan struct{}
	file     *os.File
	message  chan []byte
	wg       sync.WaitGroup
}

func main() {
	var (
		address string
		file    string
	)

	log.SetFlags(log.Llongfile)

	flag.StringVar(&address, "l", "127.0.0.1:8000", "listening address of the TCP server")
	flag.StringVar(&file, "f", "temp.txt", "file to write")
	flag.Parse()

	// 接收信号，触发 server 正确退出。
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	s := NewServer(address, file)
	go func() {
		<-term
		s.Stop()
	}()

	s.Serv()
}

// NewServer 初始化 server。
func NewServer(address, file string) *Server {
	s := &Server{
		quit:    make(chan struct{}),
		message: make(chan []byte),
	}
	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("open file err: %s", err.Error())
	}
	l, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("listen address err: %s", err.Error())
	}
	s.file = f
	s.listener = l
	return s
}

// 停止 server，等待在途的 goroutine 执行完毕，关闭 channel。
func (s *Server) Stop() {
	close(s.quit)
	s.listener.Close()
	s.wg.Wait()
	close(s.message)
}

// 启动服务，如果是正常关闭触发的监听关闭就返回。
func (s *Server) Serv() {
	s.wg.Add(1)
	go s.writeContent()

	log.Println("start serving...")
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.quit:
				log.Println("closing server gracefully...")
				return
			default:
				log.Printf("accept err: %s", err.Error())
				continue
			}
		}
		s.wg.Add(1)
		go s.readConn(conn)
	}
}

func (s *Server) readConn(conn net.Conn) {
	defer s.wg.Done()
	defer conn.Close()

	buf := bufio.NewReader(conn)
	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			if err != io.EOF {
				log.Printf("read content err: %v\n", err)
			}
			return
		}
		s.message <- line
	}
}

func (s *Server) writeContent() {
	defer s.wg.Done()
	defer s.file.Close()
	for content := range s.message {
		content = append(content, []byte("\n")...)
		n, err := s.file.Write(content)
		if err != nil {
			log.Printf("write file err: %v\n", err)
		}
		log.Printf("write %d bytes\n", n)
	}
}
