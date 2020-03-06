package broker

import (
	"context"
	"encoding/json"
	"log"

	"github.com/MixinNetwork/supergroup.mixin.one/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-plugins/broker/snssqs"
)

func Connect() broker.Broker {
	cfg := config.AppConfig
	c := aws.NewConfig().
		WithRegion(cfg.AWS.Region).
		WithCredentials(credentials.NewStaticCredentials(cfg.AWS.Key, cfg.AWS.Secret, ""))

	b := snssqs.NewBroker(snssqs.SNSConfig(c), snssqs.SQSConfig(c))
	if err := b.Connect(); err != nil {
		log.Panicln(err)
	}

	return b
}

func Pub(b broker.Broker, e *Event) error {
	arnTopic := config.AppConfig.System.EventPubTopic
	if arnTopic == "" {
		return nil
	}

	data, err := json.Marshal(e)
	if err != nil {
		return err
	}

	msg := broker.Message{
		Body: data,
	}
	return b.Publish(arnTopic, &msg)
}

type Handler func(event *Event) error

func Sub(ctx context.Context, b broker.Broker, handler Handler) error {
	arnTopic := config.AppConfig.System.EventSubTopic

	h := func(event broker.Event) error {
		var e Event
		if err := json.Unmarshal(event.Message().Body, &e); err != nil {
			return event.Ack()
		}

		if err := handler(&e); err != nil {
			return err
		}
		return event.Ack()
	}

	sub, err := b.Subscribe(arnTopic, h, snssqs.VisibilityTimeout(10))
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		_ = sub.Unsubscribe()
		return ctx.Err()
	}
}
