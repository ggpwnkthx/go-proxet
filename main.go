package main

import (
	"io"
	"net"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 0 {
		for _, arg := range os.Args {
			if strings.Contains(arg, "->") {
				targets := strings.Split(arg, "->")
				t1 := strings.Split(targets[0], ",")
				t2 := strings.Split(targets[1], ",")
				go handle(t1, t2, 0)
			} else if strings.Contains(arg, "<-") {
				targets := strings.Split(arg, "<-")
				t1 := strings.Split(targets[0], ",")
				t2 := strings.Split(targets[1], ",")
				go handle(t1, t2, 1)
			} else if strings.Contains(arg, "<>") {
				targets := strings.Split(arg, "<>")
				t1 := strings.Split(targets[0], ",")
				t2 := strings.Split(targets[1], ",")
				go handle(t1, t2, 2)
			}
		}
	}
}
func handle(t1 []string, t2 []string, direction int) {
	for {
		l, _ := net.Listen(t1[0], t1[1])
		for {
			c1, err := l.Accept()
			if err != nil {
				continue
			}
			switch direction {
			case 0:
				go outputOnly(c1, t2)
			case 1:
				go inputOnly(c1, t2)
			case 2:
				go bidirectional(c1, t2)
			}
		}
	}
}

func outputOnly(c1 net.Conn, target []string) {
	defer c1.Close()
	c2, _ := net.Dial(target[0], target[1])
	io.Copy(c2, c1)
}
func inputOnly(c1 net.Conn, target []string) {
	defer c1.Close()
	c2, _ := net.Dial(target[0], target[1])
	io.Copy(c1, c2)
}
func bidirectional(c1 net.Conn, target []string) {
	defer c1.Close()
	c2, _ := net.Dial(target[0], target[1])
	go io.Copy(c1, c2)
	io.Copy(c2, c1)
}
