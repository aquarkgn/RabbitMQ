package service

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/streadway/amqp"

	"rabbitmq/internal/app"
)

func TestPublishMessage(t *testing.T) {
	// 创建 RabbitMQ 连接
	conn, err := amqp.Dial(app.RabbitMQURL)
	if err != nil {
		t.Fatalf("无法连接到 RabbitMQ: %v", err)
	}
	log.Printf("已连接到 RabbitMQ: %s", app.RabbitMQURL)
	defer conn.Close()

	// 创建通道
	channel, err := conn.Channel()
	if err != nil {
		t.Fatalf("无法创建通道: %v", err)
	}
	log.Printf("已创建通道")
	defer channel.Close()

	// 声明交换机
	log.Printf("正在声明交换机: %s", app.ExchangeName)
	err = channel.ExchangeDeclare(app.ExchangeName, "fanout", true, false, false, false, nil)
	if err != nil {
		t.Fatalf("无法声明交换机: %v", err)
	}

	// 推送消息到交换机
	for i := 0; i < 10000; i++ {
		message := fmt.Sprintf("消息 %d", i)
		err = channel.Publish(app.ExchangeName, "", false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
		if err != nil {
			t.Errorf("无法推送消息: %v", err)
		} else {
			t.Logf("已推送消息: %s", message)
		}
		time.Sleep(time.Second)
	}
}
