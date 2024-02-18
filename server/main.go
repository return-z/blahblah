package main

import (
  "fmt"
  "flag"
  "github.com/gin-gonic/gin"
  "github.com/a-h/templ/examples/integration-gin/gintemplrenderer"
  "net/http"
)

var addr = flag.String("localhost", ":5990", "http service address")
var hubs map[string]*Hub


func CORSMiddleware() gin.HandlerFunc {
  return func(c *gin.Context) {
    c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
    c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
    c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
    c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

    if c.Request.Method == "OPTIONS" {
      c.AbortWithStatus(204)
            return
    }

    c.Next()
  }
}

func _createHubs(){
  hubs = make(map[string]*Hub)
  for i:=1; i<=4; i++{
    chatroomName := fmt.Sprintf("test_chatroom_%d", i)
    hub := newHub()
    hubs[chatroomName] = hub
  }
}

func main(){
  _createHubs()
  go dbInit()
  router := gin.Default()
  router.Use(CORSMiddleware())
  router.GET("/", func(c *gin.Context){
    r := gintemplrenderer.New(c.Request.Context(), http.StatusOK, Auth())
    c.Render(http.StatusOK, r)
  })
  router.GET("/chat", func(c *gin.Context){
    r := gintemplrenderer.New(c.Request.Context(), http.StatusOK, Chat())
    c.Render(http.StatusOK, r)
  })
  router.GET("/ws/:chatroom", func(c *gin.Context){
    chatroom := c.Param("chatroom")
    if hub, ok := hubs[chatroom]; ok {
      go hub.run()
      serveWS(hub, c)
    } else{
      c.JSON(http.StatusNotFound, gin.H{"message":"Chatroom not found!"})
    }
  })
  router.POST("/auth", func(c *gin.Context){
    results, err := userAuthDB(c)
    if err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"message": err})
      return
    }
    r := gintemplrenderer.New(c.Request.Context(), http.StatusOK, Chatrooms(results))
    c.Render(http.StatusOK, r)
  })
  router.Run("localhost:5990")
}
