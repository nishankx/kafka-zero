package protocol

type MessageType byte

const (
	Produce     MessageType = 0x01
	Fetch       MessageType = 0x02
	Metadata    MessageType = 0x03
	CreateTopic MessageType = 0x04

	ProduceResponse MessageType = 0x11
	FetchResponse   MessageType = 0x12
	MetadataResponse MessageType = 0x13
)

type ProduceResult struct {
	Offset int64
	Err    error
}

type Record struct {
	Offset  int64
	Message []byte
}

type FetchResult struct {
	Records []Record
	Err     error
}



type BrokerMetadata struct {
	ID   int32
	Host string
	Port int32
}

type PartitionMetadata struct {
	ID        int32
	LeaderID  int32
	Replicas  []int32
}

type TopicMetadata struct {
	Name       string
	Partitions []PartitionMetadata
}

type MetadataResult struct {
	Brokers []BrokerMetadata
	Topics  []TopicMetadata
	Err     error
}
