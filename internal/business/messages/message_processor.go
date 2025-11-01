package messages

import (
	"assignment/config"
	messageProcessor "assignment/internal/pkg/message_processor"
	"assignment/internal/repository/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"time"
)

type Processor interface {
	Send(payload messageProcessor.MessagePayload) (bool, error)
}

type AutoMessageProcessor struct {
	repo      repo
	cache     *redis.Client
	cfg       *config.Config
	processor Processor
	Running   bool
}

func NewAuthMessageProcessor(
	cfg *config.Config,
	messageRepo repo,
	processor Processor,
	cache *redis.Client,
) *AutoMessageProcessor {
	return &AutoMessageProcessor{
		repo:      messageRepo,
		cfg:       cfg,
		processor: processor,
		cache:     cache,
	}
}

func (p *AutoMessageProcessor) Process(ctx context.Context) error {
	log.Info("start processing in background")
	ticker := time.NewTicker(p.cfg.Message.SendInterval)

	p.Running = true

	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				p.Running = false

				log.Info("Processor stopped")
				return
			case t := <-ticker.C:
				fmt.Println("Processing next batch at :", t)
				p.ProcessInBatch(ctx, p.cfg.Message.BatchProcessCount)
			}
		}
	}()

	return nil
}

// what happens if this pod is at scale and we read from same database ?
// lets imagine someone entered 10s and 1000 per 10s what happens,
// I would definetly use a worker pool but for now I think it is enough to have this simple code...
func (p *AutoMessageProcessor) ProcessInBatch(ctx context.Context, batchCount int) {
	msgs, err := p.repo.GetMessagesByStatuses(
		batchCount,
		[]models.MessageStatus{models.StatusPending, models.StatusFailed},
	)
	if err != nil {
		log.Errorf("failed to fetch pending/failed messages")
	}

	log.Infof("Batch count: %d", len(msgs))

	for _, message := range msgs {
		log.Infof("processing recipient: %s", message.UUID)

		messageId := message.UUID.String()
		exists, _ := p.cache.Exists(ctx, fmt.Sprintf("message:%s", messageId)).Result()
		if exists > 0 {
			log.Infof("message %s already processed, skipping", message.UUID)
			continue
		}

		sent, err := p.processor.Send(messageProcessor.MessagePayload{
			To:      message.RecipientPhoneNumber,
			Content: message.Content,
		})
		if err != nil {
			log.Errorf("failed to send message to recipient:%s, saving for retry, err: %+v", message.UUID, err)
			if p.cfg.Message.Retry {
				p.repo.MessageToRetry(messageId, p.cfg.Message.RetryFailCount)
			}

			continue
		}

		if sent {
			p.repo.UpdateStatus(messageId, models.StatusSent)
			p.CacheMessageInfo(ctx, messageId)
		}

		log.Info("sent!")
	}
}

func (p *AutoMessageProcessor) CacheMessageInfo(ctx context.Context, messageID string) {
	type MessageCache struct {
		MessageID string    `json:"message_id"`
		SentAt    time.Time `json:"sent_at"`
	}

	cacheData := MessageCache{
		MessageID: messageID,
		SentAt:    time.Now().UTC(),
	}

	data, _ := json.Marshal(cacheData)

	ttl := 2 * time.Hour
	key := fmt.Sprintf("message:%s", messageID)

	log.Info("caching id: ", messageID)

	err := p.cache.Set(ctx, key, data, ttl).Err()
	if err != nil {
		log.Errorf("failed to cache message %s: %v", messageID, err)
	}
}
