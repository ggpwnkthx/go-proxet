package main

import (
	"io"
	"net"
	"os"
	"strings"
	"sync"
)

func main() {
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
	for {
		l, err := net.Listen(t1[0], t1[1])
		if err != nil {
			return
		}
		for {
			c1, err := l.Accept()
			if err != nil {
				continue
			}
			go connect(c1, dial)
		}
	}
}
func connect(c1 net.Conn, dial string) {
	t2 := strings.Split(dial, ",")
	c2, err := net.Dial(t2[0], t2[1])
	if err != nil {
		return
	}
	go io.Copy(c1, c2)
	io.Copy(c2, c1)
}
