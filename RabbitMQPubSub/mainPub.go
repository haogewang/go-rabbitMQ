package main

import (
	"imooc-rabbitmq/RabbitMQ"
	"strconv"
	"time"
)

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQPubSub("e")
	for i := 0; i < 100; i++ {
		rabbitmq.PublishPub("hello pub sub" + strconv.Itoa(i))
		time.Sleep(time.Second * 1)
	}
}
