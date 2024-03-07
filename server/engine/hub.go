package engine

import (
  "fmt"
)

type BroadcastMessage struct {
  sender []byte
  message []byte
}

func NewBroadcastMessage(sender string, message string) *BroadcastMessage {
  return &BroadcastMessage{
    sender: []byte(sender),
    message: []byte(message),
  }
}

func DecodeBroadcastMessage(b *BroadcastMessage) ([]byte, []byte) {
  return b.sender, b.message
}

type Hub struct {
  name string
  clients map[*ImClient]bool
  register chan *ImClient
  deregister chan *ImClient
  broadcast chan *BroadcastMessage
}


func newHub(name string) *Hub {
  fmt.Println("Creating a new hub...")
  return &Hub{
    name: name,
    clients: make(map[*ImClient]bool), 
    register: make(chan *ImClient),
    deregister: make(chan *ImClient),
    broadcast: make(chan *BroadcastMessage),
  }
}


func (hub *Hub) runHub(){
  for {
    select {
    case client := <- hub.register:
        client.setHub(hub)
        hub.clients[client] = true
        client.send <- NewBroadcastMessage(client.name, "Joined the hub")
    case client := <- hub.deregister:
        if _,ok := hub.clients[client]; ok {
            client.send <- NewBroadcastMessage(client.name, "Client left the hub")
            delete(hub.clients, client)
            client.setHub(nil)
        }
    case broadcastMessage := <- hub.broadcast:
      for client := range(hub.clients){
        select {
        case client.send <- broadcastMessage:
        default:
          close(client.send)
          delete(hub.clients, client)
        }
      }
    }
  }
}


func (hub *Hub) run(){
  go hub.runHub()
}
