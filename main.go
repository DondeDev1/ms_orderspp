package main

import (
	"ProductionOrders/Publishing"
	"ProductionOrders/elastic"
	"ProductionOrders/sap"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

func main() {

	// load .env file
	err := godotenv.Load("Test/.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	//MemCache
	//Crea el servidor cache
	//mc := memcache.New("localhost:11211")
	//delta := os.Getenv("DELTA_INIT")
	/*errm := mc.Add(&memcache.Item{Key: "Delta", Value: []byte(delta)})
	if errm != nil {
		log.Fatal(errm)
	}*/

	repo := sap.NewOrderRepository()
	orderService := Publishing.NewService(repo)
	elk_host := os.Getenv("ELK_HOST_HTTP")
	elk_port := os.Getenv("ELK_PORT")
	elk_index := os.Getenv("ELK_INDEX")
	elk_alias := os.Getenv("ELK_ALIAS")
	es, err := elastic.NewElasticRepository(fmt.Sprintf("http://%s:%s", elk_host, elk_port), elk_index, elk_alias, time.Second*30)
	if err != nil {
		fmt.Print("Error al conectar con elastic cache ")
		log.Fatal(err)
	}

	elasticService := Publishing.NewElasticRepository(es)

	err = Publishing.StarCronjob(context.Background(), orderService, elasticService)
	if err != nil {
		fmt.Println("Error ", err)
	}

}
