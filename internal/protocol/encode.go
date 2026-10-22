package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
)

func EncodeProduceRequest(topic string, partition int32, message []byte) ([]byte, error) {
	topicBytes := []byte(topic)

	if len(topicBytes) > math.MaxUint16 {
		return nil, errors.New("topic name too long")
	}

	if len(message) > math.MaxUint32 {
		return nil, errors.New("message too long")
	}

	var buf bytes.Buffer

	buf.WriteByte(byte(Produce))

	if err := binary.Write(&buf, binary.BigEndian, uint16(len(topicBytes))); err != nil {
		return nil, fmt.Errorf("encode topic length: %w", err)
	}

	buf.Write(topicBytes)

	if err := binary.Write(&buf, binary.BigEndian, partition); err != nil {
		return nil, fmt.Errorf("encode partition: %w", err)
	}

	if err := binary.Write(&buf, binary.BigEndian, int32(len(message))); err != nil {
		return nil, fmt.Errorf("encode message length: %w", err)
	}

	buf.Write(message)

	return buf.Bytes(), nil
}

func EncodeFetchRequest(topic string, partition int32, offset, maxBytes int32) ([]byte, error) {
	//Purpose: Encodes a request to fetch (read) messages from a specific topic, partition, starting from a given offset
	//Key parameters: Topic, partition, starting offset, and maximum bytes to fetch
	//Structure: Similar to produce request but includes offset and maxBytes parameters

	var buf bytes.Buffer

	buf.WriteByte(byte(Fetch))

	if err := binary.Write(&buf, binary.BigEndian, uint16(len(topic))); err != nil {
		return nil, fmt.Errorf("encode topic length: %w", err)
	}

	buf.Write([]byte(topic))

	if err := binary.Write(&buf, binary.BigEndian, partition); err != nil {
		return nil, fmt.Errorf("encode partition: %w", err)
	}

	if err := binary.Write(&buf, binary.BigEndian, offset); err != nil {
		return nil, fmt.Errorf("encode offset: %w", err)
	}

	if err := binary.Write(&buf, binary.BigEndian, maxBytes); err != nil {
		return nil, fmt.Errorf("encode maxBytes: %w", err)
	}

	return buf.Bytes(), nil
}

func EncodeMetadataRequest() ([]byte, error) {
	// Purpose: Encodes a request to retrieve metadata about brokers and topics in the cluster
	// Structure: Simplest request — just the request type byte (0x03)

	var buf bytes.Buffer

	buf.WriteByte(byte(Metadata))

	return buf.Bytes(), nil
}

func EncodeCreateTopicRequest(topic string, numPartitions int32, replicationFactor int16) ([]byte, error) {
	//Purpose: Encodes a request to create a new topic with specified partitions and replication
	// Key parameters: Topic name, number of partitions, and replication factor

	topicBytes := []byte(topic)

	if len(topicBytes) > math.MaxUint16 {
		return nil, errors.New("topic name too long")
	}

	var buf bytes.Buffer

	buf.WriteByte(byte(CreateTopic))

	if err := binary.Write(&buf, binary.BigEndian, uint16(len(topicBytes))); err != nil {
		return nil, fmt.Errorf("encode topic length: %w", err)
	}

	buf.Write(topicBytes)

	if err := binary.Write(&buf, binary.BigEndian, numPartitions); err != nil {
		return nil, fmt.Errorf("encode number of partitions: %w", err)
	}

	if err := binary.Write(&buf, binary.BigEndian, replicationFactor); err != nil {
		return nil, fmt.Errorf("encode replication factor: %w", err)
	}

	return buf.Bytes(), nil
}

func EncodeReplicateRequest(
	topic string,
	partition int32,
	offset int64,
	message []byte,
) ([]byte, error) {
	topicBytes := []byte(topic)

	if len(topicBytes) > math.MaxUint16 {
		return nil, errors.New("topic name too long")
	}
	if partition < 0 {
		return nil, errors.New("partition cannot be negative")
	}
	if offset < 0 {
		return nil, errors.New("offset cannot be negative")
	}
	if len(message) > math.MaxInt32 {
		return nil, errors.New("message too large")
	}

	var buf bytes.Buffer

	buf.WriteByte(byte(Replicate))

	if err := binary.Write(&buf, binary.BigEndian, uint16(len(topicBytes))); err != nil {
		return nil, fmt.Errorf("encode topic length: %w", err)
	}

	buf.Write(topicBytes)

	if err := binary.Write(&buf, binary.BigEndian, partition); err != nil {
		return nil, fmt.Errorf("encode partition: %w", err)
	}

	if err := binary.Write(&buf, binary.BigEndian, offset); err != nil {
		return nil, fmt.Errorf("encode offset: %w", err)
	}

	if err := binary.Write(&buf, binary.BigEndian, int32(len(message))); err != nil {
		return nil, fmt.Errorf("encode message length: %w", err)
	}

	buf.Write(message)

	return buf.Bytes(), nil
}

func EncodeTopicNotification(topic string) ([]byte, error) {
	topicBytes := []byte(topic)

	if len(topicBytes) > math.MaxUint16 {
		return nil, errors.New("topic name too long")
	}

	var buf bytes.Buffer

	buf.WriteByte(byte(TopicNotification))

	if err := binary.Write(&buf, binary.BigEndian, uint16(len(topicBytes))); err != nil {
		return nil, fmt.Errorf("encode topic length: %w", err)
	}

	buf.Write(topicBytes)

	return buf.Bytes(), nil
}

func EncodeErrorResponse(message string) ([]byte, error) {
	messageBytes := []byte(message)

	if len(messageBytes) > math.MaxUint16 {
		return nil, errors.New("error message too long")
	}

	var buf bytes.Buffer

	buf.WriteByte(byte(ErrorResponse))

	if err := binary.Write(
		&buf,
		binary.BigEndian,
		uint16(len(messageBytes)),
	); err != nil {
		return nil, fmt.Errorf("encode error message length: %w", err)
	}

	buf.Write(messageBytes)

	return buf.Bytes(), nil
}
