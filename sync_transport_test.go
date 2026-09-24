package gadb

import (
	"bytes"
	"errors"
	"io"
	"net"
	"strings"
	"testing"
	"testing/iotest"
)

func TestSendStreamSendsDataBeforeReadError(t *testing.T) {
	readErr := errors.New("read failed")
	for name, wantErr := range map[string]error{"EOF": nil, "read error": readErr} {
		t.Run(name, func(t *testing.T) {
			reader := io.Reader(strings.NewReader("apk"))
			if wantErr != nil {
				reader = io.MultiReader(reader, iotest.ErrReader(wantErr))
			}
			writer, receiver := net.Pipe()
			defer receiver.Close()
			done := make(chan error, 1)
			go func() {
				done <- (syncTransport{sock: writer}).SendStream(iotest.DataErrReader(reader))
				writer.Close()
			}()
			got, err := io.ReadAll(receiver)
			if err != nil {
				t.Fatal(err)
			}
			if err := <-done; !errors.Is(err, wantErr) {
				t.Fatalf("SendStream error = %v, want %v", err, wantErr)
			}
			if want := []byte("DATA\x03\x00\x00\x00apk"); !bytes.Equal(got, want) {
				t.Fatalf("wire bytes = %q, want %q", got, want)
			}
		})
	}
}
