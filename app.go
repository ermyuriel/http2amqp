package http2amqp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ermyuriel/amqphelper"
	lambdaproxy "github.com/ermyuriel/go-lambda-proxy"
)

type HTTPRequest struct {
	Method    string              `json:"method"`
	Header    map[string][]string `json:"header"`
	Body      interface{}         `json:"body"`
	Source    string              `json:"source"`
	Target    string              `json:"target"`
	Timestamp int64               `json:"timestamp"`
}

const (
	contentType = "application/json"
)

var host string
var queues map[string]map[string]*amqphelper.Queue

func init() {
	host = os.Getenv("AMQP_HOST")
	if host == "" {
		log.Panicln("No instance defined")
	}
	queues = make(map[string]map[string]*amqphelper.Queue)

}

func main() {
	lambda.Start(lambdaproxy.ProxyFunction(transformMessage))
}

func push(message []byte, exchange, routingKey string) error {
	var queue *amqphelper.Queue

	if m, e := queues[exchange]; !e || m == nil {
		queues[exchange] = make(map[string]*amqphelper.Queue)
	}

	if q, e := queues[exchange][routingKey]; e {
		queue = q
	} else {
		conf := amqphelper.Configuration{Host: host, ContentType: contentType, RoutingKey: routingKey, Exchange: exchange, Durable: true}

		q, err := amqphelper.GetQueue(&conf)
		if err != nil {
			return err
		}
		queues[exchange][routingKey] = q
		queue = q
	}

	return queue.Publish(message, nil, true, false)
}

func transformMessage(ctx context.Context) (interface{}, error) {

	results := make(chan error)

	go func(results chan error) {
		var body interface{}

		request, is := ctx.Value("request").(*http.Request)

		if !is {
			results <- errors.New("No valid request object in context")
			return
		}

		if request == nil {
			results <- errors.New("No valid request object in context")
			return
		}

		exchange := request.URL.Query().Get("exchange")

		routingKey := request.URL.Query().Get("key")

		err := json.NewDecoder(request.Body).Decode(&body)

		if err != nil {
			results <- err
			return
		}

		toSend := HTTPRequest{Method: request.Method, Header: request.Header, Body: body, Source: request.RemoteAddr, Target: request.RequestURI, Timestamp: time.Now().Unix()}
		message, err := json.Marshal(toSend)

		if err != nil {
			results <- err
			return
		}

		results <- push(message, routingKey, exchange)

	}(results)

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("Request terminated before publish")
	case err := <-results:
		return nil, err
	}

}
