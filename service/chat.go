package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	uuid "github.com/satori/go.uuid"

	websocket "github.com/gorilla/websocket"
)

type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	id     string
	socket *websocket.Conn
	send   chan []byte
}

type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

var manager = ClientManager{
	broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[*Client]bool),
}

func (manager *ClientManager) start() {
	fmt.Println("start the start fun")

	for {
		select {
		case conn := <-manager.register:
			fmt.Println("register conn:", conn)

			manager.clients[conn] = true
			// jsonMessage, _ := json.Marshal(&Message{Content: "/A new socket has connected."})
			jsonMessage, _ := json.Marshal(&Message{Content: "/一个新时空隧道已连接."})
			manager.send(jsonMessage, conn)
		case conn := <-manager.unregister:
			fmt.Println("unregister conn:", conn)

			if _, ok := manager.clients[conn]; ok {
				close(conn.send)
				delete(manager.clients, conn)
				// jsonMessage, _ := json.Marshal(&Message{Content: "/A socket has disconnected."})
				jsonMessage, _ := json.Marshal(&Message{Content: "/一个时空隧道连接断开."})
				manager.send(jsonMessage, conn)
			}
		case message := <-manager.broadcast:
			fmt.Println("broadcast message:", message)

			for conn := range manager.clients {
				select {
				case conn.send <- message:
				default:
					fmt.Println("broadcast message: default!!!")

					close(conn.send)
					delete(manager.clients, conn)
				}
			}
		}
	}
}

func (manager *ClientManager) send(message []byte, ignore *Client) {
	for conn := range manager.clients {
		if conn != ignore {
			conn.send <- message
		}
	}
}

func (c *Client) read() {
	defer func() {
		manager.unregister <- c
		c.socket.Close()
	}()

	for {
		fmt.Println("In read for loop")

		_, message, err := c.socket.ReadMessage()
		fmt.Println("In read for loop: after init message:", message)

		if err != nil {
			manager.unregister <- c
			c.socket.Close()
			fmt.Println("In read for loop: err != nil")

			break
		}
		fmt.Println("In read for loop: before get jsonMessage")

		jsonMessage, _ := json.Marshal(&Message{Sender: c.id, Content: string(message)})
		fmt.Println("In read for loop jsonMessage:", jsonMessage)

		manager.broadcast <- jsonMessage
	}
}

func (c *Client) write() {
	defer func() {
		c.socket.Close()
	}()

	for {
		fmt.Println("In write for loop")

		select {
		case message, ok := <-c.send:
			fmt.Println("In write for loop: init message, ok:", message, ok)

			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				fmt.Println("In write for loop: !ok ")
				return
			}

			c.socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func main() {
	fmt.Println("Starting application...")
	go manager.start()
	fmt.Println("After start")

	http.HandleFunc("/myws", wsPage)
	http.ListenAndServe("127.0.0.1:12345", nil)

	fmt.Println("End of main")

}

func wsPage(res http.ResponseWriter, req *http.Request) {
	fmt.Println("start the wsPage fun") // "RES:", res, "REQ:", req,

	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if error != nil {
		http.NotFound(res, req)
		return
	}
	client := &Client{id: uuid.NewV4().String(), socket: conn, send: make(chan []byte)}

	manager.register <- client

	go client.read()
	go client.write()

	fmt.Println("After read and write")

}
