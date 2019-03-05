package main

import (
	"bufio"
	"encoding/base64"
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
	"strings"
)

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, "usage: base64 text...\n\tor from pipe: echo \"text\" | base64\n")
}

func fromArgs(w io.Writer) (err error) {
	str := strings.Join(os.Args[1:], "")
	enc := base64.NewEncoder(base64.StdEncoding, w)
	defer func() { err = enc.Close() }()
	_, err = enc.Write([]byte(str))
	return
}

func fromReader(r io.Reader, w io.Writer) (err error) {
	sc := bufio.NewScanner(r)
	enc := base64.NewEncoder(base64.StdEncoding, w)
	defer func() { err = enc.Close() }()
	for sc.Scan() {
		_, err = enc.Write(sc.Bytes())
	}
	return
}

func main() {
	flag.Usage = usage
	flag.Parse()

	var err error
	switch {
	case terminal.IsTerminal(int(os.Stdin.Fd())) && flag.NArg() == 0:
		usage()
		os.Exit(1)
	case flag.NArg() == 0:
		//_, err := os.Stdin.Stat()
		err = fromReader(os.Stdin, os.Stdout)
	default:
		err = fromArgs(os.Stdout)
	}
	if err != nil {
		fmt.Println(err)
		usage()
		os.Exit(1)
	}
	os.Exit(0)
}
