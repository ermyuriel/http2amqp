package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/eucj/amqphelper"
)

func main() {
	http.HandleFunc(os.Getenv("HEALTHCHECK_PATH"), healthcheck)
	exchanges := strings.Split(os.Getenv("AMQP_EXCHANGES"), ",")
	paths := strings.Split(os.Getenv("RECEIVE_PATHS"), ",")
	routingKeys := strings.Split(os.Getenv("AMQP_ROUTING_KEYS"), ",")

	for i := range paths {
		q, err := amqphelper.GetQueue(os.Getenv("AMQP_HOST"), routingKeys[i], true, false, false, false, nil)
		if err != nil {
			log.Panicln(err)
		}
		sq := serverQueue{Queue: q, Path: paths[i], Exchange: exchanges[i], RoutingKey: routingKeys[i]}
		http.HandleFunc(sq.Path, sq.serve)
	}

	log.Fatal(http.ListenAndServe(":"+os.Getenv("RECEIVE_PORT"), nil))
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("{\"Health\":\"OK\"}"))
}
