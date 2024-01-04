package protocol

type MessageType string

const (
	MessageTypeClientWelcome  MessageType = "client_welcome"
	MessageTypeServerQuestion             = "server_question"
	MessageTypeClientAnswer               = "client_answer"
	MessageTypeServerResult               = "server_result"
)

type Message struct {
	Type MessageType
}

type ClientWelcomeRequest struct {
	Message
}

type ServerQuestionRequest struct {
	Message
	Prefix     int64
	Difficulty int
}

type ClientAnswerResponse struct {
	Message
	Nonce      int64
	Prefix     int64
	Difficulty int
}

type ServerResultResponse struct {
	Message
	Success bool
	Payload string
}
