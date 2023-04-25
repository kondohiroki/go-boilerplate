package queue

import (
	"github.com/kondohiroki/go-boilerplate/internal/helper/queue"
	"github.com/kondohiroki/go-boilerplate/internal/repository"
)

type QueueApp interface {
	GetQueues() ([]GetQueueDTO, error)
	// GetQueueByID(GetQueueDTI) (GetQueueDTO, error)
}

type queueApp struct {
	Repo repository.Repository
}

func NewQueueApp() QueueApp {
	return &queueApp{
		Repo: *repository.NewRepository(),
	}
}

type GetQueueDTI struct {
	Key string `json:"key"`
}

type GetQueueDTO struct {
	Key              string `json:"key,omitempty"`
	KeyWithoutPrefix string `json:"key_without_prefix,omitempty"`
	NumberOfItems    int64  `json:"number_of_items"`
}

func (app *queueApp) GetQueues() ([]GetQueueDTO, error) {
	var queues []GetQueueDTO

	qs, err := queue.ListQueueKeysAndLengths()
	if err != nil {
		return nil, err
	}
	for _, q := range qs {
		queue := GetQueueDTO{
			Key:              q.Key,
			KeyWithoutPrefix: q.KeyWithoutPrefix,
			NumberOfItems:    q.NumberOfItems,
		}
		queues = append(queues, queue)
	}

	return queues, nil
}
