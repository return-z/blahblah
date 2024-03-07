package engine

import (
  "fmt"
  "net/http"
  "bytes"
  "time"
  "github.com/gorilla/websocket"
  "encoding/json"
)


type ImClient struct {
  engine *Engine
  name string
  hub *Hub
  conn *websocket.Conn
  send chan *BroadcastMessage
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

func htmxizedMessage(broadcastMessage *BroadcastMessage) []byte{
  sender, message := DecodeBroadcastMessage(broadcastMessage) 
  htmxMessage := [][]byte{
    []byte(`<div id="messages" hx-swap-oob="beforeend">`), 
    []byte(fmt.Sprintf(`<div class="flex p-1"><p class="text-green-400 font-mono">%s:</p> <p class="text-blue-400 font-mono ml-1">%s</p></div>`, string(sender), string(message))), 
    []byte("</p></div>")}
  return bytes.Join(htmxMessage, []byte(""))
}

var (
  newline = []byte{'\n'}
  space = []byte{' '}
)


func (c *ImClient) setHub(hub *Hub){
  c.hub = hub
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
      c.hub.broadcast <- NewBroadcastMessage(c.name, msg.Msg) 
    } else {
      c.send <- NewBroadcastMessage(c.name, msg.Msg) 
    }
    c.parseCommand(msg.Msg)
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
    case broadcastMessage,ok := <-c.send:
      c.conn.SetWriteDeadline(time.Now().Add(writeWait)) 
      if !ok{
        c.conn.WriteMessage(websocket.CloseMessage, []byte{})
        return
      }
      w,err := c.conn.NextWriter(websocket.TextMessage)
      if err != nil{
        return
      }
      w.Write(htmxizedMessage(broadcastMessage))
      n := len(c.send)
      for i:=0; i<n; i++{
        w.Write(newline)
        broadcastMessage := <-c.send
        w.Write(htmxizedMessage(broadcastMessage))
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

