package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {
	for _, a := range os.Args {
		fmt.Println(a)
	}
	for i := 0; i < len(os.Args); i += 2 {
		t1 := strings.Split(os.Args[i], ",")
		t2 := strings.Split(os.Args[1+1], ",")
		go handle(t1, t2)
	}
}
func handle(listen []string, dial []string) {
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
			go connect(c1, dial)
		}
	}
}
func connect(c1 net.Conn, target []string) {
	defer c1.Close()
	c2, err := net.Dial(target[0], target[1])
	if err != nil {
		fmt.Println(err.Error())
	}
	go io.Copy(c1, c2)
	io.Copy(c2, c1)
}
