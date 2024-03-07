package api

import (
  "net/http"
  "github.com/gin-gonic/gin"
  "assets"
  "github.com/a-h/templ/examples/integration-gin/gintemplrenderer"
  "engine"
)

type Server struct {
  address string
}

func NewServer(address string) *Server {
  return &Server{
    address: address,
  }
}

func (s *Server) Start (engine *engine.Engine){ 
  
  router := gin.Default()
  router.Use(CORSMiddleware())
  
  ginHTMLRenderer := router.HTMLRender
  router.HTMLRender = &gintemplrenderer.HTMLTemplRenderer{FallbackHtmlRenderer: ginHTMLRenderer}

  router.GET("/", func(c *gin.Context){
    c.HTML(http.StatusOK, "", assets.Home())
  })
  router.GET("/login", func(c *gin.Context){
    c.HTML(http.StatusOK, "", assets.LoginForm(nil))
  })
  router.GET("/register", func(c *gin.Context){
    c.HTML(http.StatusOK, "", assets.RegisterForm(nil))
  })
  router.GET("/ws/:username", func(c *gin.Context){
    engine.ServeWS(c)
  })

  router.POST("/register", func(c *gin.Context){
    name := c.PostForm("username")
    err := engine.RegisterUserToDB(name)
    if err != nil{
      c.HTML(http.StatusOK, "", assets.RegisterForm(err))
    } else {
      c.HTML(http.StatusOK, "", assets.LoginForm(nil))
    }
  })
  router.POST("/login", func(c *gin.Context){
    name := c.PostForm("username")
    err := engine.LoginUser(name)
    if err != nil{
      c.HTML(http.StatusOK, "", assets.LoginForm(err))
    } else {
      c.HTML(http.StatusOK, "", assets.Chat(name))
    }
  })

  router.Run(s.address)
}
