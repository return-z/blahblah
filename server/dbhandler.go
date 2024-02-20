package main

import (
  "fmt"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  "go.mongodb.org/mongo-driver/bson"
  "context"
  "time"
  "encoding/json"
  "os"
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

func userAuthDB(name string) (error){
  if name == "" {
    return errors.New("invalid username")
  }
  coll := connDB.Database("chat-app").Collection("chatters")
  filter := bson.M{"username": name}
  fmt.Println(filter)
  var res Chatter
  err := coll.FindOne(context.TODO(), filter).Decode(&res)
  fmt.Println(res)
  if err != nil{
    if err == mongo.ErrNoDocuments {
      return errors.New("User not found")
    }
  }
  return nil
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
