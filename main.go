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
	c1 *net.Conn
	c2 *net.Conn
}

var Relays relays

func main() {
	Relays = relays{
		list: map[string]*relay{},
	}
	for i := 1; i < len(os.Args); i += 2 {
		Relays.list[os.Args[i]+";"+os.Args[i+1]] = &relay{}
	}
	for len(Relays.list) > 0 {
		for handle, r := range Relays.list {
			targets := strings.Split(handle, ";")
			t1 := strings.Split(targets[0], ",")
			if r.c1 == nil {
				fmt.Println("opening: " + t1[1])
				l, err := net.Listen(t1[0], t1[1])
				if err != nil {
					delete(Relays.list, handle)
					continue
				}
				for {
					c1, err := l.Accept()
					Relays.Lock()
					Relays.list[handle].c1 = &c1
					Relays.Unlock()
					if err != nil {
						continue
					}
					go connect(targets[0], targets[1])
				}
			}
		}
	}
	CleanUp()
}

func connect(listen string, dial string) {
	t2 := strings.Split(dial, ",")
	c2, err := net.Dial(t2[0], t2[1])
	if err != nil {
		return
	}
	Relays.Lock()
	defer Relays.Unlock()
	Relays.list[listen+";"+dial].c2 = &c2
	go process(listen + ";" + dial)
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
