package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
)

type relays struct {
	sync.RWMutex
	sync.WaitGroup
	list map[string]*relay
}

type relay struct {
	c1 *net.Conn
	c2 *net.Conn
}

var Relays relays

func main() {
	Relays = relays{
		list: map[string]*relay{},
	}
	for i := 1; i < len(os.Args); i += 2 {
		Relays.Add(1)
		go Proxet(os.Args[i], os.Args[i+1])
	}
	Relays.Wait()
	for len(Relays.list) > 0 {
		for _, r := range Relays.list {
			if r.c2 == nil {
				fmt.Println("waiting for connection to " + (*r.c1).LocalAddr().String())
			} else {
				fmt.Println("relaying " + (*r.c1).LocalAddr().String() + " to " + (*r.c2).LocalAddr().String())
			}
		}
	}
}

func Proxet(listen string, dial string) {
	defer Relays.Done()
	t1 := strings.Split(listen, ",")
	Relays.Lock()
	Relays.list[listen+";"+dial] = &relay{}
	Relays.Unlock()
	for {
		l, err := net.Listen(t1[0], t1[1])
		if err != nil {
			return
		}
		for {
			c1, err := l.Accept()
			Relays.Lock()
			Relays.list[listen+";"+dial].c1 = &c1
			Relays.Unlock()
			if err != nil {
				continue
			}
			go connect(listen, dial)
		}
	}
}
func connect(listen string, dial string) {
	t2 := strings.Split(dial, ",")
	c2, err := net.Dial(t2[0], t2[1])
	if err != nil {
		return
	}
	Relays.Lock()
	Relays.list[listen+";"+dial].c2 = &c2
	Relays.Unlock()
	go process(listen + ";" + dial)
}
func process(handle string) {
	Relays.RLock()
	defer Relays.RUnlock()
	go io.Copy((*Relays.list[handle].c1), (*Relays.list[handle].c2))
	io.Copy((*Relays.list[handle].c2), (*Relays.list[handle].c1))
}
