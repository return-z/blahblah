package main

import (
  "fmt"
  "flag"
  "github.com/gin-gonic/gin"
  "github.com/a-h/templ/examples/integration-gin/gintemplrenderer"
  "net/http"
  "errors"
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
    go hub.run()
    hubs[chatroomName] = hub
  }
}

func main(){
  err := dbInit()
  if err != nil {
    fmt.Println("Problem initializing DB")
    return
  }
  _createHubs()
  gin.SetMode(gin.ReleaseMode)
  router := gin.Default()
  router.Use(CORSMiddleware())
  router.Static("/assets", "./")
  router.GET("/", func(c *gin.Context){
    r := gintemplrenderer.New(c.Request.Context(), http.StatusOK, Home())
    c.Render(http.StatusOK, r)
  })
  router.GET("/chat", func(c *gin.Context){
    r := gintemplrenderer.New(c.Request.Context(), http.StatusOK, Home())
    c.Render(http.StatusOK, r)
  })
  router.GET("/login", func(c *gin.Context){
    r := gintemplrenderer.New(c.Request.Context(), http.StatusOK, LoginForm(nil))
    c.Render(http.StatusOK, r)
  })
  router.GET("/register", func(c *gin.Context){
    r := gintemplrenderer.New(c.Request.Context(), http.StatusOK, RegisterForm(nil))
    c.Render(http.StatusOK, r)
  })
  router.GET("/ws", func(c *gin.Context){
    serveWS(c)
  })
  router.POST("/register", func(c *gin.Context){
    name := c.PostForm("username")
    err := userAuthDB(name)
    fmt.Println(err)
    if err != nil{
      registerUser(name)
      r := gintemplrenderer.New(c.Request.Context(), http.StatusOK,LoginForm(nil))
      c.Render(http.StatusOK, r)
    } else {
      r := gintemplrenderer.New(c.Request.Context(), http.StatusOK,  RegisterForm(errors.New("User already exists")))
      c.Render(http.StatusOK, r)
    }
  })
  router.POST("/auth", func(c *gin.Context){
    name := c.PostForm("username")
    err := userAuthDB(name)
    if err != nil{
      r := gintemplrenderer.New(c.Request.Context(), http.StatusOK, LoginForm(err))
      c.Render(http.StatusOK, r)
    } else {
      username = name
      r := gintemplrenderer.New(c.Request.Context(), http.StatusOK, Chat())
      c.Render(http.StatusOK, r)
    }
  })
  router.Run(":5990")
}
