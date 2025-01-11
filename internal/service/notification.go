package service

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/iBoBoTi/aqua-sec-inventory/internal/domain"
	"github.com/iBoBoTi/aqua-sec-inventory/internal/repository"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Notifier interface {
    Publish(message domain.Notification) error
    Listen() error
    Close()
}

type RabbitMQNotifier struct {
    conn    *amqp.Connection
    channel *amqp.Channel
    queue   amqp.Queue
    notificationRepo repository.NotificationRepository
}

func NewRabbitMQNotifier(amqpURL string, notificationRepo repository.NotificationRepository) (*RabbitMQNotifier, error) {
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
        notificationRepo: notificationRepo,
    }, nil
}

func (n *RabbitMQNotifier) Publish(payload domain.Notification) error {
    body, err := json.Marshal(&payload)
	if err != nil {
		log.Println(fmt.Errorf("error marshalling payload: %v", err))
		return err
	}

    return n.channel.Publish(
        "", 
        n.queue.Name,
        false, 
        false, 
        amqp.Publishing{
            ContentType: "text/plain",
            Body:        body,
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
            log.Printf("[NotificationService] Received: %s", d.Body)
            var payload domain.Notification
            err := json.Unmarshal(d.Body, &payload)
			
            if err != nil {
                log.Printf("Failed to decode message: %s", err)
                continue
            }
            if payload.Event == "notification" && payload.UserID != 0 && payload.Message != "" {
                log.Println("notification payload: ", payload)
                if err := n.notificationRepo.Create(&payload); err != nil {
                    log.Println("error creating notification: ", err)
                }
            }
            
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
