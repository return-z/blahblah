package main

import (
  "fmt"
)


type Hub struct {
  clients map[*ImClient]bool
  register chan *ImClient
  deregister chan *ImClient
  broadcast chan []byte
}


func newHub() *Hub {
  fmt.Println("Creating a new hub...")
  return &Hub{
    clients: make(map[*ImClient]bool), 
    register: make(chan *ImClient),
    deregister: make(chan *ImClient),
    broadcast: make(chan []byte),
  }
}


func (hub *Hub) runHub(){
  for {
    select {
    case client := <- hub.register:
        client.setHub(hub)
        hub.clients[client] = true
        client.send <- []byte("Joined the hub!")
    case client := <- hub.deregister:
        if _,ok := hub.clients[client]; ok {
            client.send <- []byte(fmt.Sprintf("%s left the hub", username))
            delete(hub.clients, client)
            client.setHub(nil)
        }
    case message := <- hub.broadcast:
      messagesDB <- message 
      fmt.Println("Put message in db")
      for client := range(hub.clients){
        select {
        case client.send <- message:
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
  go handleDB()
}
