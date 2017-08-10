package gremlin

import (
	"encoding/json"
	"errors"
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
	errorListener     map[string]chan error
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
	client.errorListener = make(map[string]chan error)
	client.ErrChannel = make(chan error)
	client.lock = &sync.Mutex{}
	return client
}

func (client *Client) SendRequest(gr *GremlinRequest) (interface{}, error) {

	data, err := gr.PackageRequest()
	id := gr.RequestID
	defer client.unregisterErrorListener(id)
	client.registerErrorListener(id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	client.requests <- data
	client.responseListeners[id] = make(chan bool)
	var response interface{}
	for {
		select {
		case resp := <-client.responseListeners[id]:
			if resp {
				log.Println(client.allResults[id])
				response = client.allResults[id]
			} else {
				err = errors.New("FAILED")
			}
			return response, err
		case err1 := <-client.errorListener[gr.RequestID]:
			err = err1
			log.Println(err)
			return response, err
		}
	}

}

func (client *Client) broadCastError(err error) {
	log.Println("broadcast error")
	for _, v := range client.errorListener {
		v <- err
	}
}

func (client *Client) registerErrorListener(id string) {
	client.errorListener[id] = make(chan error)
}

func (client *Client) unregisterErrorListener(id string) {
	delete(client.errorListener, id)
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
	defer func() {
		client.conn.Close()
		client.IsConnected = false
	}()
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
		case err := <-client.ErrChannel:
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Println("closed")
				client.IsConnected = false
			}
			client.broadCastError(err)
			log.Println(err)
			break
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
