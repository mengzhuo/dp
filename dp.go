package main

import (
	"flag"
	"io"
	"log"
	"os"
	"time"

	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/pub"
	"github.com/go-mangos/mangos/protocol/sub"
	"github.com/go-mangos/mangos/transport/tcp"
)

var (
	url    = flag.String("u", "tcp://0.0.0.0:9999", "URL that listen or attach")
	topic  = flag.String("t", "", "Topic that publish/subscribe")
	listen = flag.Bool("l", false, "pipeline listen at")
)

type writer struct {
	sock mangos.Socket
}

func (w *writer) Write(p []byte) (n int, err error) {
	err = w.sock.Send(p)
	return len(p), err
}

func publish() {
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
	}
}

func subscribe() {
	sock, err := sub.NewSocket()
	if err != nil {
		log.Fatal(err)
	}
	defer sock.Close()

	sock.AddTransport(tcp.NewTransport())

	if err = sock.Dial(*url); err != nil {
		log.Fatal(err)
	}

	err = sock.SetOption(mangos.OptionSubscribe, []byte(*topic))
	if err != nil {
		log.Fatal(err)
	}

	for {
		var (
			msg []byte
			err error
		)

		if msg, err = sock.Recv(); err != nil {
			log.Fatal(err)
		} else {
			os.Stdout.Write(msg)
		}
	}
}

func main() {
	log.SetPrefix("dp: ")
	flag.Parse()
	if *listen {
		publish()
	} else {
		subscribe()
	}
}
