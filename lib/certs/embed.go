package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
)

//go:generate go run embed.go

var head = `package certs

var Path = "/etc/ssl/certs/ca-certificates.crt"
var Contents = 
`

func main() {
	var err error
	if err = os.MkdirAll("content", 0775); nil != err {
		log.Fatalln(err)
	}

	var contents []byte
	if contents, err = ioutil.ReadFile("ca-certificates.crt"); nil != err {
		log.Fatalln(err)
	}

	parts := [][]byte{
		[]byte(head),
		[]byte("`"),
		contents,
		[]byte("\n`"),
	}

	contents = bytes.Join(parts, []byte(""))

	if ioutil.WriteFile("content/content.gen.go", contents, 0666); nil != err {
		log.Fatalln(err)
	}
}
