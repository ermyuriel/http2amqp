package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/eucj/amqphelper"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	exchanges := strings.Split(os.Getenv("AMQP_EXCHANGES"), ",")
	paths := strings.Split(os.Getenv("RECEIVE_PATHS"), ",")
	routingKeys := strings.Split(os.Getenv("AMQP_ROUTING_KEYS"), ",")

	for i := range paths {
		conf := amqphelper.Configuration{Host: os.Getenv("AMQP_HOST"), ContentType: "application/json", RoutingKey: routingKeys[i], Exchange: exchanges[i]}
		q, err := amqphelper.GetQueue(&conf)
		if err != nil {
			log.Panicln(err)
		}
		sq := serverQueue{Queue: q, Path: paths[i]}
		http.HandleFunc(sq.Path, sq.serve)
	}

	http.HandleFunc(os.Getenv("HEALTHCHECK_PATH"), healthcheck)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("RECEIVE_PORT"), nil))
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("{\"Health\":\"OK\"}"))
}
