package Publishing

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
)

type Router interface {
	Handle(method, path string, handler http.Handler)
}

func RegisterRoutes(router *httprouter.Router, s Service) {
	getOrderHandler := kithttp.NewServer(
		MakeEnpointGetOrders(s),
		decodeGetArticleRequest,
		encondeGetOrderResponse,
	)

	router.Handler(http.MethodGet, "/orders/:fecha/:hora", getOrderHandler)
}

func decodeGetArticleRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	params := httprouter.ParamsFromContext(ctx)
	return GetOrdersRequestModel{
		fecha: params.ByName("fecha"),
		hora:  params.ByName("hora"),
	}, nil
}

func encondeGetOrderResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res, ok := response.(GetOrdersResponseModel)

	if !ok {
		return fmt.Errorf("encodeGetOrdersResponse failed cast response")
	}
	formatted := formatGetOrdenResponse(res)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(formatted)
}
