package protocol

type MessageType string

const (
	MessageTypeUndefined MessageType = "undefined"
	MessageTypeRequest               = "request"
	MessageTypeResponse              = "response"
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
