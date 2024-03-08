package engine

import (
  "fmt"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  "go.mongodb.org/mongo-driver/bson"
  "context"
  "time"
  "os"
  "errors"
)

type ResponseData struct {
  Username string `json:"username"`
}


func (e *Engine) RegisterUserToDB(name string) (error){
  if name == ""{
    return errors.New("invalid name")
  }
  coll := e.dbConn.Database("chat-app").Collection("chatters")
  chatrooms := make([]string, 0)
  newChatter := Chatter{Username: name, CreatedAt: time.Now(), Chatrooms: chatrooms}
  result, err := coll.InsertOne(context.TODO(), newChatter)
  if err != nil{
    return errors.New("DB error")
  }
  fmt.Println("Inserted user in the db with id: ", result.InsertedID)
  return nil
}

func (e *Engine) LoginUser(name string) (error){
  if name == "" {
    return errors.New("invalid username")
  }
  coll := e.dbConn.Database("chat-app").Collection("chatters")
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

func (e *Engine) getHubsFromDB() ([]string, error){
  var results []string
  coll := e.dbConn.Database("chat-app").Collection("chatrooms")
  filter := bson.D{{}}
  opts := options.Find().SetProjection(bson.D{{"name", 1}})
  cursor, err := coll.Find(context.TODO(), filter, opts)
  if err != nil {
    return results, err
  }
  defer cursor.Close(context.Background())
  var elem struct {
    Name string `bson:"name"`
  }
  for cursor.Next(context.Background()){
    err := cursor.Decode(&elem)
    if err != nil {
      return results, err
    }
    results = append(results, elem.Name)
  }
  return results, nil
}


func dbInit() (*mongo.Client, error){
  uri, err := getURI()
  if err != nil{
    return nil, errors.New("Error fetching URI")
  }
  serverAPI := options.ServerAPI(options.ServerAPIVersion1)
  opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
  dbConn, err := mongo.Connect(context.TODO(), opts)
  if err != nil {
    panic(err)
  } 
  var pingResult bson.M
  if err := dbConn.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&pingResult); err != nil {
    return nil, err
  }
  fmt.Println(pingResult)
  
  return dbConn, nil
}
