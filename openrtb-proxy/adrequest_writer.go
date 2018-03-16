package main

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"time"
	"log"
)

type AdRequestWriter struct {
	KafkaProducer *kafka.Producer
	KafkaTopic    string
}

//easyjson:json
type AdRequestResponseMessage struct {
	PageviewId string      `json:"pageview_id"`
	Ibbid      string      `json:"ibbid"`
	Ip         string      `json:"ip"`
	Url        string      `json:"url"`
	UserAgent  string      `json:"user_agent"`
	Positions  [] Position `json:"positions"`
	Prebid     bool        `json:"prebid"`
	Timestamp  int64       `json:"timestamp"`
	Status     int         `json:"status"`
}

//easyjson:json
type Position struct {
	ImpressionId string `json:"impression_id"`
	PositionId   string `json:"position_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
}

func NewAdRequestWriter(kafkaHosts string, topic string) *AdRequestWriter {
	kafkaProducer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":   kafkaHosts,
		"go.delivery.reports": false,
		"go.batch.producer":   true,
	})

	if err != nil {
		log.Fatalf("Failed to create producer: %s\n", err)
		panic(err)
	}

	return &AdRequestWriter{
		kafkaProducer,
		topic,
	}
}

func (writer *AdRequestWriter) write(request *AdRequest, adResponseStatus int) {
	message, err := NewAdRequestResponseMessage(request, adResponseStatus).MarshalJSON()

	if err != nil {
		log.Printf("Failed to serialize AdRequestResponseMessage: %s\n", err)
		return
	}

	deliveryChan := make(chan kafka.Event)
	err = writer.KafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &writer.KafkaTopic, Partition: kafka.PartitionAny},
		Value: message,
	}, deliveryChan)

	if err != nil {
		log.Printf("Failed to produce AdRequestResponseMessage: %s\n", err)
		return
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		log.Printf("Failed to deliver AdRequestResponseMessage: %s\n", m.TopicPartition.Error)
	}

	close(deliveryChan)
}

func (writer *AdRequestWriter) Close() {
	writer.KafkaProducer.Close()
}

func NewAdRequestResponseMessage(request *AdRequest, adResponseStatus int) (*AdRequestResponseMessage) {
	positions := make([]Position, len(request.Positions))
	for i, p := range request.Positions {
		positions[i] = Position{
			p.ImpressionId,
			p.PositionId,
			p.Width,
			p.Height,
		}
	}
	return &AdRequestResponseMessage{
		request.PageviewId,
		request.Ibbid,
		request.Ip,
		request.Url,
		request.UserAgent,
		positions,
		false,
		time.Now().UTC().Unix(),
		adResponseStatus,
	}
}
