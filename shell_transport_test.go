package gadb

import (
	"net"
	"testing"
	"time"
)

func TestShellReadClearsHandshakeDeadline(t *testing.T) {
	client, server := net.Pipe()
	defer client.Close()
	defer server.Close()
	// A silent command can outlast the deadline left by its handshake.
	if err := client.SetReadDeadline(time.Now().Add(-time.Second)); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() {
		_, err := server.Write([]byte{byte(shellExit), 1, 0, 0, 0, 0})
		done <- err
	}()
	transport := newShellTransport(client, time.Second)
	command, data, err := transport.Read()
	if err != nil || command != shellExit || len(data) != 1 || data[0] != 0 {
		t.Fatalf("command=%d, data=%v, error=%v", command, data, err)
	}
	if err := <-done; err != nil {
		t.Fatal(err)
	}
}
