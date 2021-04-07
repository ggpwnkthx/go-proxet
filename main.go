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

type relays struct {
	sync.WaitGroup
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
	defer CleanUp()
	Relays.Add(1)
	for i := 1; i < len(os.Args); i += 2 {
		handle := os.Args[i] + ";" + os.Args[i+1]
		Relays.list[handle] = &relay{}
		fmt.Println("init: " + handle)
		go proxet(handle)
	}
	// Cancellation Context
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		CleanUp()
		os.Exit(0)
	}()
	Relays.Wait()
}
func proxet(handle string) {
	targets := strings.Split(handle, ";")
	t1 := strings.Split(targets[0], ",")
	if t1[0] == "unix" {
		os.Remove(t1[1])
	}
	l1, err := net.Listen(t1[0], t1[1])
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	Relays.list[handle].l1 = &l1
	for {
		if Relays.list[handle] == nil {
			break
		}
		c1, err := (*Relays.list[handle].l1).Accept()
		Relays.list[handle].c1 = &c1
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		go handler(handle)
	}
}
func handler(handle string) {
	targets := strings.Split(handle, ";")
	t2 := strings.Split(targets[1], ",")
	if t2[0] == "unix" {
		os.Remove(t2[1])
	}
	c2, err := net.Dial(t2[0], t2[1])
	if err != nil {
		defer Relays.Done()
		fmt.Println(err.Error())
		return
	}
	Relays.list[handle].c2 = &c2
	go copy(Relays.list[handle].c1, Relays.list[handle].c2) // c1 -> c2
	copy(Relays.list[handle].c2, Relays.list[handle].c1)    // c2 -> c1
}
func copy(writer *net.Conn, reader *net.Conn) {
	_, err := io.Copy(*writer, *reader)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func CleanUp() {
	for handle, r := range Relays.list {
		if r.l1 != nil {
			(*r.l1).Close()
		}
		if r.c1 != nil {
			//(*r.c1).Close()
		}
		if r.c2 != nil {
			(*r.c2).Close()
		}
		delete(Relays.list, handle)
	}
}
