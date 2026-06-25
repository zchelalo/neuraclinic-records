package rabbitmq

import (
	"context"
	"errors"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zchelalo/neuraclinic-records/internal/modules/attachments/ports"
	"go.uber.org/zap"
)

type Config struct {
	URL        string
	Exchange   string
	Queue      string
	RoutingKey string
	DLX        string
	DLQ        string
	Prefetch   int
}

type Handler func(ctx context.Context, routingKey string, body []byte) error

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	cfg     Config
	handler Handler
	logger  *zap.Logger
}

func NewConsumer(cfg Config, handler Handler, logger *zap.Logger) (*Consumer, error) {
	if logger == nil {
		logger = zap.NewNop()
	}

	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("dial rabbitmq: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("open rabbitmq channel: %w", err)
	}

	consumer := &Consumer{
		conn:    conn,
		channel: ch,
		cfg:     cfg,
		handler: handler,
		logger:  logger,
	}
	if err := consumer.declareTopology(); err != nil {
		_ = consumer.Close()
		return nil, err
	}

	return consumer, nil
}

func (c *Consumer) Run(ctx context.Context) error {
	deliveries, err := c.channel.Consume(
		c.cfg.Queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("start consuming: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case delivery, ok := <-deliveries:
			if !ok {
				return fmt.Errorf("rabbitmq delivery channel closed")
			}
			c.handleDelivery(ctx, delivery)
		}
	}
}

func (c *Consumer) Close() error {
	var err error
	if c.channel != nil {
		err = c.channel.Close()
	}
	if c.conn != nil {
		if closeErr := c.conn.Close(); err == nil {
			err = closeErr
		}
	}
	return err
}

func (c *Consumer) declareTopology() error {
	if err := c.channel.ExchangeDeclare(c.cfg.Exchange, "topic", true, false, false, false, nil); err != nil {
		return fmt.Errorf("declare exchange: %w", err)
	}
	if err := c.channel.ExchangeDeclare(c.cfg.DLX, "topic", true, false, false, false, nil); err != nil {
		return fmt.Errorf("declare dlx: %w", err)
	}

	if _, err := c.channel.QueueDeclare(c.cfg.DLQ, true, false, false, false, nil); err != nil {
		return fmt.Errorf("declare dlq: %w", err)
	}
	if err := c.channel.QueueBind(c.cfg.DLQ, c.deadLetterRoutingKey(), c.cfg.DLX, false, nil); err != nil {
		return fmt.Errorf("bind dlq: %w", err)
	}

	args := amqp.Table{
		"x-dead-letter-exchange":    c.cfg.DLX,
		"x-dead-letter-routing-key": c.deadLetterRoutingKey(),
	}
	if _, err := c.channel.QueueDeclare(c.cfg.Queue, true, false, false, false, args); err != nil {
		return fmt.Errorf("declare queue: %w", err)
	}
	if err := c.channel.QueueBind(c.cfg.Queue, c.cfg.RoutingKey, c.cfg.Exchange, false, nil); err != nil {
		return fmt.Errorf("bind queue: %w", err)
	}
	if err := c.channel.Qos(c.cfg.Prefetch, 0, false); err != nil {
		return fmt.Errorf("set qos: %w", err)
	}

	return nil
}

func (c *Consumer) handleDelivery(ctx context.Context, delivery amqp.Delivery) {
	err := c.handler(ctx, delivery.RoutingKey, delivery.Body)
	if err == nil {
		if ackErr := delivery.Ack(false); ackErr != nil {
			c.logger.Warn("rabbitmq ack failed", zap.Error(ackErr))
		}
		return
	}

	if errors.Is(err, ports.ErrInvalidEvent) {
		c.logger.Warn("invalid event rejected", zap.Error(err), zap.String("routing_key", delivery.RoutingKey))
		if rejectErr := delivery.Reject(false); rejectErr != nil {
			c.logger.Warn("rabbitmq reject failed", zap.Error(rejectErr))
		}
		return
	}

	c.logger.Warn("event processing failed; requeueing delivery", zap.Error(err), zap.String("routing_key", delivery.RoutingKey))
	if nackErr := delivery.Nack(false, true); nackErr != nil {
		c.logger.Warn("rabbitmq nack failed", zap.Error(nackErr))
	}
}

func (c *Consumer) deadLetterRoutingKey() string {
	return c.cfg.Queue + ".dead"
}
