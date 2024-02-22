package main

import (
  "fmt"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  "go.mongodb.org/mongo-driver/bson"
  "context"
  "time"
  "os"
  "errors"
  "github.com/joho/godotenv"
)


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

func registerUser(name string) (error){
  if name == ""{
    return errors.New("invalid name")
  }
  coll := connDB.Database("chat-app").Collection("chatters")
  chatrooms := make([]string, 0)
  newChatter := Chatter{Username: name, CreatedAt: time.Now(), Chatrooms: chatrooms}
  result, err := coll.InsertOne(context.TODO(), newChatter)
  if err != nil{
    return errors.New("DB error")
  }
  fmt.Println("Inserted user in the db with id: ", result.InsertedID)
  return nil
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

func getURI() (string, error){
  return os.Getenv("URI"), nil
}

func dbInit() (error){
  uri, err := getURI()
  if err != nil{
    return errors.New("Error fetching URI")
  }
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
  
  return nil
}
