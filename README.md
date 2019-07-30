# http2amqp

 A simple Go AWS lambda module to receive HTTP requests and route them to AMQP exchanges, for use with webhooks and other data integration strategies.

## Configuration

Set up environment variables:

```bash
AMQP_HOST=amqp://ulbiqljr:OuyUmV@dddsf.rmq.cloudamqp.com/sdasd
```
Each path will receive HTTP requests and route them to the specified exchange and routing key, passed in as URL parameters

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
