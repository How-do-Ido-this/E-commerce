package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
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

// Publish envía un mensaje a un exchange.
func (c *Client) Publish(exchange, routingKey string, body []byte) error {
	return c.channel.Publish(
		exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

// Consume se suscribe a una cola y devuelve un canal de entregas.
func (c *Client) Consume(queueName string) (<-chan amqp.Delivery, error) {
	q, err := c.channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	return c.channel.Consume(
		q.Name,
		"",    // consumer tag
		true,  // auto-ack
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
