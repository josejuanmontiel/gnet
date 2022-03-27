package test

import (
	"crypto/rand"
	"io"
	"net"
	"testing"
	"time"
)

func TestHandleMessage(t *testing.T) {
	server, client := net.Pipe()

	// Set deadline so test can detect if HandleMessage does not return
	client.SetDeadline(time.Now().Add(time.Second))

	// Configure a go routine to act as the server
	go func() {
		HandleMessage(server)
		server.Close()
	}()

	_, err := client.Write([]byte("test\n"))
	if err != nil {
		t.Fatalf("failed to write: %s", err)
	}

	// As the go routine closes the connection ReadAll is a simple way to get the response
	in, err := io.ReadAll(client)
	if err != nil {
		t.Fatalf("failed to read: %s", err)
	}

	// Using an Assert here will also work (if using a library that provides that functionality)
	if string(in) != "test whatever\n" {
		t.Fatalf("expected `test` got `%s`", in)
	}

	client.Close()
}

func TestDecode(t *testing.T) {
	server, client := net.Pipe()

	// Set deadline so test can detect if HandleMessage does not return
	client.SetDeadline(time.Now().Add(60 * time.Second))

	codec := *&SimpleCodec{}

	// TODO
	// Build Error: go test -c -o /home/jose/git/gnet/test/simple_protocol/protocol/__debug_bin -gcflags all=-N -l .
	// # protocol [protocol.test]
	// ./protocol_test.go:52:15: cannot use server (type net.Conn) as type gnet.Conn in argument to codec.Decode:
	// 	net.Conn does not implement gnet.Conn (missing Context method) (exit status 2)

	// Configure a go routine to act as the server
	go func() {
		codec.Decode(server)
		server.Close()
	}()

	req := make([]byte, 1024)
	_, err := rand.Read(req)
	buf, _ := codec.Encode(req)

	_, err2 := client.Write(buf)
	if err2 != nil {
		t.Fatalf("failed to write: %s", err2)
	}

	// As the go routine closes the connection ReadAll is a simple way to get the response
	in, err := io.ReadAll(client)
	if err != nil {
		t.Fatalf("failed to read: %s", err)
	}

	// Using an Assert here will also work (if using a library that provides that functionality)
	if string(in) != "test whatever\n" {
		t.Fatalf("expected `test` got `%s`", in)
	}

	client.Close()
}
