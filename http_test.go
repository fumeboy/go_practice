package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/http"
	"testing"
)

func TestHTTP(t *testing.T) {
	ln, err := net.Listen("tcp", ":8002")
	if err != nil{
		panic(err)
	}
	tl := ln.(*net.TCPListener)
	for {
		if conn, err := tl.Accept();err == nil{
			var r = make([]byte, 1024)
			//var w = make([]byte, 1024)
			conn.Read(r)
			req, _ := http.ReadRequest(bufio.NewReader(bytes.NewReader(r)))
			fmt.Println(req)
			req.Write(bytes.NewBuffer([]byte("123")))
			conn.Write(req.)
		}
	}
}
