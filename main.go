package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

var unixSocketPaths []string

func main() {
	var wg sync.WaitGroup
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		DeleteFiles(&wg)
		os.Exit(0)
	}()

	for i := 1; i < len(os.Args); i += 2 {
		wg.Add(1)
		t1 := strings.Split(os.Args[i], ",")
		t2 := strings.Split(os.Args[1+1], ",")
		go handle(t1, t2, &wg)
	}
	wg.Wait()
}
func handle(listen []string, dial []string, wg *sync.WaitGroup) {
	defer wg.Done()
	if listen[0] == "unix" {
		defer os.Remove(listen[1])
		unixSocketPaths = append(unixSocketPaths, listen[1])
	}
	for {
		fmt.Println("starting listener of type " + listen[0] + " at " + listen[1])
		l, err := net.Listen(listen[0], listen[1])
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		for {
			c1, err := l.Accept()
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			fmt.Println("starting dialer of type " + dial[0] + " at " + dial[1])
			go connect(c1, dial)
		}
	}
}
func connect(c1 net.Conn, target []string) {
	if target[0] == "unix" {
		defer os.Remove(target[1])
		unixSocketPaths = append(unixSocketPaths, target[1])
	}
	c2, err := net.Dial(target[0], target[1])
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	go io.Copy(c1, c2)
	io.Copy(c2, c1)
}

func DeleteFiles(wg *sync.WaitGroup) {
	for _, socketPath := range unixSocketPaths {
		if _, err := os.Stat(socketPath); err == nil {
			os.Remove(socketPath)
			wg.Done()
		}
	}
}
