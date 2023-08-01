package notify

import (
	"encoding/json"
	"github.com/nsqio/go-nsq"
)

const DEVICE_BREAK_NOTIFY = "device_break_notify"

var producer *nsq.Producer

func init() {
	config := nsq.NewConfig()
	var err error
	producer, err = nsq.NewProducer("127.0.0.1:4150", config)
	if err != nil {
		panic(err)
	}
}

func sendMessage(topic string, msg *MsgBase) error {
	byteContent, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return producer.Publish(topic, byteContent)
}

type MsgBase struct {
	Source    string `json:"source"`
	Timestamp int64  `json:"timestamp"`
}
