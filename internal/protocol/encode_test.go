package protocol

import (
	"bytes"
	"testing"
)

func TestEncodeProduceRequest(t *testing.T) {
	got, err := EncodeProduceRequest(
		"orders",
		0,
		[]byte("Hello"),
	)
	if err != nil {
		t.Fatalf("EncodeProduceRequest() error: %v", err)
	}

	expected := []byte{
		0x01,
		0x00, 0x06,
		'o', 'r', 'd', 'e', 'r', 's',
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x05,
		'H', 'e', 'l', 'l', 'o',
	}

	if !bytes.Equal(got, expected) {
		t.Fatalf("encoded bytes mismatch\n got: %v\nwant: %v", got, expected)
	}
}
