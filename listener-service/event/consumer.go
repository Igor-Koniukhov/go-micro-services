package event

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"net/http"
)

const webPort = "8082"

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}
	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	return declareExchange(channel)
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}
	for _, topic := range topics {
		err = ch.QueueBind(
			q.Name,
			topic,
			"logs_topic",
			false,
			nil,
		)
		if err != nil {
			return err
		}

	}

	messages, err := ch.Consume(
		q.Name, "",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	forever := make(chan bool)
	go func() {
		for message := range messages {
			var payload Payload
			_ = json.Unmarshal(message.Body, &payload)
			go handlePayload(payload)

		}
	}()
	fmt.Printf(" Waiting for messages [Exchange, Queue] [logs_topic, %s]\n", q.Name)
	<-forever
	return nil
}

func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		err := logEvent(payload)
		if err != nil {
			fmt.Println(err)
		}

	case "auth":
		fmt.Println("auth topic")
	default:
		err := logEvent(payload)
		if err != nil {
			fmt.Println(err)
		}

	}
}

func logEvent(entry Payload) error {
	jsonData, err := json.MarshalIndent(entry, "", "\t")
	logServiceURL := fmt.Sprintf("http://logger-service:%v/log", webPort)

	req, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		return errors.New("bad status code")
	}

	return nil
}
