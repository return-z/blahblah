package engine

import (
  "fmt"
  "errors"
  "strings"
)


func (c *ImClient)parseCommand(message string) (error){
  if !strings.HasPrefix(message, "!"){
    return nil
  }
  args := strings.Fields(message)
  fmt.Println(args)
  switch cmd := args[0]; cmd {
    case "!join":
      if c.hub != nil {
        return errors.New("Already in a hub")
      }
      if len(args) > 1 {
        chatroom := args[1]
        if hub, ok := c.engine.Hubs[chatroom]; ok {
          hub.register <- c
        }
      }
      return errors.New("The chatroom entered does not exist")
    case "!leave":
      fmt.Println("Trying to leave")
      if c.hub != nil {
        c.hub.deregister <- c
      }
      return errors.New("Not joined to a chatroom yet!")
    case "!create":
      fmt.Println("Creating a chatroom")
      if len(args) > 1 {
        chatroom := args[1]
        if _, ok := c.engine.Hubs[chatroom]; ok {
          return errors.New("The chatroom already exists")
        }
        hub := newHub(chatroom)
        c.engine.Hubs[chatroom] = hub
        hub.run()
      }
    case "!rooms":
      rooms := make([]string, 0, len(c.engine.Hubs))
      for hub := range c.engine.Hubs {
        rooms = append(rooms, hub)
      }
      c.send <- []byte(strings.Join(rooms, ", "))
    default:
      return nil
  }
  return nil
}
