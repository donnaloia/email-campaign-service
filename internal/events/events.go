package events

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama" // Updated import path
)

type EventPublisher struct {
	producer sarama.SyncProducer
}

func NewEventPublisher(brokers []string) (*EventPublisher, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	return &EventPublisher{producer: producer}, nil
}

// CampaignLaunchedEvent represents the data structure for a campaign launch
type CampaignLaunchedEvent struct {
	CampaignID     string   `json:"campaign_id"`
	OrganizationID string   `json:"organization_id"`
	EmailAddresses []string `json:"email_addresses"`
	TemplateIDs    []string `json:"template_ids"`
}

func (p *EventPublisher) PublishCampaignLaunched(event CampaignLaunchedEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("error marshaling event: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: "campaign.launched",
		Value: sarama.StringEncoder(payload),
		Key:   sarama.StringEncoder(event.CampaignID),
	}

	_, _, err = p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("error publishing event: %w", err)
	}

	return nil
}

func (p *EventPublisher) Close() error {
	return p.producer.Close()
}
