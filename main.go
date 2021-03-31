package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	for i := 1; i < len(os.Args); i += 2 {
		wg.Add(1)
		t1 := strings.Split(os.Args[i], ",")
		if t1[0] == "unix" {
			defer os.Remove(t1[1])
		}
		t2 := strings.Split(os.Args[1+1], ",")
		if t2[0] == "unix" {
			defer os.Remove(t2[1])
		}
		go handle(t1, t2, &wg)
	}
	wg.Wait()
}
func handle(listen []string, dial []string, wg *sync.WaitGroup) {
	for {
		fmt.Println("starting listener of type " + listen[0] + " at " + listen[1])
		l, err := net.Listen(listen[0], listen[1])
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		for {
			c1, err := l.Accept()
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			fmt.Println("starting dialer of type " + dial[0] + " at " + dial[1])
			go connect(c1, dial, wg)
		}
	}
}
func connect(c1 net.Conn, target []string, wg *sync.WaitGroup) {
	defer c1.Close()
	defer wg.Done()
	fmt.Println("starting dialer of type " + target[0] + " at " + target[1])
	c2, err := net.Dial(target[0], target[1])
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	go io.Copy(c1, c2)
	io.Copy(c2, c1)
}
