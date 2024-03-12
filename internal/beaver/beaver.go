package beaver

import (
	"context"
	"encoding/json"

	"github.com/BelyaevEI/wallet/internal/config"
	"github.com/BelyaevEI/wallet/internal/logger"
	"github.com/BelyaevEI/wallet/internal/models"
	"github.com/BelyaevEI/wallet/internal/store"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Implementation check
var _ Beaverer = Beaver{}

type Beaverer interface {
	RunBeaver() error
	transferFunds(walls []byte) error
}

type Beaver struct {
	log       *logger.Logger
	channel   *amqp.Channel
	queueName string
	store     store.Storer
}

func NewBeaver(log *logger.Logger, cfg config.Config) (Beaverer, error) {

	// Connect to broker messanges
	conn, err := amqp.Dial(cfg.Rabbit)
	if err != nil {
		return nil, err
	}

	// Create channel for message
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// Create a new queue
	q, err := ch.QueueDeclare(
		"test", // name
		true,   // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)

	if err != nil {
		return nil, err
	}

	store, err := store.Connect(cfg)
	if err != nil {
		return nil, err
	}

	return Beaver{log: log, channel: ch, queueName: q.Name, store: store}, nil
}

func (beaver Beaver) RunBeaver() error {

	messages, err := beaver.channel.Consume(
		beaver.queueName, // queue
		"",               // consumer
		true,             // auto-ack
		false,            // exclusive
		false,            // no-local
		false,            // no-wait
		nil,              // args
	)
	if err != nil {
		return err
	}

	for message := range messages {

		go func(mes amqp.Delivery) {
			err := beaver.transferFunds(mes.Body)
			if err != nil {
				beaver.log.Log.Info("transfer funds is failed: ", err)
			}
		}(message)
	}
	return nil
}

func (beaver Beaver) transferFunds(walls []byte) error {

	var wallets models.Transfer

	if err := json.Unmarshal(walls, &wallets); err != nil {
		return err
	}

	if err := beaver.store.TransferFunds(context.Background(), wallets); err != nil {
		return err
	}

	return nil
}
