package gremlin

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Connection struct {
	host                   string
	conn                   *websocket.Conn
	IsConnected            bool
	clientRequestChannel   chan GremlinRequest
	clientResponseListener chan GremlinResponse
	poolErrorListener      chan *Connection
	requestErrorListener   chan string
	errorListener          chan error
	CurrentRequestIds      map[string]bool
	lock                   *sync.Mutex
}

func NewConnection(host string, reqChannel chan GremlinRequest, respChannel chan GremlinResponse,
	poolErrorListener chan *Connection, requestErrorListener chan string) *Connection {
	connection := new(Connection)
	connection.host = "ws://" + host
	connection.clientRequestChannel = reqChannel
	connection.clientResponseListener = respChannel
	connection.poolErrorListener = poolErrorListener
	connection.requestErrorListener = requestErrorListener
	connection.IsConnected = false
	connection.CurrentRequestIds = make(map[string]bool)
	connection.errorListener = make(chan error)
	connection.lock = &sync.Mutex{}
	return connection
}

func (conn *Connection) Connect() error {
	d := websocket.Dialer{
		WriteBufferSize: 8192,
		ReadBufferSize:  8192,
	}
	var err error
	conn.conn, _, err = d.Dial(conn.host, http.Header{})
	if err == nil {
		conn.IsConnected = true
		go conn.writeHelper()
		go conn.readHelper()
	}
	return err
}

func (conn *Connection) errorAllRequests(err error) {
	for id, _ := range conn.CurrentRequestIds {
		conn.requestErrorListener <- id
	}
	conn.lock.Lock()
	for k := range conn.CurrentRequestIds {
		delete(conn.CurrentRequestIds, k)
	}
	conn.lock.Unlock()
}

func (conn *Connection) addRequestId(id string) {
	conn.lock.Lock()
	conn.CurrentRequestIds[id] = true
	conn.lock.Unlock()
}

func (conn *Connection) onRequestHandled(response GremlinResponse) {
	id := response.getRequestId()
	conn.lock.Lock()
	delete(conn.CurrentRequestIds, id)
	conn.lock.Unlock()
	log.Println("onRequestHAndled ", response)
	conn.clientResponseListener <- response
}

func (conn *Connection) writeHelper() {
	defer func() {
		conn.conn.Close()
		conn.IsConnected = false
	}()
	for {
		select {
		case gr := <-conn.clientRequestChannel:
			log.Println("conn recieved ", gr)
			if !conn.IsConnected {
				err := conn.Connect()
				if err != nil {
					log.Println(err)
					conn.poolErrorListener <- conn
					break
				}
			}
			log.Println(gr)
			id := gr.RequestID
			conn.CurrentRequestIds[id] = true
			data, err1 := gr.PackageRequest()
			if err1 != nil {
				conn.requestErrorListener <- id
			}
			err2 := conn.conn.WriteMessage(2, data)
			if err2 != nil {
				log.Println(err2)

			}
		case err := <-conn.errorListener:
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Println("closed")
				conn.IsConnected = false
			}
			conn.errorAllRequests(err)
			log.Println(err)
			conn.poolErrorListener <- conn
			return
		}
	}
}

func (conn *Connection) readHelper() {
	for {
		_, msg, err := conn.conn.ReadMessage()
		log.Println("recvd ", string(msg), err, conn.host)
		fmt.Println(string(msg))
		fmt.Println(err)
		if err != nil {
			log.Println(err)
			conn.errorListener <- err
			break
		}
		if msg != nil {
			var resp GremlinResponse
			err := json.Unmarshal(msg, &resp)
			if err != nil {
				log.Println(err)
				conn.errorListener <- err
			}
			conn.onRequestHandled(resp)
		}
	}
}
