package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/eucj/amqphelper"
	"github.com/eucj/gostructs"
)

type serverQueue struct {
	Queue      *amqphelper.Queue
	Path       string
	Exchange   string
	RoutingKey string
}

func (s *serverQueue) serve(w http.ResponseWriter, r *http.Request) {
	var body interface{}
	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logRequest(err.Error(), r)
		return
	}

	toSend := &gostructs.HTTPRequest{Method: r.Method, Header: r.Header, Body: body, Source: r.RemoteAddr, Target: r.RequestURI, Timestamp: time.Now().Unix()}
	message, err := json.Marshal(toSend)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logRequest(err.Error(), r)
		return
	}

	err = s.Queue.Publish(string(message), "application/json", s.Exchange, s.RoutingKey, true, false)

	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		log.Println(err)
		return
	}

	w.Write([]byte(fmt.Sprintf("{\"Exchange\":\"%s\",\"RoutingKey\":\"%s\"}", s.Exchange, s.RoutingKey)))
	w.WriteHeader(http.StatusOK)
}
func logRequest(prefix string, r *http.Request) string {
	req, _ := httputil.DumpRequest(r, true)
	lm := fmt.Sprintf("%s:\n%s", prefix, string(req))
	log.Println(lm)
	return lm

}
