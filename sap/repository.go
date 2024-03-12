package sap

import (
	"ProductionOrders/order"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Define structs to represent the JSON structure Mapped to SAP Result
type Metadata struct {
	ID   string `json:"id"`
	URI  string `json:"uri"`
	Type string `json:"type"`
}

type ProductionOrder struct {
	Metadata      Metadata `json:"__metadata"`
	Order         string   `json:"order"`
	ClaseOrden    string   `json:"claseOrden"`
	Amases        string   `json:"amases"`
	FechaInicio   string   `json:"fechaInicio"`
	HoraInicio    string   `json:"horaInicio"`
	FechaFin      string   `json:"fechaFin"`
	HoraFin       string   `json:"horaFin"`
	VarNombre     string   `json:"varNombre"`
	PresNombre    string   `json:"presNombre"`
	LineaNombre   string   `json:"lineaNombre"`
	Variedad      string   `json:"variedad"`
	Presentacion  string   `json:"presentacion"`
	Linea         string   `json:"linea"`
	Cantidad      string   `json:"cantidad"`
	Uma           string   `json:"uma"`
	CantUnidad    string   `json:"cantUnidad"`
	PuestoTrabajo string   `json:"puestoTrabajo"`
	Estatus       string   `json:"estatus"`
}

type Results struct {
	Metadata Metadata          `json:"__metadata"`
	Results  []ProductionOrder `json:"results"`
}

type ResponseBody struct {
	D Results `json:"d"`
}

type orderRepository struct {
	orderspp map[string]order.Orden
}

func NewOrderRepository() *orderRepository {
	return &orderRepository{
		orderspp: map[string]order.Orden{},
	}
}

func (r *orderRepository) GetOrders(ctx context.Context, fecha string, hora string) ([]order.Orden, error) {
	ordersList := []order.Orden{}
	fmt.Println("Consultando: ...")
	username := os.Getenv("SAP_USERNAME")
	password := os.Getenv("SAP_PASSWORD")
	sap_host := os.Getenv("SAP_HOST_HTTPS")
	sap_port := os.Getenv("SAP_PORT_HTTPS")
	sap_env := os.Getenv("SAP_ENV")

	//SAP API
	//Golang net/http not encode https automatically only http, with https u must encode manually
	url := "https://" + sap_host + ":" + sap_port + "/sap/opu/odata/sap/ZQM_ODATA_SRV/ProductionOrdersSet?$filter=fechaInicio%20eq%20%27" + fecha + "%27%20and%20horaInicio%20eq%20%27" + hora + "%27&sap-client=" + sap_env + "&$format=json"
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // <--- Problem
	}
	client := &http.Client{Transport: tr}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	// Add Basic Authentication header
	req.SetBasicAuth(username, password)

	// Make the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	/*if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}*/
	var response ResponseBody
	//Marshal content
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	// Access the parsed data
	for _, productionOrder := range response.D.Results {
		ordersList = append(ordersList, order.Orden{
			Order:         productionOrder.Order,
			ClaseOrden:    productionOrder.ClaseOrden,
			Amases:        productionOrder.Amases,
			FechaInicio:   productionOrder.FechaInicio,
			HoraInicio:    productionOrder.HoraInicio,
			FechaFin:      productionOrder.FechaFin,
			HoraFin:       productionOrder.HoraFin,
			VarNombre:     productionOrder.VarNombre,
			PresNombre:    productionOrder.PresNombre,
			LineaNombre:   productionOrder.LineaNombre,
			Variedad:      productionOrder.Variedad,
			Presentacion:  productionOrder.Presentacion,
			Linea:         productionOrder.Linea,
			Cantidad:      productionOrder.Cantidad,
			Uma:           productionOrder.Uma,
			CantUnidad:    productionOrder.CantUnidad,
			PuestoTrabajo: productionOrder.PuestoTrabajo,
			Estatus:       productionOrder.Estatus,
		})
	}
	fmt.Println("data recolectada ", len(ordersList))
	return ordersList, nil
}
