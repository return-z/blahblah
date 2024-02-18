package main

import (
  "fmt"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  "go.mongodb.org/mongo-driver/bson"
  "context"
  "time"
  "net/http"
  "encoding/json"
  "os"
  "github.com/gin-gonic/gin"
  "errors"
)

type Secret struct{
  URI string `json:"uri"`
}

var connDB *mongo.Client
var messagesDB chan []byte

type ResponseData struct {
  Username string `json:"username"`
}

func handleDB(){
  defer func() {
    if err := connDB.Disconnect(context.TODO()); err != nil {
      panic(err)
    }
  }()
  conn := connDB.Database("chat-app").Collection("messages")
  for{
    select{ 
    case message := <-messagesDB:
      strMessage := string(message)
      fmt.Println(strMessage)
      doc := Message{Body: strMessage, CreatedAt: time.Now()}
      result, err := conn.InsertOne(context.TODO(), doc)
      if err != nil {
        panic(err)
      }
      fmt.Println("Inserted message in db with _id: ", result.InsertedID)
    }
  }
}

func userAuthDB(c *gin.Context) ([]Chatroom, error){
  username := c.PostForm("username")
  if username := c.PostForm("username"); username == "" {
    return nil, errors.New("invalid username")
  }
  conn := connDB.Database("chat-app").Collection("chatters")
  filter := bson.M{"username": username}
  fmt.Println(filter)
  var res Chatter
  err := conn.FindOne(context.TODO(), filter).Decode(&res)
  fmt.Println(res)
  if err != nil{
    if err == mongo.ErrNoDocuments {
      c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
      return nil, errors.New("User not found")
    }
    panic(err)
  }
  fmt.Println(res.Chatrooms)
  collection := connDB.Database("chat-app").Collection("chatrooms")
  chatroomFilter := bson.M{"name" : bson.M{"$in": res.Chatrooms}}
  cursor, err := collection.Find(context.TODO(), chatroomFilter)
  if err != nil {
    fmt.Println(err)
    return nil, errors.New("Error retrieving chatrooms")
  }
  defer cursor.Close(context.TODO())

  var results []Chatroom
  if err = cursor.All(context.TODO(), &results); err != nil {
    return nil, errors.New("Error processing chatrooms")
  }

  fmt.Println(results)
  return results, nil
}

func getURI() string{
  secret, err := os.ReadFile("secret.json")
  if err != nil {
    fmt.Println("File doesn't exist, pwnie")
    return ""
  }
  var uri Secret
  if err := json.Unmarshal(secret, &uri); err != nil{
    fmt.Println("Error parsing uri: ", err)
    return ""
  }
  return uri.URI
}

func dbInit(){
  uri := getURI()
  serverAPI := options.ServerAPI(options.ServerAPIVersion1)
  opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
  conn, err := mongo.Connect(context.TODO(), opts)
  connDB = conn
  messagesDB = make(chan []byte)
  fmt.Println("Successfully connected to the DB")
  if err != nil {
    panic(err)
  }
  var pingResult bson.M
  if err := connDB.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&pingResult); err != nil {
    panic(err)
  }
  fmt.Println(pingResult)
  go handleDB()
}
