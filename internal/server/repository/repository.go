// the repository layer is designed for working
// with data and processing it
package repository

import (
	"context"
	"encoding/json"

	"github.com/BelyaevEI/wallet/internal/config"
	"github.com/BelyaevEI/wallet/internal/models"
	"github.com/BelyaevEI/wallet/internal/store"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Implementation check
var _ Repositorer = Repository{}

type Repositorer interface {
	CheckExists(ctx context.Context, id uint32) (bool, error)
	GetBalanceByID(ctx context.Context, id uint32) (int, error)
	CheckFundsByID(ctx context.Context, id uint32, amount float64) (bool, error)
	SendMessageToWorker(ctx context.Context, mes models.Transfer) error
	Shutdown()
}

// Repository layer
type Repository struct {
	Store     store.Storer
	channel   *amqp.Channel
	queueName string
}

// Create new repository for service
func NewRepo(cfg config.Config) (Repositorer, error) {

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

	// Connect to database
	store, err := store.Connect(cfg)
	if err != nil {
		return nil, err
	}

	return Repository{Store: store, channel: ch, queueName: q.Name}, nil
}

// Check exists wallet by id
func (repository Repository) CheckExists(ctx context.Context, id uint32) (bool, error) {
	return repository.Store.CheckExists(ctx, id)
}

// Getting  balance by id
func (repository Repository) GetBalanceByID(ctx context.Context, id uint32) (int, error) {
	return repository.Store.GetBalanceByID(ctx, id)
}

// Closing open connection
func (repository Repository) Shutdown() {
	repository.Store.CloseConnection2DB()
}

// Check funds
func (repository Repository) CheckFundsByID(ctx context.Context, id uint32, amount float64) (bool, error) {
	balance, err := repository.Store.GetBalanceByID(ctx, id)
	if err != nil {
		return false, err
	}

	return balance > int(amount), nil
}

func (repository Repository) SendMessageToWorker(ctx context.Context, mes models.Transfer) error {

	// Marshal message to slice byte
	byteMes, err := json.Marshal(mes)
	if err != nil {
		return err
	}

	// Publish new message
	err = repository.channel.PublishWithContext(ctx,
		"",                   // exchange
		repository.queueName, // routing key
		false,                // mandatory
		false,                // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         byteMes,
			DeliveryMode: amqp.Persistent,
		})

	return err
}
