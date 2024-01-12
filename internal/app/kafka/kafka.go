// Package kafka implements general TCP kafka
package kafka

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/kirill-a-belov/test_task_framework/internal/app/kafka/pkg/config"
	"github.com/kirill-a-belov/test_task_framework/internal/pkg/protocol"
	"github.com/kirill-a-belov/test_task_framework/pkg/context_helper"
	"github.com/kirill-a-belov/test_task_framework/pkg/logger"
	"github.com/kirill-a-belov/test_task_framework/pkg/tracer"
)

func New(ctx context.Context, config *config.Config) *Kafka {
	_, span := tracer.Start(ctx, "internal.app.kafka.New")
	defer span.End()

	return &Kafka{
		config:   config,
		stopChan: make(chan struct{}),
		logger:   logger.New("kafka"),
	}
}

type Kafka struct {
	config   *config.Config
	stopChan chan struct{}
	logger   logger.Logger

	producer         sarama.SyncProducer
	consumer         sarama.Consumer
	partConsumerList []sarama.PartitionConsumer
}

const (
	topicA = "example_topic_a"
	topicB = "example_topic_B"
)

func (k *Kafka) Start(ctx context.Context) error {
	_, span := tracer.Start(ctx, "internal.app.kafka.Kafka.Start")
	defer span.End()

	k.logger.Info(fmt.Sprintf("config: %+v", *k.config))

	kafkaServiceAddressList := []string{k.config.Address}
	p, err := sarama.NewSyncProducer(kafkaServiceAddressList, nil)
	if err != nil {
		return errors.Wrap(err, "creating producer")
	}
	k.producer = p
	go k.processor(ctx, sender(k.producer))

	c, err := sarama.NewConsumer(kafkaServiceAddressList, nil)
	if err != nil {
		return errors.Wrap(err, "creating consumer")
	}
	k.consumer = c

	topicList := []string{topicA, topicB}
	k.partConsumerList = make([]sarama.PartitionConsumer, len(topicList))
	for i, topic := range topicList {
		if k.partConsumerList[i], err = k.consumer.ConsumePartition(topic, 0, sarama.OffsetNewest); err != nil {
			return errors.Wrapf(err, "creating consumer for topic (%s)", topic)
		}
		go k.processor(ctx, receiver(k.partConsumerList[i]))
	}

	return nil
}

func (k *Kafka) processor(ctx context.Context, interactor func() error) {
	_, span := tracer.Start(ctx, "internal.app.kafka.Kafka.processor")
	defer span.End()

	for {
		select {
		case <-k.stopChan:
			k.logger.Info("processor terminated")

			return
		default:
			if err := context_helper.RunWithTimeout(k.config.ConnTTL, func() error {
				return interactor()
			}); err != nil {
				k.logger.Error(err, "kafka interaction")
			}

			time.Sleep(k.config.Delay)
		}
	}
}

func (k *Kafka) Stop(ctx context.Context) {
	_, span := tracer.Start(ctx, "internal.app.kafka.Kafka.Stop")
	defer span.End()

	for _, pc := range k.partConsumerList {
		_ = pc.Close()
	}
	_ = k.consumer.Close()
	_ = k.producer.Close()

	close(k.stopChan)
}

func sender(producer sarama.SyncProducer) func() error {
	return func() error {
		uuidA := uuid.NewString()
		msgA := &bytes.Buffer{}
		if err := gob.NewEncoder(msgA).
			Encode(protocol.Message{
				Type: protocol.MessageTypeRequest,
			},
			); err != nil {
			return errors.Wrap(err, "encoding msg A")
		}
		uuidB := uuid.NewString()
		msgB := &bytes.Buffer{}
		if err := gob.NewEncoder(msgB).
			Encode(protocol.Message{
				Type: protocol.MessageTypeResponse,
			},
			); err != nil {
			return errors.Wrap(err, "encoding msg B")
		}

		msgList := []*sarama.ProducerMessage{
			{
				Topic: topicA,
				Key:   sarama.StringEncoder(uuidA),
				Value: sarama.ByteEncoder(msgA.Bytes()),
			},
			{
				Topic: topicB,
				Key:   sarama.StringEncoder(uuidB),
				Value: sarama.ByteEncoder(msgB.Bytes()),
			},
		}

		if err := producer.SendMessages(msgList); err != nil {
			return errors.Wrap(err, "sending messages")
		}

		logger.New("kafka.sender").Info("sending messages",
			"sent keys", uuidA, uuidB,
		)

		return nil
	}
}

func receiver(consumer sarama.PartitionConsumer) func() error {
	return func() error {
		msg, ok := <-consumer.Messages()
		if !ok {
			return errors.New("reading msg from consumer chan")
		}

		msgValue := &protocol.Request{}
		msgB := &bytes.Buffer{}
		if _, err := msgB.Write(msg.Value); err != nil {
			return errors.Wrap(err, "writing msg buffer")
		}
		if err := gob.NewDecoder(msgB).Decode(msgValue); err != nil {
			return errors.Wrap(err, "decoding msg value")
		}

		logger.New("kafka.receiver").Info("receiving message",
			"received key", string(msg.Key))

		return nil
	}
}
