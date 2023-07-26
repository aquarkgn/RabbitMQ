package main

import (
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/streadway/amqp"

	"rabbitmq/cmd/service"

	"rabbitmq/internal/app"
)

func main() {
	// 创建 RabbitMQ 连接
	conn, err := amqp.Dial(app.RabbitMQURL)
	if err != nil {
		log.Fatalf("无法连接到 RabbitMQ: %v", err)
	}
	log.Printf("已连接到 RabbitMQ: %s", app.RabbitMQURL)
	defer conn.Close()

	// 创建通道
	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("无法创建通道: %v", err)
	}
	log.Printf("已创建通道")
	defer channel.Close()

	// 声明交换机
	log.Printf("正在声明交换机: %s", app.ExchangeName)
	err = channel.ExchangeDeclare(app.ExchangeName, "fanout", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("无法声明交换机: %v", err)
	}

	// 创建队列
	log.Printf("正在创建队列: %s", app.QueueName)
	_, err = channel.QueueDeclare(app.QueueName, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("无法创建队列: %v", err)
	}

	// 绑定队列到交换机
	log.Printf("正在绑定队列到交换机 %s %s", app.QueueName, app.ExchangeName)
	err = channel.QueueBind(app.QueueName, "", app.ExchangeName, false, nil)
	if err != nil {
		log.Fatalf("无法绑定队列到交换机: %v", err)
	}

	// 使用 WaitGroup 来等待消费者 goroutine 完成
	var wg sync.WaitGroup

	// 启动消费者
	wg.Add(1)
	log.Printf("正在启动消费者")
	go service.ConsumerInstance.Message(channel, &wg)

	// 等待中断信号，优雅地关闭连接
	log.Printf("按 CTRL+C 退出")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt

	log.Printf("正在关闭通道")
	// 等待消费者 goroutine 完成
	wg.Wait()
}
