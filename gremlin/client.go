package gremlin

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
)

type Client struct {
	host              string
	conn              *websocket.Conn
	requests          chan []byte
	allResults        map[string]interface{}
	responseListeners map[string]chan bool
	ErrChannel        chan error
	lock              *sync.Mutex
	IsConnected       bool
}

func NewClient(host string) *Client {
	client := new(Client)
	client.host = "ws://" + host
	client.IsConnected = false
	client.requests = make(chan []byte)
	client.allResults = make(map[string]interface{})
	client.responseListeners = make(map[string]chan bool)
	client.lock = &sync.Mutex{}
	return client
}

func (client *Client) SendRequest(gr *GremlinRequest) (interface{}, error) {
	//log.Println(msg)
	data, err := gr.PackageRequest()
	id := gr.RequestID
	if err != nil {
		log.Println(err)
		return nil, err
	}
	client.requests <- data
	client.responseListeners[id] = make(chan bool)
	_ = <-client.responseListeners[id]
	log.Println(client.allResults[id])
	return client.allResults[id], err
}

func (client *Client) OnResponse(data []byte) {
	var resp *GremlinResponse
	resp = new(GremlinResponse)
	err := json.Unmarshal(data, resp)
	if err != nil {
		log.Println(err)
		client.ErrChannel <- err
		return
	}
	id := resp.getRequestId()
	if client.responseListeners[id] == nil {
		client.responseListeners[id] = make(chan bool)
	}
	status := resp.getStatusCode()
	if status == 200 || status == 304 || status == 204 || status == 206 {
		client.allResults[id] = resp.GetResultData()
		client.responseListeners[id] <- true
	} else {
		client.responseListeners[id] <- false
	}

}

func (client *Client) NewGremlinRequest() *GremlinRequest {
	gremlinRequest := new(GremlinRequest)
	gremlinRequest.RequestID = uuid.NewV4().String()
	gremlinRequest.Op = "eval"
	gremlinRequest.Processor = ""
	var args Arguments
	args.Language = "gremlin-groovy"
	args.Bindings = make(map[string]interface{})
	args.Gremlin = ""
	gremlinRequest.Args = args
	return gremlinRequest
}

func (client *Client) writeHelper() {
	for {
		select {
		case msg := <-client.requests:
			if !client.IsConnected {
				err := client.Connect()
				if err != nil {
					log.Println(err)
					client.ErrChannel <- err
					break
				}
			}
			log.Println(string(msg))
			err2 := client.conn.WriteMessage(2, msg)
			if err2 != nil {
				log.Println(err2)
				break
			}
		}
	}
}

func (client *Client) readHelper() {
	for {
		_, msg, err := client.conn.ReadMessage()
		fmt.Println(string(msg))
		fmt.Println(err)
		if err != nil {
			log.Println(err)
			client.ErrChannel <- err
			break
		}
		if msg != nil {
			client.OnResponse(msg)
		}
	}
}

func (client *Client) Connect() (err error) {
	d := websocket.Dialer{
		WriteBufferSize: 8192,
		ReadBufferSize:  8192,
	}
	client.conn, _, err = d.Dial(client.host, http.Header{})
	if err == nil {
		client.IsConnected = true
		go client.writeHelper()
		go client.readHelper()
	}
	return err
}
