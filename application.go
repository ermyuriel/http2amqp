package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/eucj/amqphelper"
)

//HTTPRequest struct will be initialized with data from request, marshalled to JSON and pushed to the configured queues
type HTTPRequest struct {
	Method    string              `json:"method"`
	Header    map[string][]string `json:"header"`
	Body      interface{}         `json:"body"`
	Source    string              `json:"source"`
	Target    string              `json:"target"`
	Timestamp int64               `json:"timestamp"`
}

func main() {

	exchanges := strings.Split(os.Getenv("AMQP_EXCHANGES"), ",")
	paths := strings.Split(os.Getenv("RECEIVE_PATHS"), ",")
	routingKeys := strings.Split(os.Getenv("AMQP_ROUTING_KEYS"), ",")

	for i := range paths {
		conf := amqphelper.Configuration{Host: os.Getenv("AMQP_HOST"), ContentType: "application/json", RoutingKey: routingKeys[i], Exchange: exchanges[i], Durable: true}
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
