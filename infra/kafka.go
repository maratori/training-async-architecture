package infra

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Shopify/sarama"
)

const (
	KafkaAddr    = "$SVC_KAFKA_HOST:$SVC_KAFKA_PORT"
	KafkaGroupID = "$SVC_KAFKA_GROUP_ID"
)

func newConfig() (*sarama.Config, []string) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_1_0
	config.Producer.Return.Successes = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config, []string{os.ExpandEnv(KafkaAddr)}
}

type KafkaSender struct {
	producer sarama.SyncProducer
}

func NewKafkaSender() (*KafkaSender, error) {
	config, addrs := newConfig()

	producer, err := sarama.NewSyncProducer(addrs, config)
	if err != nil {
		return nil, fmt.Errorf("sarama.NewSyncProducer: %w", err)
	}

	sender := &KafkaSender{
		producer: producer,
	}

	return sender, nil
}

func (s *KafkaSender) Send(topic string, msg []byte) error {
	_, _, err := s.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(msg),
	})
	if err != nil {
		return fmt.Errorf("producer.SendMessage: %w", err)
	}
	return nil
}

func (s *KafkaSender) Close() {
	err := s.producer.Close()
	if err != nil {
		log.Printf("Can't close kafka producer: %+v\n", err)
	}
}

type KafkaHandler interface {
	HandleMessage(context.Context, []byte) error
}

type KafkaHandlerFunc func(context.Context, []byte) error

func (f KafkaHandlerFunc) HandleMessage(ctx context.Context, msg []byte) error {
	return f(ctx, msg)
}

type KafkaProcessor struct {
	consumerGroup sarama.ConsumerGroup
	handlers      map[string]KafkaHandler
}

func NewKafkaProcessor() (*KafkaProcessor, error) {
	config, addrs := newConfig()
	groupID := os.ExpandEnv(KafkaGroupID)

	consumerGroup, err := sarama.NewConsumerGroup(addrs, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("sarama.NewConsumerGroup: %w", err)
	}

	processor := &KafkaProcessor{
		consumerGroup: consumerGroup,
		handlers:      make(map[string]KafkaHandler),
	}

	return processor, nil
}

func (c *KafkaProcessor) RegisterHandler(topic string, handler KafkaHandler) {
	if _, ok := c.handlers[topic]; ok {
		panic(fmt.Sprintf("handler already registered for topic %q", topic))
	}

	c.handlers[topic] = handler
}

func (c *KafkaProcessor) RegisterHandlerFunc(topic string, handler func(context.Context, []byte) error) {
	c.RegisterHandler(topic, KafkaHandlerFunc(handler))
}

func (c *KafkaProcessor) Start(ctx context.Context) {
	for topic, handler := range c.handlers {
		go c.handleTopic(ctx, topic, handler)
	}
}

func (c *KafkaProcessor) Close() {
	err := c.consumerGroup.Close()
	if err != nil {
		log.Printf("Can't close kafka consumer group: %+v\n", err)
	}
}

func (c *KafkaProcessor) handleTopic(ctx context.Context, topic string, handler KafkaHandler) {
	for {
		err := c.consumerGroup.Consume(ctx, []string{topic}, NewSimpleConsumerGroupHandler(handler))
		if err != nil {
			log.Printf("Rejoin to consumer group: %+v\n", err)
		}
	}
}

type SimpleConsumerGroupHandler struct {
	handler     KafkaHandler
	stopOnError bool
}

func NewSimpleConsumerGroupHandler(processor KafkaHandler) *SimpleConsumerGroupHandler {
	return &SimpleConsumerGroupHandler{
		handler:     processor,
		stopOnError: false,
	}
}

func (c SimpleConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c SimpleConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c SimpleConsumerGroupHandler) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) (err error) {
	defer func() {
		if v := recover(); v != nil {
			err = fmt.Errorf("panic in handler %+v", v)
		}
	}()

	ctx := session.Context()
	messages := claim.Messages()

	for {
		select {
		case msg := <-messages:
			errP := c.handler.HandleMessage(ctx, msg.Value)
			if errP != nil {
				if c.stopOnError {
					return fmt.Errorf("handler.HandleMessage: %w", errP)
				}

				log.Printf("Can't handle message: %+v\n", errP)
			}
			session.MarkMessage(msg, "")

		// https://github.com/Shopify/sarama/issues/1192
		case <-ctx.Done():
			return nil
		}
	}
}
