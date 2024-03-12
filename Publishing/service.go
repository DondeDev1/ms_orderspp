package Publishing

import (
	"ProductionOrders/order"
	"context"
)

// SAP Search
type OrdersRepository interface {
	GetOrders(ctx context.Context, fecha string, hora string) ([]order.Orden, error)
}
type service struct {
	repo OrdersRepository
}

func NewService(repo OrdersRepository) *service {
	return &service{repo: repo}
}

func (s *service) GetOrders(ctx context.Context, fecha string, hora string) ([]order.Orden, error) {
	return s.repo.GetOrders(ctx, fecha, hora)
}

// ElasticSearch
type OrderStorer interface {
	Insert(ctx context.Context, order order.Orden) error
	Update(ctx context.Context, order order.Orden) error
	FindOne(ctx context.Context, id string) (order.Orden, error)
}

type serviceElastic struct {
	elasticRepository OrderStorer
}

func NewElasticRepository(r OrderStorer) *serviceElastic {
	return &serviceElastic{elasticRepository: r}
}

func (s *serviceElastic) Insert(ctx context.Context, orden order.Orden) error {

	if err := s.elasticRepository.Insert(ctx, orden); err != nil {
		return err
	}
	return nil
}

func (s *serviceElastic) Update(ctx context.Context, orden order.Orden) error {
	if err := s.elasticRepository.Update(ctx, orden); err != nil {
		return err
	}
	return nil
}

func (s *serviceElastic) FindOne(ctx context.Context, id string) (order.Orden, error) {
	ord, err := s.elasticRepository.FindOne(ctx, id)
	if err != nil {
		return order.Orden{}, err
	}

	return ord, nil
}
