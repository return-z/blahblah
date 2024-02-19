package main

import (
  "fmt"
  "bytes"
  "time"
  "github.com/gin-gonic/gin"
  "github.com/gorilla/websocket"
  "encoding/json"
  "strings"
  "errors"
)

func parseCommand(message string) (string, []string, error){
  if !strings.HasPrefix(message, "!"){
    return "", nil, errors.New("Not a command")
  }
  parts := strings.Fields(message)
  if len(parts) == 0 {
    return "", nil, errors.New("No command found")
  }
  command := parts[0][1:]
  var args []string
  if len(parts) > 1 {
    args = parts[1:]
  }
  return command, args, nil
}

func handleCommand(command string, args []string, c *ImClient) (*Hub, error){
  switch command {
    case "join":
      roomName := args[0]
      if hub, ok := hubs[roomName]; ok{
          fmt.Println("Trying to join hub: ", roomName)
          return hub, nil
      } else {
        return nil, errors.New("Not a valid hub")    
      }
    default:
      return nil, errors.New("Something went wrong with the command")
  }
}

func (c *ImClient) serverSocketReadPump(){
  defer func() {
    c.serverHub.deregister <- c
    c.conn.Close()
  }()
  c.conn.SetReadLimit(maxMessageSize)
  c.conn.SetReadDeadline(time.Now().Add(pongWait))
  c.conn.SetPongHandler(func(string) error {c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil})
  for{
    _, message, err := c.conn.ReadMessage()
    if err != nil {
      return
    }
    message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
    var msg ReceivedMessage 
    err = json.Unmarshal(message, &msg)
    if err != nil {
      fmt.Println(err)
      panic(err)
    }
    fmt.Println(msg.Msg)
    c.send <- []byte(msg.Msg)
  }
}

func (c *ImClient) serverSocketWritePump(){
  ticker := time.NewTicker(pingPeriod)
  for{
    select{
    case message,ok := <-c.send:
      c.conn.SetWriteDeadline(time.Now().Add(writeWait)) 
      if !ok{
        c.conn.WriteMessage(websocket.CloseMessage, []byte{})
        return
      }
      w,err := c.conn.NextWriter(websocket.TextMessage)
      if err != nil{
        return
      }
      c.mu.Lock()
      w.Write(htmxized(message))
      c.mu.Unlock()
      if command, args, err := parseCommand(string(message)); err == nil{ 
        if hub, err := handleCommand(command, args, c); err == nil {
          go joinHub(c, hub) 
        }
      }
      if err := w.Close(); err != nil{
        return
      }
    case <-ticker.C:
      c.conn.SetWriteDeadline(time.Now().Add(writeWait))
      if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
        return
      }
    }
  }
  defer func() {
    ticker.Stop()
    c.conn.Close()
  }()
}

func serveWS(serverHub *ServerHub, c *gin.Context){
  conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
  if err != nil {
    return
  }
  imClient := &ImClient{serverHub: serverHub, hub: nil, conn: conn, send: make(chan []byte, 256)}
  imClient.serverHub.register <- imClient
  go imClient.serverSocketReadPump()
  go imClient.serverSocketWritePump()
}
