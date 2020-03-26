package main

import (
	"fmt"
	"imooc-rabbitmq/RabbitMQ"
	"strconv"
	"time"
)

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("imoocSimple")
	for i := 0; i <100; i++ {
		rabbitmq.PublishSimple("Hello imooc!" + strconv.Itoa(i))
		time.Sleep(time.Second * 1)
		fmt.Println("send success!")
	}
}
