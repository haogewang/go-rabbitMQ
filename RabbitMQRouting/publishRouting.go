package main

import (
	"imooc-rabbitmq/RabbitMQ"
	"strconv"
	"time"
)

func main() {
	imoocOne := RabbitMQ.NewRabbitMQRouting("exImooc", "imooc_one")
	imoocTwo := RabbitMQ.NewRabbitMQRouting("exImooc", "imooc_two")
	for i :=0; i < 100; i++{
		imoocOne.PublishRouting("hello one" + strconv.Itoa(i))
		imoocTwo.PublishRouting("hello two" + strconv.Itoa(i))
		time.Sleep(time.Second * 1)
	}
}
