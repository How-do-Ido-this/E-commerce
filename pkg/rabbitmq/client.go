package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	OrdersExchange          = "orders"
	OrdersExchangeType      = "direct"
	OrderCreatedKey         = "order.created"
	InventoryReservedKey    = "inventory.reserved"
	InventoryFailedKey      = "inventory.failed"
	InventoryQueue          = "inventory.queue"
	OrderQueue              = "order.queue"
)

// Client envuelve la conexión y el canal de RabbitMQ.
type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewClient conecta a RabbitMQ y abre un canal.
func NewClient(amqpURL string) (*Client, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open rabbitmq channel: %w", err)
	}

	return &Client{conn: conn, channel: ch}, nil
}

// SetupTopology declara el exchange, las colas y sus bindings.
func (c *Client) SetupTopology() error {
	if err := c.channel.ExchangeDeclare(
		OrdersExchange,
		OrdersExchangeType,
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // args
	); err != nil {
		return fmt.Errorf("failed to declare exchange %s: %w", OrdersExchange, err)
	}

	if _, err := c.channel.QueueDeclare(
		InventoryQueue,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // args
	); err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", InventoryQueue, err)
	}

	if err := c.channel.QueueBind(
		InventoryQueue,
		OrderCreatedKey,
		OrdersExchange,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("failed to bind queue %s to %s: %w", InventoryQueue, OrderCreatedKey, err)
	}

	if _, err := c.channel.QueueDeclare(
		OrderQueue,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // args
	); err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", OrderQueue, err)
	}

	for _, routingKey := range []string{InventoryReservedKey, InventoryFailedKey} {
		if err := c.channel.QueueBind(
			OrderQueue,
			routingKey,
			OrdersExchange,
			false,
			nil,
		); err != nil {
			return fmt.Errorf("failed to bind queue %s to %s: %w", OrderQueue, routingKey, err)
		}
	}

	return nil
}

// Publish envía un mensaje a un exchange con headers (para el CorrelationID).
func (c *Client) Publish(exchange, routingKey string, body []byte, headers map[string]interface{}) error {
	return c.channel.Publish(
		exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Headers:     amqp.Table(headers),
		},
	)
}

// Consume se suscribe a una cola con Ack manual.
func (c *Client) Consume(queueName, consumerTag string) (<-chan amqp.Delivery, error) {
	return c.channel.Consume(
		queueName,
		consumerTag,
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
}

// Close cierra la conexión y el canal.
func (c *Client) Close() {
	c.channel.Close()
	c.conn.Close()
}
