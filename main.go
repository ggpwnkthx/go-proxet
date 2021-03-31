package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

type proxet struct {
	listener *net.Listener
	dialer   *net.Conn
}

var Proxettes = struct {
	sync.RWMutex
	sync.WaitGroup
	list map[string]*proxet
}{
	list: map[string]*proxet{},
}

func main() {
	// Cancellation Context
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		CleanUp()
		os.Exit(0)
	}()

	// Parse args
	for i := 1; i < len(os.Args); i += 2 {
		Proxettes.Add(1)
		go listen(os.Args[i], os.Args[i+1])
	}
	Proxettes.Wait()
	for len(Proxettes.list) > 0 {
		for k, p := range Proxettes.list {
			if p.dialer == nil {
				c1, err := (*p.listener).Accept()
				if err != nil {
					fmt.Println(err.Error())
					continue
				}
				dial := strings.Split(k, ";")[1]
				go connect(c1, dial)
			}
		}
	}
	CleanUp()
}
func listen(listen string, dial string) {
	defer Proxettes.Done()
	fmt.Println("opening " + listen)
	t1 := strings.Split(listen, ",")
	p := new(proxet)
	l, err := net.Listen(t1[0], t1[1])
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	p.listener = &l
	Proxettes.Lock()
	Proxettes.list[listen+";"+dial] = p
	Proxettes.Unlock()
}
func connect(c1 net.Conn, target string) {
	fmt.Println("opening " + target)
	t2 := strings.Split(target, ",")
	c2, err := net.Dial(t2[0], t2[1])
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	Proxettes.Lock()
	Proxettes.list[c1.LocalAddr().String()+";"+target].dialer = &c2
	Proxettes.Unlock()
}

func CleanUp() {
	for k, p := range Proxettes.list {
		fmt.Println("closing " + strings.Split(k, ";")[0])
		(*p.listener).Close()
		if p.dialer == nil {
			fmt.Println("closing " + strings.Split(k, ";")[1])
			(*p.dialer).Close()
		}
	}
}
