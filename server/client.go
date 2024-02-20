package main

import (
  "fmt"
  "net/http"
  "bytes"
  "time"
  "github.com/gin-gonic/gin"
  "github.com/gorilla/websocket"
  "encoding/json"
  "strings"
)

var username string

type ImClient struct {
  hub *Hub
  conn *websocket.Conn
  send chan []byte
}

type ReceivedMessage struct {
  Msg string `json:"message"`
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
  ReadBufferSize: 1024,
  WriteBufferSize: 1024,
  CheckOrigin: func(r *http.Request) bool {
    return true
  },
}

func htmxized(b []byte) []byte{
  htmxMessage := [][]byte{[]byte(`<div id="messages" hx-swap-oob="beforeend">`), []byte(fmt.Sprintf(`<div class="flex p-1"><p class="text-green-400 font-mono">%s:</p> <p class="text-blue-400 font-mono ml-1">%s</p></div>`, username, string(b))), []byte("</p></div>")}
  return bytes.Join(htmxMessage, []byte(""))
}

var (
  newline = []byte{'\n'}
  space = []byte{' '}
)

func isJoinCommand(message string) (bool, *Hub) {
  args := strings.Fields(message)
  fmt.Println(args)
  if len(args) > 1  {
    if args[0] == "!join" {
      roomname := args[1]
      if hub, ok := hubs[roomname]; ok{
        return true, hub
      }
    }
  }
  return false, nil
}

func (c *ImClient) setHub(hub *Hub){
  c.hub = hub
  c.hub.register <- c
}

func (c *ImClient) socketReadPump(){
  defer func() {
    if c.hub != nil {
      c.hub.deregister <- c
    }
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
    fmt.Println(msg.Msg)
    if err != nil {
      fmt.Println(err)
      panic(err)
    }
    if c.hub != nil {
      c.hub.broadcast <- []byte(msg.Msg)
    } else {
      c.send <- []byte(msg.Msg)
    }
    if yes, hub := isJoinCommand(msg.Msg); yes {
      c.setHub(hub)
    }
  }
}

func (c *ImClient) socketWritePump(){
  ticker := time.NewTicker(pingPeriod)
  defer func() {
    ticker.Stop()
    c.conn.Close()
  }()
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
      w.Write(htmxized(message))
      n := len(c.send)
      for i:=0; i<n; i++{
        w.Write(newline)
        message := <-c.send
        w.Write(htmxized(message))
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
}

func serveWS(c *gin.Context){
  conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
  if err != nil {
    return
  }
  imClient := &ImClient{hub: nil, conn: conn, send: make(chan []byte, 256)}
  go imClient.socketReadPump()
  go imClient.socketWritePump()
}
