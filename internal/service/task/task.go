package task

import (
	"context"
	"fmt"

	publisher "github.com/danieeelfr/swd-challenge/internal/async/publisher/notification"

	broker "github.com/danieeelfr/swd-challenge/internal/async/pubsub/broker"
	sub "github.com/danieeelfr/swd-challenge/internal/async/pubsub/subscriber"
	"github.com/danieeelfr/swd-challenge/internal/config"
	"github.com/danieeelfr/swd-challenge/internal/models"
	repo "github.com/danieeelfr/swd-challenge/internal/repository"

	log "github.com/sirupsen/logrus"
)

// Task holds the performed task implementation
type Task struct {
	broker     *broker.Broker
	sub        *sub.Subscriber
	repository *repo.MySQLRepo
}

// NewTaskService returns a new task implementation or error
func NewTaskService(cfg *config.Config) (*Task, error) {
	log.Info("starting MySQL connection...")

	r := repo.NewMySQLRepo(cfg.MySQLRepositoryConfig)
	err := r.Connect()
	if err != nil {
		return nil, err
	}

	b := broker.NewBroker()
	sub1 := b.AddSubscriber()
	b.Subscribe(sub1, "MANAGER")

	log.Info("listening async topics...")
	go sub1.Listen()

	return &Task{repository: r, broker: b, sub: sub1}, nil
}

// Find finds tasks according to the filter
func (t *Task) Find(ctx context.Context, filter *models.TaskFilter) ([]models.Task, error) {

	var tasks []models.Task

	t.repository.DB.Table("tasks")

	if filter.UserID != "" {
		if err := t.repository.DB.Where("user_id = ?", filter.UserID).Find(&tasks).Error; err != nil {
			return nil, err
		}
		return tasks, nil
	}

	if err := t.repository.DB.Find(&tasks).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}

// Save a new task and notify managers if it is the case
func (t *Task) Save(ctx context.Context, task *models.Task, user *models.User) error {
	t.repository.DB.Save(task)

	if user.Role == "TECHNICIAN" {
		msg := fmt.Sprintf(
			"The tech %s performed the task %s on date %s",
			user.Name, task.ID, task.PerformedAt.String(),
		)

		go publisher.Notify(t.broker, msg)
	}

	return nil
}
