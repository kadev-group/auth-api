package rabbitmq

import (
	"auth-api/internal/models"
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type ProducerClient struct {
	chanel *amqp.Channel
	queues map[string]amqp.Queue
}

func NewProducerClient(config *models.Config, log *zap.Logger) *ProducerClient {
	log = log.Named("[RabbitMQ]")

	conn, err := amqp.Dial(config.RabbitMQ.ServerURL)
	if err != nil {
		log.Fatal(fmt.Sprintf("connect: %s", err))
	}

	chanel, err := conn.Channel()
	if err != nil {
		log.Fatal(fmt.Sprintf("open channel: %s", err))
	}

	return &ProducerClient{
		chanel: chanel,
		queues: map[string]amqp.Queue{},
	}
}

func (pc *ProducerClient) Send(ctx context.Context, qName string, message []byte, args ...amqp.Table) (err error) {
	q, ok := pc.queues[qName]
	if !ok {
		var arg amqp.Table
		if len(args) != 0 {
			arg = args[0]
		}

		q, err = pc.chanel.QueueDeclare(qName,
			true, false, false, false, arg)
		if err != nil {
			return err
		}
	}

	err = pc.chanel.PublishWithContext(ctx, "", q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		})
	if err != nil {
		return err
	}
	return nil
}
