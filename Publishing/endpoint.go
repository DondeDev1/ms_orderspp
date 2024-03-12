package Publishing

import (
	"ProductionOrders/order"
	"context"
	"errors"
	"fmt"
	"github.com/go-kit/kit/endpoint"
)

type Service interface {
	GetOrders(ctx context.Context, fecha string, hora string) ([]order.Orden, error)
}

type ServiceElastic interface {
	Insert(ctx context.Context, order order.Orden) error
	Update(ctx context.Context, order order.Orden) error
	FindOne(ctx context.Context, id string) (order.Orden, error)
}

// Get Orders
type GetOrdersRequestModel struct {
	fecha string
	hora  string
}

type GetOrdersResponseModel struct {
	Orders []order.Orden
}

func MakeEnpointGetOrders(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(GetOrdersRequestModel)
		if !ok {
			return nil, errors.New("MakeEndpointGetOrders failed cast request")
		}

		a, err := s.GetOrders(ctx, req.fecha, req.hora)
		if err != nil {
			return nil, fmt.Errorf("MakeEndpointGetOrders: %w", err)
		}

		return GetOrdersResponseModel{
			Orders: a,
		}, nil
	}
}
