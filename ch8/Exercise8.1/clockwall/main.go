package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"text/scanner"
	"time"
)

type TimeZoneIP struct {
	addr string
	ip   string
	port string
}

type lexer struct {
	scan  scanner.Scanner
	token rune
}

func (lex *lexer) next() {
	lex.token = lex.scan.Scan()
}

func (lex *lexer) text() string {
	return lex.scan.TokenText()
}

func tokenType(tok rune) string {
	switch tok {
	case scanner.Ident:
		return "Ident"
	case scanner.Int:
		return "Int"
	case scanner.String:
		return "String"
	case scanner.Float:
		return "Float"
	case scanner.RawString:
		return "RawString"
	case scanner.Char:
		return "Char"
	case scanner.Comment:
		return "Comment"
	default:
		return fmt.Sprintf("'%c'", tok)
	}
}

func parse(input string) (*TimeZoneIP, error) {
	// fmt.Println(input)
	res := &TimeZoneIP{}
	lexer := new(lexer)
	lexer.scan.Init(strings.NewReader(input))
	lexer.scan.Mode = scanner.ScanIdents | scanner.ScanInts | scanner.ScanStrings
	lexer.next()
	if lexer.token == scanner.Ident {
		res.addr = lexer.text()
		lexer.next()
	}
	if lexer.token == '=' {
		lexer.next()
	}
	if lexer.token == scanner.Ident {
		res.ip = lexer.text()
		lexer.next()
	}
	if lexer.token == ':' {
		lexer.next()
	}
	if lexer.token == scanner.Int {
		res.port = lexer.text()
		lexer.next()
	}
	if lexer.token != scanner.EOF {
		return nil, fmt.Errorf("invalid input:%s", input)
	}
	return res, nil
}

var timeMap map[string]string
var addrs []string
var lock sync.Mutex

func clockWall(tzi TimeZoneIP) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", tzi.ip, tzi.port))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	s := bufio.NewScanner(conn)
	for s.Scan() {
		lock.Lock()
		timeMap[tzi.addr] = s.Text()
		lock.Unlock()
	}
	fmt.Println(tzi.addr, "done")
}

func showTable() {
	fmt.Println(strings.Join(addrs, "\t"))
	for {
		str := ""
		for _, addr := range addrs {
			str += timeMap[addr]
			str += "\t"
		}
		fmt.Println(str)
		time.Sleep(1 * time.Second)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("%v", fmt.Errorf("len(args) must be greater than %d", 1))
		os.Exit(1)
	}
	args := os.Args[1:]
	timeMap = make(map[string]string)
	for _, arg := range args {
		if tzi, err := parse(arg); err == nil {
			addrs = append(addrs, tzi.addr)
			go clockWall(*tzi)
		}
	}
	go showTable()
	time.Sleep(5 * time.Minute)
}
