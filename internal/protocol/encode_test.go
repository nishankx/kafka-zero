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

func TestEncodeErrorResponse(t *testing.T) {
	got, err := EncodeErrorResponse("unknown topic")
	if err != nil {
		t.Fatalf("encode failed: %v", err)
	}

	expected := []byte{
		byte(ErrorResponse),
		0x00, 0x0D,
		'u', 'n', 'k', 'n', 'o', 'w', 'n', ' ',
		't', 'o', 'p', 'i', 'c',
	}

	if !bytes.Equal(got, expected) {
		t.Fatalf("got %v, want %v", got, expected)
	}
}
