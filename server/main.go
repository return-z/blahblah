package main

import (
  "flag"
  "engine"
  "api"
)

var addr = flag.String("addr", ":5990", "http service address")

func main(){
  flag.Parse()
  
  e := engine.NewEngine()
  e.Run()

  s := api.NewServer(*addr)
  s.Start(e)
}
