package RabbitMQ

import (
	"github.com/streadway/amqp"
	"log"
)

const MQURL = "amqp://imoocuser:imoocuser@127.0.0.1:5672/imooc"

type RabbitMQ struct {
	conn  *amqp.Connection
	channel *amqp.Channel
	//队列名称
	QueueName string
	//交换机
	Exchange string
	key string
	//连接信息
	Mqurl string
}

func NewRabbitMQ(queueName, exchange, key string) *RabbitMQ {
	rabbitmq := &RabbitMQ{
		QueueName: queueName,
		Exchange:  exchange,
		key:       key,
		Mqurl:     MQURL,
	}
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnError(err, "create conn error")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnError(err, "get channel error")
	return rabbitmq
}

func (r *RabbitMQ) Destroy() {
	r.channel.Close()
	r.conn.Close()
}

func (r *RabbitMQ) failOnError(err error, message string)  {
	if err != nil {
		log.Fatalf("%s,%s", message, err)
		return
	}
}

//简单模式step1：创建简单模式下rabbitMQ
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	return NewRabbitMQ(queueName, "", "")
}

//简单模式step2：生产
func (r *RabbitMQ) PublishSimple(message string)  {
	//1.申请队列,不存在则创建,存在则跳过
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false,//是否持久化
		false,
		false ,//排他性
		false, //是否阻塞
		nil)
	r.failOnError(err, "publish error")
	//2.发送消息到队列中
	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		false, //如果为true,根据exchange类型和routkey规则，
		//如果无法找到符合条件的队列那么会把发送的消息返回给发送者
		false,
		amqp.Publishing{
			ContentType:     "text/plain",
			Body:            []byte(message),
		})
}
//简单模式step3：消费
func (r *RabbitMQ)ConsumerSimple()  {
	//1.申请队列,不存在则创建,存在则跳过
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false,//是否持久化
		false,
		false ,//排他性
		false, //是否阻塞
		nil)
	r.failOnError(err, "publish error")
	//2.接受消息
	msgs, err := r.channel.Consume(
		r.QueueName,
		"",//用来区分多个消费者
		true, //是否自动应答
		false,
		false,
		false,
		nil)
	r.failOnError(err, "consume error")
	forever := make(chan bool)
	//3.启用携程消费
	go func() {
		for d := range msgs {
			log.Printf("receive message. %s", d.Body)
		}
	}()
	log.Printf("wait for message, to exit press CTRL+C")
	<-forever
}
//订阅模式创建RabbitMQ实例
func NewRabbitMQPubSub(exchangeName string) *RabbitMQ {
	rabbitmq := NewRabbitMQ("", exchangeName, "")
	return rabbitmq
}
//订阅模式生产
func (r *RabbitMQ) PublishPub(message string)  {
	//1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout", //广播类型
		true,//是否持久化
		false,
		false ,//排他性
		false, //是否阻塞
		nil)
	r.failOnError(err, "declare exchange error.")
	//2.发送消息
	r.channel.Publish(
		r.Exchange,
		"",
		false, //如果为true,根据exchange类型和routkey规则，
		//如果无法找到符合条件的队列那么会把发送的消息返回给发送者
		false,
		amqp.Publishing{
			ContentType:     "text/plain",
			Body:            []byte(message),
		})
}

//订阅模式消费
func (r *RabbitMQ)ConsumerPub()  {
	var err error
	//1.创建交换机
	err = r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout",
		true,//是否持久化
		false,
		false ,
		false, //是否阻塞
		nil)
	r.failOnError(err, "failed to declare exchange error")
	//2.创建队列  注意队列名称不要写
	queue, err := r.channel.QueueDeclare(
		"", //随机生成队列名称
		false,
		false,
		true,
		false,
		nil)
	r.failOnError(err, "failed to declare queue error")
	//3. 绑定队列到exchange中
	err = r.channel.QueueBind(
		queue.Name,
		//在pub/sub模式下，这里key为空
		"",
		r.Exchange,
		false,
		nil)
	//4.消费消息
	msgs, err := r.channel.Consume(
		queue.Name,
		"",//用来区分多个消费者
		true, //是否自动应答
		false,
		false,
		false,
		nil)
	r.failOnError(err, "consume error")
	forever := make(chan bool)
	//3.启用携程消费
	go func() {
		for d := range msgs {
			log.Printf("receive message. %s", d.Body)
		}
	}()
	log.Printf("wait for message, to exit press CTRL+C")
	<-forever
}

//路由模式
//创建RabbitMQ实例
func NewRabbitMQRouting(exchangeName, routingKey string)  *RabbitMQ{
	return NewRabbitMQ("", exchangeName, routingKey)
}

//路由模式发送消息
func (r *RabbitMQ) PublishRouting(message string)  {
	//1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct", //路由模式要设置为direct
		true,
		false,
		false,
		false,
		nil)
	r.failOnError(err, "declare exchange failed.")
	//2.发送消息
	err = r.channel.Publish(
		r.Exchange,
		r.key,
		false,
		false,
		amqp.Publishing{
			ContentType:     "text/plain",
			Body:            []byte(message),
		})
	r.failOnError(err, "publish failed.")
}

func (r *RabbitMQ)ConsumerRouting()  {
	//1.试探性创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil)
	r.failOnError(err, "declare exchange failed.")
	//2.试探性创建队列，队列名称不要写
	queue, err := r.channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil)
	//3.绑定队列到交换机
	err = r.channel.QueueBind(
		queue.Name,
		r.key,
		r.Exchange,
		false,
		nil)
	r.failOnError(err, "QueueBind failed.")
	//4.消费消息
	msgs, err := r.channel.Consume(
		queue.Name,
		"",//用来区分多个消费者
		true, //是否自动应答
		false,
		false,
		false,
		nil)
	r.failOnError(err, "consume error")
	forever := make(chan bool)
	//3.启用携程消费
	go func() {
		for d := range msgs {
			log.Printf("receive message. %s", d.Body)
		}
	}()
	log.Printf("wait for message, to exit press CTRL+C")
	<-forever
}


