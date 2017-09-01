package gremlin

import (
	"errors"
	"log"
	"strconv"
	"sync"

	uuid "github.com/satori/go.uuid"
)

type Client struct {
	hosts                []string
	connections          []*Connection
	requests             chan GremlinRequest
	responseChannel      chan GremlinResponse
	poolErrorListener    chan *Connection
	requestErrorListener chan string
	allResults           map[string]interface{}
	responseListeners    map[string]chan bool
	errorListener        map[string]chan error
	lock                 *sync.Mutex
}

func NewClient(hosts []string) *Client {
	client := new(Client)
	client.hosts = hosts
	client.requests = make(chan GremlinRequest)
	client.responseChannel = make(chan GremlinResponse)
	client.poolErrorListener = make(chan *Connection)
	client.requestErrorListener = make(chan string)
	client.allResults = make(map[string]interface{})
	client.responseListeners = make(map[string]chan bool)
	client.errorListener = make(map[string]chan error)
	client.lock = &sync.Mutex{}
	client.connections = make([]*Connection, len(hosts))
	for _, host := range hosts {
		connection := NewConnection(host, client.requests, client.responseChannel,
			client.poolErrorListener, client.requestErrorListener)
		connection.Connect()
		client.connections = append(client.connections, connection)
	}
	go client.messageReciever()
	return client
}

func (client *Client) SendRequest(gr *GremlinRequest) (interface{}, error) {

	id := gr.RequestID
	defer func() {
		client.unregisterErrorListener(id)
		if _, ok := client.allResults[id]; ok {
			delete(client.allResults, id)
		}
		if _, ok := client.responseListeners[id]; ok {
			close(client.responseListeners[id])
			delete(client.responseListeners, id)
		}
	}()
	client.registerErrorListener(id)
	client.responseListeners[id] = make(chan bool)
	client.requests <- *gr
	var response interface{}
	var err error
	for {
		select {
		case resp := <-client.responseListeners[id]:
			if resp {
				log.Println(client.allResults[id])
				response = client.allResults[id]
			} else {
				errCodde := strconv.FormatFloat(client.allResults[id].(float64), 'g', -1, 64)
				err = errors.New("Got status code " + errCodde)
			}

			return response, err
		case err1 := <-client.errorListener[gr.RequestID]:
			err = err1
			log.Println(err)
			return response, err
		}
	}
}

func (client *Client) messageReciever() {
	for {
		select {
		case requestId := <-client.requestErrorListener:
			client.sendErrorToRequest(requestId, errors.New("Failed to send request"))
		case resp := <-client.responseChannel:
			id := resp.getRequestId()
			if client.responseListeners[id] == nil {
				client.responseListeners[id] = make(chan bool)
			}
			status := resp.getStatusCode()
			if status == 200 || status == 304 || status == 204 || status == 206 {
				client.lock.Lock()
				client.allResults[id] = resp.GetResultData()
				client.lock.Unlock()
				client.responseListeners[id] <- true
			} else {
				client.lock.Lock()
				client.allResults[id] = resp.getStatusCode()
				client.lock.Unlock()
				client.responseListeners[id] <- false
			}
		case channelErr := <-client.poolErrorListener:
			log.Println(channelErr)
		}
	}

}

func (client *Client) sendErrorToRequest(id string, err error) {
	if v, ok := client.errorListener[id]; ok {
		v <- err
	}
}

func (client *Client) registerErrorListener(id string) {
	client.errorListener[id] = make(chan error)
}

func (client *Client) unregisterErrorListener(id string) {
	if _, ok := client.errorListener[id]; ok {
		close(client.errorListener[id])
		delete(client.errorListener, id)
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

/*func (client *Client) writeHelper() {
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
}*/
