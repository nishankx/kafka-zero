package protocol

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestDecodeProduceResponseSuccess(t *testing.T) {
	data := make([]byte, 13)
	data[0] = byte(ProduceResponse)

	binary.BigEndian.PutUint64(data[1:9], 42)
	binary.BigEndian.PutUint16(data[9:11], 0)
	binary.BigEndian.PutUint16(data[11:13], 0)

	result, err := DecodeProduceResponse(data)
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}

	if result.Offset != 42 {
		t.Errorf("offset: got %d, want 42", result.Offset)
	}

	if result.Err != nil {
		t.Errorf("unexpected broker error: %v", result.Err)
	}
}

func TestDecodeProduceResponseBrokerError(t *testing.T) {
	message := []byte("unknown topic")
	data := make([]byte, 13+len(message))

	data[0] = byte(ProduceResponse)
	binary.BigEndian.PutUint64(data[1:9], 0)
	binary.BigEndian.PutUint16(data[9:11], 1)
	binary.BigEndian.PutUint16(data[11:13], uint16(len(message)))
	copy(data[13:], message)

	result, err := DecodeProduceResponse(data)
	if err != nil {
		t.Fatalf("decoding should succeed: %v", err)
	}

	if result.Err == nil {
		t.Fatal("expected broker error, got nil")
	}
}

func TestDecodeFetchResponseSuccess(t *testing.T) {
	var buf bytes.Buffer

	buf.WriteByte(byte(FetchResponse))
	_ = binary.Write(&buf, binary.BigEndian, uint16(0)) // Error code
	_ = binary.Write(&buf, binary.BigEndian, uint16(0)) // Error-message length
	_ = binary.Write(&buf, binary.BigEndian, uint32(2)) // Record count

	// Record 1
	_ = binary.Write(&buf, binary.BigEndian, int64(10))
	_ = binary.Write(&buf, binary.BigEndian, int32(5))
	buf.WriteString("Hello")

	// Record 2
	_ = binary.Write(&buf, binary.BigEndian, int64(11))
	_ = binary.Write(&buf, binary.BigEndian, int32(5))
	buf.WriteString("World")

	result, err := DecodeFetchResponse(buf.Bytes())
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}

	if result.Err != nil {
		t.Fatalf("unexpected broker error: %v", result.Err)
	}

	if len(result.Records) != 2 {
		t.Fatalf("got %d records, want 2", len(result.Records))
	}

	if result.Records[0].Offset != 10 ||
		string(result.Records[0].Message) != "Hello" {
		t.Errorf("unexpected first record: %+v", result.Records[0])
	}

	if result.Records[1].Offset != 11 ||
		string(result.Records[1].Message) != "World" {
		t.Errorf("unexpected second record: %+v", result.Records[1])
	}
}

func TestDecodeFetchResponseBrokerError(t *testing.T) {
	var buf bytes.Buffer

	buf.WriteByte(byte(FetchResponse))
	_ = binary.Write(&buf, binary.BigEndian, uint16(1))
	_ = binary.Write(&buf, binary.BigEndian, uint16(len("unknown topic")))
	buf.WriteString("unknown topic")
	_ = binary.Write(&buf, binary.BigEndian, uint32(0))

	result, err := DecodeFetchResponse(buf.Bytes())
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}

	if result.Err == nil {
		t.Fatal("expected broker error")
	}

	if len(result.Records) != 0 {
		t.Fatalf("got %d records, want 0", len(result.Records))
	}
}

func TestEncodeReplicateRequest(t *testing.T) {
	got, err := EncodeReplicateRequest("orders", 1, 42, []byte("Hello"))
	if err != nil {
		t.Fatalf("encode failed: %v", err)
	}

	var expected bytes.Buffer
	expected.WriteByte(byte(Replicate))
	_ = binary.Write(&expected, binary.BigEndian, uint16(6))
	expected.WriteString("orders")
	_ = binary.Write(&expected, binary.BigEndian, int32(1))
	_ = binary.Write(&expected, binary.BigEndian, int64(42))
	_ = binary.Write(&expected, binary.BigEndian, int32(5))
	expected.WriteString("Hello")

	if !bytes.Equal(got, expected.Bytes()) {
		t.Fatalf("encoded bytes mismatch\ngot: %v\nwant: %v", got, expected.Bytes())
	}
}

func TestEncodeTopicNotification(t *testing.T) {
	got, err := EncodeTopicNotification("orders")
	if err != nil {
		t.Fatalf("encode failed: %v", err)
	}

	expected := []byte{
		byte(TopicNotification),
		0x00, 0x06,
		'o', 'r', 'd', 'e', 'r', 's',
	}

	if !bytes.Equal(got, expected) {
		t.Fatalf("encoded bytes mismatch\ngot: %v\nwant: %v", got, expected)
	}
}
