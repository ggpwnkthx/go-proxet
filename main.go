package main

import (
	"io"
	"net"
	"os"
	"strings"
	"sync"
)

type relays struct {
	sync.Mutex
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
	wg := sync.WaitGroup{}
	for i := 1; i < len(os.Args); i += 2 {
		wg.Add(1)
		go Proxet(os.Args[i], os.Args[i+1], &wg)
	}
	wg.Wait()
}

func Proxet(listen string, dial string, wg *sync.WaitGroup) {
	defer wg.Done()
	t1 := strings.Split(listen, ",")
	Relays.list[listen+";"+dial] = &relay{}
	for {
		l, err := net.Listen(t1[0], t1[1])
		if err != nil {
			return
		}
		for {
			c1, err := l.Accept()
			Relays.list[listen+";"+dial].c1 = &c1
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
	Relays.list[listen+";"+dial].c2 = &c2
	go process(listen + ";" + dial)
}
func process(handle string) {
	go io.Copy((*Relays.list[handle].c1), (*Relays.list[handle].c2))
	io.Copy((*Relays.list[handle].c2), (*Relays.list[handle].c1))
}
