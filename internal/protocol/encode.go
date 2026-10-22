package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
)

func EncodeProduceRequest(topic string, partition int, message []byte) ([]byte, error) {
	topicBytes := []byte(topic)

	if len(topicBytes) > math.MaxUint16 {
		return nil, errors.New("topic name too long")
	}

	if len(message) > math.MaxUint32 {
		return nil, errors.New("message too long")
	}

	var buf bytes.Buffer

	buf.WriteByte(byte(Produce))

	if err := binary.Write(&buf, binary.BigEndian, uint16(len(topicBytes)));
	err != nil {
		return nil, fmt.Errorf("encode topic length: %w", err)
	}

	buf.Write(topicBytes)

	if err := binary.Write(&buf, binary.BigEndian, partition);
	err != nil {
		return nil, fmt.Errorf("encode partition: %w", err)
	}

	if err := binary.Write(&buf, binary.BigEndian, int32(len(message)));
	err != nil {
		return nil, fmt.Errorf("encode message length: %w", err)
	}

	buf.Write(message)

	return buf.Bytes(), nil
}