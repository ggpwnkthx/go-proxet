package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

type relays struct {
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
		Relays.list[handle].l1 = &l1
	}
	for {
		if Relays.list[handle].c1 == nil {
			c1, err := (*Relays.list[handle].l1).Accept()
			Relays.list[handle].c1 = &c1
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
		Relays.list[handle].c2 = &c2
	}
	go process(handle)
}
func process(handle string) {
	go io.Copy((*Relays.list[handle].c1), (*Relays.list[handle].c2))
	_, err := io.Copy((*Relays.list[handle].c2), (*Relays.list[handle].c1))
	if err != nil {
		fmt.Println(handle + ": " + err.Error())
	}
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
