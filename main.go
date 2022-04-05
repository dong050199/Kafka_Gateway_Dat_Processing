package main

import (
	D "MQTTClient/driver"
	"MQTTClient/kafka"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	D.Connect()
	go func() {
		kafka.Kafka_Sub(D.ClientMQ)
	}()
	r := mux.NewRouter()
	log.Fatal(http.ListenAndServe(":8080", r))

}
