package service

import (
	"log"
	"sync"

	"github.com/streadway/amqp"

	"rabbitmq/internal/app"
)

var ConsumerInstance = new(Consumer)

type Consumer struct{}

func (consumer *Consumer) Message(channel *amqp.Channel, wg *sync.WaitGroup) {
	defer wg.Done()

	msgs, err := channel.Consume(app.QueueName, "", true, false, false, false, nil)
	if err != nil {
		log.Printf("无法启动消费者: %v", err)
		return
	}
	log.Printf("已启动消费者")

	for msg := range msgs {
		log.Printf("收到消息: %s", string(msg.Body))
	}
}
