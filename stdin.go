package main

import (
	"flag"
	"io"
	"log"
	"os"
	"time"

	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/pub"
	"github.com/go-mangos/mangos/transport/tcp"
)

var (
	url   = flag.String("u", "tcp://127.0.0.1:9999", "URL that listen or attach")
	topic = flag.String("t", "dp", "Topic that publish/subscribe")
)

type writer struct {
	sock mangos.Socket
}

func (w *writer) Write(p []byte) (n int, err error) {
	err = w.sock.Send(p)
	return len(p), err
}

func main() {
	sock, err := pub.NewSocket()
	if err != nil {
		log.Fatal(err)
	}
	defer sock.Close()
	sock.AddTransport(tcp.NewTransport())

	if err = sock.Listen(*url); err != nil {
		log.Fatalf("%s %s ", err, *url)
	}

	w := &writer{sock}

	for {
		n, err := io.Copy(w, os.Stdin)
		if err != nil {
			log.Fatal("err:%s", err)
		}
		if n == 0 {
			time.Sleep(time.Millisecond * 500)
		}
		log.Printf("Write:%d", n)
	}
}
