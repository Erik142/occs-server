package message

type MessageType string

const (
	Subscribe   MessageType = "subscribe"
	Unsubscribe             = "unsubscribe"
	Publish                 = "publish"
)

type Message struct {
	Type MessageType
	Data string
}
