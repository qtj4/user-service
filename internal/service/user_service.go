package service

import (
	"context"

	"github.com/yourusername/user-service/internal/entity"
	"github.com/yourusername/user-service/internal/event"
	"github.com/yourusername/user-service/internal/repository"
)

type UserService interface {
	GetByID(ctx context.Context, id int) (*entity.User, error)
	GetUserOrders(ctx context.Context, userID int) ([]interface{}, error)
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id int) error
}

type userService struct {
	repo           repository.UserRepository
	eventPublisher *event.EventPublisher
}

func NewUserService(repo repository.UserRepository, eventPublisher *event.EventPublisher) UserService {
	return &userService{repo: repo, eventPublisher: eventPublisher}
}

func (s *userService) GetByID(ctx context.Context, id int) (*entity.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *userService) GetUserOrders(ctx context.Context, userID int) ([]interface{}, error) {
	// Placeholder: In a real scenario, this would fetch orders from an Order Service via gRPC or RabbitMQ
	return []interface{}{map[string]interface{}{"order_id": 1, "book_id": 101, "status": "completed"}}, nil
}

func (s *userService) Create(ctx context.Context, user *entity.User) error {
	err := s.repo.Create(ctx, user)
	if err != nil {
		return err
	}
	go s.eventPublisher.Publish(event.UserCreated, user)
	return nil
}

func (s *userService) Update(ctx context.Context, user *entity.User) error {
	err := s.repo.Update(ctx, user)
	if err != nil {
		return err
	}
	go s.eventPublisher.Publish(event.UserUpdated, user)
	return nil
}

func (s *userService) Delete(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	go s.eventPublisher.Publish(event.UserDeleted, map[string]int{"id": id})
	return nil
}