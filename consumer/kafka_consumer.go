package consumer

import (
	"encoding/json"
	"fmt"
	"hex/phoenix/config"
	email_factory "hex/phoenix/factories"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func StartConsumer() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.KafkaBroker,
		"group.id":          config.ConsumerGroupId,
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		fmt.Printf("Failed to connect to Kafka: %v\n", err)
	}

	err = c.Subscribe(config.Topic, nil)

	if err != nil {
		fmt.Printf("Failed to subscribe to topic %v: ,%v", config.Topic, err)
	}

	for {
		ev := c.Poll(100)

		switch e := ev.(type) {
		case *kafka.Message:
			var newEmail email_factory.Email

			err := json.Unmarshal(e.Value, &newEmail)

			if err != nil {
				fmt.Printf("Error unmarshaling json: %v", err)
				return
			}

			email_factory.NewEmail().Send(newEmail)
		case *kafka.Error:
			fmt.Printf("%+v\n", e)
		}
	}
}
