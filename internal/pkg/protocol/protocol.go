package protocol

type MessageType string

const (
	MessageTypeRequest  MessageType = "request"
	MessageTypeResponse MessageType = "response"
)

type Message struct {
	Type MessageType
}

type Request struct {
	Message
	Payload []int
}

type Response struct {
	Message
	Payload int
}
