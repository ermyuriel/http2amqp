# http2amqp

 A simple Go server module to receive HTTP requests and route them to AMQP exchanges, for use with webhooks and other data integration strategies

## Configuration

Set up environment variables:

```bash
AMQP_HOST=amqp://ulbiqljr:OuyUmV@dddsf.rmq.cloudamqp.com/sdasd
RECEIVE_PATHS=/c,/m
AMQP_ROUTING_KEYS=calls,messages
AMQP_EXCHANGES=webhooks,webhooks
HEALTHCHECK_PATH=/healthcheck
RECEIVE_PORT=5000
```
Each path will receive HTTP requests and route them to the specified exchange with the specified routing key. Parameters are provided as comma separated lists.

In the provided example, http://server/c will send messages generated from non empty requests (i.e. ContentLength>0) to the 'webhooks' exchange with the routing key 'calls'. 

The generated AMQP message will conform to this structure:

```go
type HTTPRequest struct {
	Method    string      `json:"method"`
	Header    map[string][]string `json:"header"`
	Body      interface{} `json:"body"`
	Source    string      `json:"source"`
	Target    string      `json:"target"`
	Timestamp int64       `json:"timestamp"`
}
```

## Contributing
Just tell me what I did wrong.

## License
[MIT](https://choosealicense.com/licenses/mit/)
