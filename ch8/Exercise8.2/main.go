package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func writeln(c net.Conn, s ...interface{}) {
	s = append(s, "\r\n")
	fmt.Fprint(c, s...)
}

func showDir(c net.Conn, args []string) {
	path := "."
	if len(args) != 0 {
		path = args[0]
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		writeln(c, fmt.Sprintf("%v", err))
	}
	res := ""
	for _, entry := range entries {
		res += entry.Name() + "\t"
	}
	res += "\n"
	writeln(c, res)
}

type downloadFile struct {
	file os.FileInfo
}

func (f *downloadFile) String() string {
	str := ""
	str += fmt.Sprintf("name:%s\n", f.file.Name())
	str += fmt.Sprintf("isDir:%t\n", f.file.IsDir())
	str += fmt.Sprintf("mode:%s\n", f.file.Mode())
	str += fmt.Sprintf("size:%d\n", f.file.Size())
	str += fmt.Sprintf("Modify Time:%s\n", f.file.ModTime())
	return str
}

func getFile(c net.Conn, args []string) {
	if len(args) == 0 {
		writeln(c, "len(args) must be greater than 0,please input a path that you want to get")
	}
	path := args[0]
	info, err := os.Stat(path)
	df := downloadFile{info}
	if err != nil {
		writeln(c, err)
	}
	writeln(c, df.String())
}

func handleFTPConn(c net.Conn) {
	writeln(c, "220 Ready.")
	s := bufio.NewScanner(c)
	for s.Scan() {
		fileds := strings.Fields(s.Text())
		if len(fileds) == 0 {
			continue
		}
		var args []string
		cmd := fileds[0]
		if len(fileds) > 1 {
			args = fileds[1:]
		}
		cmd = strings.ToLower(cmd)
		switch cmd {
		case "ls":
			showDir(c, args)
		case "get":
			getFile(c, args)
		case "close":
			writeln(c, "Bye!")
			c.Close()
			return
		default:
			writeln(c, "it is not cmd")
		}
	}
}

func main() {
	var port int
	flag.IntVar(&port, "port", 8000, "listen port")

	ln, err := net.Listen("tcp4", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal("Opening main listener:", err)
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			log.Print("Accepting new connection:", err)
		}
		go handleFTPConn(c)
	}
}
