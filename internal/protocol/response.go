package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const maxMessageSize = 10 * 1024 * 1024

func DecodeProduceResponse(data []byte) (ProduceResult, error) {
	var result ProduceResult
	r := bytes.NewReader(data)

	msgType, err := r.ReadByte()
	if err != nil {
		return result, fmt.Errorf("error reading message type: %w", err)
	}
	if MessageType(msgType) != ProduceResponse {
		return result, fmt.Errorf("expected ProduceResponse, got %v", msgType)
	}

	if err := binary.Read(r, binary.BigEndian, &result.Offset); err != nil {
		return result, fmt.Errorf("error reading offset: %w", err)
	}

	var errorCode uint16
	if err := binary.Read(r, binary.BigEndian, &errorCode); err != nil {
		return result, fmt.Errorf("error reading error code: %w", err)
	}

	var messageLen uint16
	if err := binary.Read(r, binary.BigEndian, &messageLen); err != nil {
		return result, fmt.Errorf("read error-message length: %w", err)
	}

	message := make([]byte, int(messageLen))
	if _, err := io.ReadFull(r, message); err != nil {
		return result, fmt.Errorf("error reading message: %w", err)
	}

	if r.Len() != 0 {
		return ProduceResult{}, errors.New("unexpected trailing data in response")
	}

	if errorCode != 0 {
		result.Err = fmt.Errorf(
			"broker produced error code %d: %s",
			errorCode,
			string(message),
		)
	} else if messageLen != 0 {
		return ProduceResult{}, errors.New("successful response contains an error message")
	}
	return result, nil
}

func DecodeFetchResponse(data []byte) (FetchResult, error) {

	var result FetchResult
	r := bytes.NewReader(data)

	msgType, err := r.ReadByte()
	if err != nil {
		return result, fmt.Errorf("error reading message type: %w", err)
	}
	if MessageType(msgType) != FetchResponse {
		return result, fmt.Errorf("expected FetchResponse, got %v", msgType)
	}

	var errorCode uint16
	if err := binary.Read(r, binary.BigEndian, &errorCode); err != nil {
		return result, fmt.Errorf("error reading error code: %w", err)
	}

	var errorMessageLen uint16
	if err := binary.Read(r, binary.BigEndian, &errorMessageLen); err != nil {
		return result, fmt.Errorf("error reading error-message length: %w", err)
	}

	errorMessage := make([]byte, int(errorMessageLen))
	if _, err := io.ReadFull(r, errorMessage); err != nil {
		return result, fmt.Errorf("error reading error-message: %w", err)
	}

	var recordCount uint32
	if err := binary.Read(r, binary.BigEndian, &recordCount); err != nil {
		return result, fmt.Errorf("error reading number of records: %w", err)
	}

	if uint64(recordCount) > uint64(r.Len()/12) {
		return result, errors.New("record count exceeds available response data")
	}

	if errorCode != 0 {
		if recordCount != 0 {
			return result, errors.New("error occurred but records are present")
		}
		if r.Len() != 0 {
			return result, errors.New("error occurred but trailing data is present")
		}
		result.Err = fmt.Errorf(
			"broker produced error code %d: %s",
			errorCode,
			string(errorMessage),
		)
		return result, nil
	}

	if errorMessageLen != 0 {
		return result, errors.New("successful response contains an error message")
	}

	result.Records = make([]Record, 0, int(recordCount))
	for i := uint32(0); i < recordCount; i++ {
		var record Record

		if err := binary.Read(r, binary.BigEndian, &record.Offset); err != nil {
			return result, fmt.Errorf("error reading record offset #%d: %w", i, err)
		}

		var messageLen uint32
		if err := binary.Read(r, binary.BigEndian, &messageLen); err != nil {
			return FetchResult{}, fmt.Errorf("read record #%d message length: %w", i, err)
		}

		if messageLen < 0 || messageLen > maxMessageSize {
			return FetchResult{}, fmt.Errorf("read record %d message length %d: %w", i, messageLen, err)
		}

		record.Message = make([]byte, int(messageLen))
		if _, err := io.ReadFull(r, record.Message); err != nil {
			return FetchResult{}, fmt.Errorf("read record %d message %d bytes: %w", i, messageLen, err)
		}

		result.Records = append(result.Records, record)
	}

	if r.Len() != 0 {
		return FetchResult{}, errors.New("unexpected trailing data in response")
	}

	return result, nil
}


func DecodeMetadataResponse(data []byte) (MetadataResult, error) {
	var result MetadataResult
	r := bytes.NewReader(data)

	readUint16 := func() (uint16, error) {
		var v uint16
		err := binary.Read(r, binary.BigEndian, &v)
		return v, err
	}

	readInt32 := func() (int32, error) {
		var v int32
		err := binary.Read(r, binary.BigEndian, &v)
		return v, err
	}

	readString := func() (string, error) {
		n, err := readUint16()
		if err != nil {
			return "", err
		}

		b := make([]byte, int(n))
		if _, err := io.ReadFull(r, b); err != nil {
			return "", err
		}

		return string(b), nil
	}

	msgType, err := r.ReadByte()
	if err != nil {
		return result, fmt.Errorf("read response type: %w", err)
	}
	if MessageType(msgType) != MetadataResponse {
		return result, fmt.Errorf("unexpected response type: 0x%02x", msgType)
	}

	brokerCount, err := readUint16()
	if err != nil {
		return result, fmt.Errorf("read broker count: %w", err)
	}

	result.Brokers = make([]BrokerMetadata, 0, int(brokerCount))

	for i := 0; i < int(brokerCount); i++ {
		id, err := readInt32()
		if err != nil {
			return MetadataResult{}, fmt.Errorf("read broker %d ID: %w", i, err)
		}

		host, err := readString()
		if err != nil {
			return MetadataResult{}, fmt.Errorf("read broker %d host: %w", i, err)
		}

		port, err := readInt32()
		if err != nil {
			return MetadataResult{}, fmt.Errorf("read broker %d port: %w", i, err)
		}

		result.Brokers = append(result.Brokers, BrokerMetadata{
			ID: id, Host: host, Port: port,
		})
	}

	topicCount, err := readUint16()
	if err != nil {
		return MetadataResult{}, fmt.Errorf("read topic count: %w", err)
	}

	result.Topics = make([]TopicMetadata, 0, int(topicCount))

	for i := 0; i < int(topicCount); i++ {
		name, err := readString()
		if err != nil {
			return MetadataResult{}, fmt.Errorf("read topic %d name: %w", i, err)
		}

		partitionCount, err := readUint16()
		if err != nil {
			return MetadataResult{}, fmt.Errorf("read topic %q partition count: %w", name, err)
		}

		topic := TopicMetadata{
			Name:       name,
			Partitions: make([]PartitionMetadata, 0, int(partitionCount)),
		}

		for j := 0; j < int(partitionCount); j++ {
			partitionID, err := readInt32()
			if err != nil {
				return MetadataResult{}, fmt.Errorf("read partition ID: %w", err)
			}

			leaderID, err := readInt32()
			if err != nil {
				return MetadataResult{}, fmt.Errorf("read partition leader: %w", err)
			}

			topic.Partitions = append(topic.Partitions, PartitionMetadata{
				ID: partitionID, LeaderID: leaderID,
			})
		}

		result.Topics = append(result.Topics, topic)
	}

	if r.Len() != 0 {
		return MetadataResult{}, errors.New("unexpected trailing data in metadata response")
	}

	return result, nil
}