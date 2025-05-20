package main

import (
	"io"
	"net"
	"time"
)

type Telnet struct {
	Address string
	Timeout time.Duration
	Conn    net.Conn
	In      io.ReadCloser
	Out     io.Writer
}

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Telnet{
		Address: address,
		Timeout: timeout,
		In:      in,
		Out:     out,
	}
}

func (t *Telnet) Connect() error {
	conn, err := net.DialTimeout("tcp", t.Address, t.Timeout)
	if err != nil {
		return err
	}

	t.Conn = conn
	return nil
}

func (t *Telnet) Close() error {
	return t.Conn.Close()
}

func (t *Telnet) Send() error {
	_, err := io.Copy(t.Conn, t.In)
	return err
}

func (t *Telnet) Receive() error {
	_, err := io.Copy(t.Out, t.Conn)
	return err
}
