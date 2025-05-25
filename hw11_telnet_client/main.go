package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", time.Second*10, "connection timeout")
}

func main() {
	flag.Parse()

	if flag.NArg() < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [--timeout=10s] <host> <port>", os.Args[0])
		os.Exit(1)
	}

	addr := flag.Arg(0) + ":" + flag.Arg(1)
	client := NewTelnetClient(addr, timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to: %v\n\n", err)
		os.Exit(1)
	}
	defer client.Close()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)

	done := make(chan struct{})

	// Отправитель
	go func() {
		err := client.Send()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to send to %v: %v\n", addr, err)
		}
		fmt.Fprintln(os.Stderr, "EOF. Closing connection")

		done <- struct{}{}
	}()

	// Слушатель
	go func() {
		err := client.Receive()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to receive from %v: %v\n", addr, err)
		}
		fmt.Fprintln(os.Stderr, "Connection closed by remote host")

		done <- struct{}{}
	}()

	select {
	case <-signalChannel:
		fmt.Fprintln(os.Stderr, "SIGINT, exiting")
	case <-done:
	}
}
