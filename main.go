package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 0 {
		for i := 0; i < len(os.Args); i += 2 {
			t1 := strings.Split(os.Args[i], ",")
			t2 := strings.Split(os.Args[1+1], ",")
			go handle(t1, t2)
		}
	}
}
func handle(t1 []string, t2 []string) {
	for {
		fmt.Println("starting listener of type " + t1[0] + " at " + t1[1])
		l, _ := net.Listen(t1[0], t1[1])
		for {
			c1, err := l.Accept()
			if err != nil {
				fmt.Fprintf(os.Stderr, err.Error())
				continue
			}
			fmt.Println("starting dialer of type " + t2[0] + " at " + t2[1])
			go connect(c1, t2)
		}
	}
}

func connect(c1 net.Conn, target []string) {
	defer c1.Close()
	c2, err := net.Dial(target[0], target[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}
	go io.Copy(c1, c2)
	io.Copy(c2, c1)
}
