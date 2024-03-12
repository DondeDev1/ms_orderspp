package elastic

import (
	"ProductionOrders/order"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"time"
)

type elasticSearch struct {
	client  *elasticsearch.Client
	index   string
	alias   string
	timeout time.Duration
}

func NewElasticRepository(url string, index string, alias string, time time.Duration) (*elasticSearch, error) {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{url},
	})

	if err != nil {
		return nil, err
	}

	return &elasticSearch{client: client, alias: index, index: alias, timeout: time}, nil
}

func (e *elasticSearch) CreateIndex(index string) error {
	e.index = index
	e.alias = index + "_alias"

	res, err := e.client.Indices.Exists([]string{e.index})

	if res.StatusCode == 200 {
		return nil
	}
	if res.StatusCode != 404 {
		return fmt.Errorf("Error al verificar existencia de indice 404: %s", res.String())
	}

	res, err = e.client.Indices.Create(e.index)
	if err != nil {
		return fmt.Errorf("No es posible crear el indice: %w", err)
	}
	if res.IsError() {
		return fmt.Errorf("Error al intentar crear el indice indices: %s", res.String())
	}

	res, err = e.client.Indices.PutAlias([]string{e.index}, e.alias)
	if err != nil {
		return fmt.Errorf("No se puede crear el alias del indice: %w", err)
	}
	if res.IsError() {
		return fmt.Errorf("Error al intentar crear el alias del indice: %s", res.String())
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////

type document struct {
	Source interface{} `json:"_source"`
}

func (o elasticSearch) Insert(ctx context.Context, orden order.Orden) error {

	bdy, err := json.Marshal(orden)
	if err != nil {
		return fmt.Errorf("insert: marshall: %w", err)
	}

	// res, err := p.elastic.client.Create()
	req := esapi.CreateRequest{
		Index:      o.alias,
		DocumentID: orden.Order,
		Body:       bytes.NewReader(bdy),
	}

	ctx, cancel := context.WithTimeout(ctx, o.timeout)
	defer cancel()

	res, err := req.Do(ctx, o.client)
	if err != nil {
		return fmt.Errorf("insert: request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 409 {
		return fmt.Errorf("Error 409")
	}

	if res.IsError() {
		return fmt.Errorf("insert: response: %s", res.String())
	}

	return nil

}

func (o elasticSearch) Update(ctx context.Context, orden order.Orden) error {
	bdy, err := json.Marshal(orden)
	if err != nil {
		return fmt.Errorf("update: marshall: %w", err)
	}

	// res, err := p.elastic.client.Update()
	req := esapi.UpdateRequest{
		Index:      o.alias,
		DocumentID: orden.Order,
		Body:       bytes.NewReader([]byte(fmt.Sprintf(`{"doc":%s}`, bdy))),
	}

	ctx, cancel := context.WithTimeout(ctx, o.timeout)
	defer cancel()

	res, err := req.Do(ctx, o.client)
	if err != nil {
		return fmt.Errorf("update: request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return fmt.Errorf("Error not found 404 %s", res.String())
	}

	if res.IsError() {
		return fmt.Errorf("update: response: %s", res.String())
	}

	return nil
}

func (o elasticSearch) FindOne(ctx context.Context, id string) (order.Orden, error) {
	req := esapi.GetRequest{
		Index:      o.alias,
		DocumentID: id,
	}

	ctx, cancel := context.WithTimeout(ctx, o.timeout)
	defer cancel()

	res, err := req.Do(ctx, o.client)
	if err != nil {
		return order.Orden{}, fmt.Errorf("find one: request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return order.Orden{}, fmt.Errorf("Error not found 404 %s", res.String())
	}

	if res.IsError() {
		return order.Orden{}, fmt.Errorf("find one: response: %s", res.String())
	}

	var (
		orden order.Orden
		body  document
	)

	body.Source = orden

	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return order.Orden{}, fmt.Errorf("Encontrado: 1 : Error decode: %w", err)
	}
	//acceder a los datos del documento a través del campo Source, utilizando una aserción de tipo (type assertion)
	src, _ := body.Source.(map[string]interface{})

	ord := order.Orden{
		Order:         fmt.Sprintf("%v", src["Order"]),
		ClaseOrden:    fmt.Sprintf("%v", src["ClaseOrden"]),
		Amases:        fmt.Sprintf("%v", src["Amases"]),
		FechaInicio:   fmt.Sprintf("%v", src["FechaInicio"]),
		HoraInicio:    fmt.Sprintf("%v", src["HoraInicio"]),
		FechaFin:      fmt.Sprintf("%v", src["FechaFin"]),
		HoraFin:       fmt.Sprintf("%v", src["HoraFin"]),
		VarNombre:     fmt.Sprintf("%v", src["VarNombre"]),
		PresNombre:    fmt.Sprintf("%v", src["PresNombre"]),
		LineaNombre:   fmt.Sprintf("%v", src["LineaNombre"]),
		Variedad:      fmt.Sprintf("%v", src["Variedad"]),
		Presentacion:  fmt.Sprintf("%v", src["Presentacion"]),
		Linea:         fmt.Sprintf("%v", src["Linea"]),
		Cantidad:      fmt.Sprintf("%v", src["Cantidad"]),
		Uma:           fmt.Sprintf("%v", src["Uma"]),
		CantUnidad:    fmt.Sprintf("%v", src["CantUnidad"]),
		PuestoTrabajo: fmt.Sprintf("%v", src["PuestoTrabajo"]),
		Estatus:       fmt.Sprintf("%v", src["Estatus"]),
	}
	return ord, nil
}
