package service

import (
    "fmt"
    "log"

    amqp "github.com/rabbitmq/amqp091-go"
)

type Notifier interface {
    Publish(message string) error
    Listen() error
    Close()
}

type RabbitMQNotifier struct {
    conn    *amqp.Connection
    channel *amqp.Channel
    queue   amqp.Queue
}

func NewRabbitMQNotifier(amqpURL string) (*RabbitMQNotifier, error) {
    conn, err := amqp.Dial(amqpURL)
    if err != nil {
        return nil, fmt.Errorf("error dailing RabbitMQ server: %w", err)
    }
    ch, err := conn.Channel()
    if err != nil {
        return nil, fmt.Errorf("error opening rabbitmq channel: %w", err)
    }
    q, err := ch.QueueDeclare(
        "notifications",
        true,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        return nil, fmt.Errorf("error declaring queue: %w", err)
    }

    return &RabbitMQNotifier{
        conn:    conn,
        channel: ch,
        queue:   q,
    }, nil
}

func (n *RabbitMQNotifier) Publish(message string) error {
    return n.channel.Publish(
        "", 
        n.queue.Name,
        false, 
        false, 
        amqp.Publishing{
            ContentType: "text/plain",
            Body:        []byte(message),
        },
    )
}

func (n *RabbitMQNotifier) Listen() error {
    msgs, err := n.channel.Consume(
        n.queue.Name,
        "",
        true,  // auto-ack
        false, // exclusive
        false, // noLocal
        false, // noWait
        nil,
    )
    if err != nil {
        return err
    }

    forever := make(chan bool)

    go func() {
        for d := range msgs {
            // For now, just log them
            log.Printf("[NotificationService] Received: %s", d.Body)
        }
    }()

    fmt.Println("[NotificationService] Listening for messages... Press CTRL+C to exit.")
    <-forever
    return nil
}

func (n *RabbitMQNotifier) Close() {
    if n.channel != nil {
        _ = n.channel.Close()
    }
    if n.conn != nil {
        _ = n.conn.Close()
    }
}
