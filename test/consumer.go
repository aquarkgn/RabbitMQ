package main

import (
	"github.com/streadway/amqp"
	"log"
	"os"
	"os/signal"
	"sync"
)

const (
	rabbitMQURL  = "amqp://guest:guest@localhost:5672/"
	exchangeName = "example_exchange"
	queueName    = "example_queue"
)

func main() {
	// 创建 RabbitMQ 连接
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("无法连接到 RabbitMQ: %v", err)
	}
	log.Printf("已连接到 RabbitMQ: %s", rabbitMQURL)
	defer conn.Close()

	// 创建通道
	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("无法创建通道: %v", err)
	}
	log.Printf("已创建通道")
	defer channel.Close()

	// 声明交换机
	log.Printf("正在声明交换机: %s", exchangeName)
	err = channel.ExchangeDeclare(exchangeName, "fanout", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("无法声明交换机: %v", err)
	}

	// 创建队列
	log.Printf("正在创建队列: %s", queueName)
	_, err = channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("无法创建队列: %v", err)
	}

	// 绑定队列到交换机
	log.Printf("正在绑定队列到交换机 %s %s", queueName, exchangeName)
	err = channel.QueueBind(queueName, "", exchangeName, false, nil)
	if err != nil {
		log.Fatalf("无法绑定队列到交换机: %v", err)
	}

	// 使用 WaitGroup 来等待消费者 goroutine 完成
	var wg sync.WaitGroup

	// 启动消费者
	wg.Add(1)
	log.Printf("正在启动消费者")
	go consumeMessage(channel, &wg)

	// 等待中断信号，优雅地关闭连接
	log.Printf("按 CTRL+C 退出")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt

	log.Printf("正在关闭通道")
	// 等待消费者 goroutine 完成
	wg.Wait()
}

func consumeMessage(channel *amqp.Channel, wg *sync.WaitGroup) {
	defer wg.Done()

	msgs, err := channel.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		log.Printf("无法启动消费者: %v", err)
		return
	}

	for msg := range msgs {
		log.Printf("收到消息: %s", string(msg.Body))
	}
}
