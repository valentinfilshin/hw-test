package main

import (
	"bytes"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		// Начинаем слушать какой-то свободный порт
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		// Клиент
		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		// Сервер - принимаем сообщение от порта и отвечаем ему
		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})
}

func TestTelnetClientConnectionTimeout(t *testing.T) {
	client := NewTelnetClient("ya.ru:666", 2*time.Second, io.NopCloser(&bytes.Buffer{}), &bytes.Buffer{})
	start := time.Now()
	err := client.Connect()
	require.Error(t, err)
	require.GreaterOrEqual(t, time.Since(start), 2*time.Second)
}

func TestTelnetClientClose(t *testing.T) {
	client := NewTelnetClient("ya.ru:80", 200*time.Millisecond, io.NopCloser(&bytes.Buffer{}), &bytes.Buffer{})
	require.NoError(t, client.Connect())
	require.NoError(t, client.Close())
}

func TestTelnetClientSend(t *testing.T) {
	client := NewTelnetClient("ya.ru:80", 200*time.Millisecond, io.NopCloser(&bytes.Buffer{}), &bytes.Buffer{})
	require.NoError(t, client.Connect())
	defer func() { require.NoError(t, client.Close()) }()
	require.NoError(t, client.Send())
}

func TestNewTelnetClientServerCloseConnection(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:")
	require.NoError(t, err)
	defer func() { require.NoError(t, l.Close()) }()

	go func() {
		conn, err := l.Accept()
		require.NoError(t, err)
		require.NoError(t, conn.Close())
	}()

	out := &bytes.Buffer{}
	client := NewTelnetClient(l.Addr().String(), 200*time.Millisecond, io.NopCloser(&bytes.Buffer{}), out)

	require.NoError(t, client.Connect())
	time.Sleep(2 * time.Second)
	require.NoError(t, client.Send())
	require.NoError(t, client.Receive())
	require.NoError(t, client.Send())
}
