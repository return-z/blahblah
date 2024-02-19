package main

import (
  "fmt"
)


type ServerHub struct {
  clients map[*ImClient]bool
  register chan *ImClient
  deregister chan *ImClient
}

func newServerHub() *ServerHub {
  fmt.Println("Creating a new server hub...")
  return &ServerHub{
    clients: make(map[*ImClient]bool), 
    register: make(chan *ImClient),
    deregister: make(chan *ImClient),
  }
}


func (serverHub *ServerHub) runServerHub(){
  for {
    select {
    case client := <- serverHub.register:
        serverHub.clients[client] = true
    case client := <- serverHub.deregister:
        if _,ok := serverHub.clients[client]; ok {
            delete(serverHub.clients, client)
        }
    }
  }
}

