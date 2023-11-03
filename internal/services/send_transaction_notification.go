package service

import (
	"context"
	"fmt"
	"strings"

	consumer "github.com/gamepkw/notification-banking-microservice/pkg/kafka/consumer"
	"gopkg.in/Shopify/sarama.v1"
)

func (n *notificationService) SendTransactionNotification(ctx context.Context) {
	topic := "sms_transaction"
	var partition int32 = 0
	var offset int64 = sarama.OffsetNewest
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Send Notification stopped")
			return
		default:

			// Consume a Kafka message
			message := consumer.RunKafkaConsumer(n.kafkaClient, topic, partition, offset)
			message_split := strings.Split(message, "|")
			if message_split[0] == "withdraw" {
				fmt.Printf("%s Baht has been withdrawn from account no: %s at %s\nRemaining balance: %s\n%s\n",
					message_split[1],
					message_split[2],
					message_split[3],
					message_split[4],
					strings.Repeat("-", 120))

			} else if message_split[0] == "deposit" {
				fmt.Printf("%s Baht has been deposited into account no: %s at %s\nRemaining balance: %s\n%s\n",
					message_split[1],
					message_split[2],
					message_split[3],
					message_split[4],
					strings.Repeat("-", 120))
			} else if message_split[0] == "transfer" {
				fmt.Printf("%s Baht has been transferred from account no: %s to account no: %s at %s\n%s\n",
					message_split[1],
					message_split[2],
					message_split[5],
					message_split[3],
					strings.Repeat("-", 120))
			}

			// offset += 1

			// time.Sleep(time.Second)
		}

	}

}
