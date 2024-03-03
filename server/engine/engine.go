package engine

import (
  "github.com/gin-gonic/gin"
  "go.mongodb.org/mongo-driver/mongo"
  "log"
  "fmt"
)

type Engine struct {
  Hubs  map[string]*Hub
  dbConn *mongo.Client
}

func NewEngine() *Engine {
  return &Engine{
    Hubs: make(map[string]*Hub),
    dbConn: nil,
  }
}

func (e *Engine)setdbConn(dbConn *mongo.Client){
  e.dbConn = dbConn
}

func (e *Engine)setHubs(hubs []string){
  for _,hubName := range hubs{
    hub := newHub(hubName)
    e.Hubs[hubName] = hub 
    go hub.run()
  }
}

func (e *Engine)ServeWS(c *gin.Context){
  conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
  if err != nil {
    return
  }
  imClient := &ImClient{engine: e, name: loggedInUser, hub: nil, conn: conn, send: make(chan []byte, 256)}
  go imClient.socketReadPump()
  go imClient.socketWritePump()
}

func (e *Engine)Run() {
  fmt.Println("Initializing DB and retreiving connection string")
  dbConn, err := dbInit()
  if err != nil {
    fmt.Println("encountered an error")
    log.Fatal(err)
    return
  }
  e.setdbConn(dbConn)
  
  fmt.Println("Finding hubs from the DB")
  hubs, err := e.getHubsFromDB()
  if err != nil {
    log.Fatal(err)
  }
  e.setHubs(hubs)
}


