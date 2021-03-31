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
	list map[string]*relay
}

type relay struct {
	l1 *net.Listener
	c1 *net.Conn
	c2 *net.Conn
}

var Relays relays

func main() {
	Relays = relays{
		list: map[string]*relay{},
	}
	for i := 1; i < len(os.Args); i += 2 {
		handle := os.Args[i] + ";" + os.Args[i+1]
		Relays.list[handle] = &relay{}
		go proxet(handle)
	}
	for len(Relays.list) > 0 {
		Relays.RLock()
		list := Relays.list
		Relays.RUnlock()
		for handle, r := range list {
			targets := strings.Split(handle, ";")
			t1 := strings.Split(targets[0], ",")
			if r.l1 == nil {
				fmt.Println("opening: " + t1[1])
				l, err := net.Listen(t1[0], t1[1])
				if err != nil {
					continue
				} else {
					r.l1 = &l
				}
			} else {
				go proxet(handle)
			}
		}
	}
	CleanUp()
}
func proxet(handle string) {
	if Relays.list[handle].l1 == nil {
		targets := strings.Split(handle, ";")
		t1 := strings.Split(targets[0], ",")
		l1, err := net.Listen(t1[0], t1[1])
		if err != nil {
			return
		}
		Relays.Lock()
		Relays.list[handle].l1 = &l1
		Relays.Unlock()
	}
	for {
		if Relays.list[handle].c1 == nil {
			Relays.Lock()
			c1, err := (*Relays.list[handle].l1).Accept()
			Relays.list[handle].c1 = &c1
			Relays.Unlock()
			if err != nil {
				continue
			}
		}
		go connect(handle)
	}
}
func connect(handle string) {
	if Relays.list[handle].c2 == nil {
		targets := strings.Split(handle, ";")
		t2 := strings.Split(targets[1], ",")
		c2, err := net.Dial(t2[0], t2[1])
		if err != nil {
			return
		}
		Relays.Lock()
		Relays.list[handle].c2 = &c2
		Relays.Unlock()
	}
	go process(handle)
}
func process(handle string) {
	Relays.RLock()
	defer Relays.RUnlock()
	go io.Copy((*Relays.list[handle].c1), (*Relays.list[handle].c2))
	io.Copy((*Relays.list[handle].c2), (*Relays.list[handle].c1))
}

func CleanUp() {
	for _, r := range Relays.list {
		if r.c1 != nil {
			(*r.c1).Close()
		}
		if r.c2 != nil {
			(*r.c2).Close()
		}
	}
}
