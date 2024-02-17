package main

import (
  "time"
//  "go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
//  ChatterID primitive.ObjectID `bson:"chatter_id"`
//  ChatroomID primitive.ObjectID `bson:"chatroom_id"`
  Body string `bson:"body"`
  CreatedAt time.Time `bson:"created_at"`
}


type Chatter struct {
// ChatterId primitive.ObjectID `bson:"_id,omitempty"`
  Username string `bson:"username"`
  CreatedAt time.Time `bson:"created_at"`
  Chatrooms []string `bson:"chatrooms"`
}

type Chatroom struct {
//  ChatroomId primitive.ObjectID  `bson:"_id,omitempty"`
  Name string `bson:"name"`  
  CreatedAt time.Time `bson:"created_at"`
//  Chatters []primitive.ObjectID `bson:"chatters"` 
}


