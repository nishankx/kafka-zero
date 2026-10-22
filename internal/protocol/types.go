package protocol

type MessageType byte

const (
	Produce      MessageType = 0x01
	Fetch        MessageType = 0x02
	Metadata     MessageType = 0x03
	CreateTopic MessageType = 0x04
	

	ProduceResponse MessageType = 0x11
	FetchResponse   MessageType = 0x12
)