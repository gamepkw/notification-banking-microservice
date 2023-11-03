package service

import (
	"context"
	"time"

	"gopkg.in/Shopify/sarama.v1"
)

type notificationService struct {
	contextTimeout time.Duration
	kafkaClient    sarama.Client
}

func NewNotificationService(timeout time.Duration, kafka sarama.Client) NotificationService {
	return &notificationService{
		contextTimeout: timeout,
		kafkaClient:    kafka,
	}
}

type NotificationService interface {
	SendTransactionNotification(ctx context.Context)
}
