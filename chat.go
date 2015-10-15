package main

// import (
// 	"github.com/gin-gonic/gin"
// 	"gopkg.in/mgo.v2"
// 	"gopkg.in/mgo.v2/bson"

// 	"log"
// 	"time"
// )

// var (
// 	mongo = &Mongo{
// 		ConnectionString:  "localhost:27017",
// 		Database:          "chat",
// 		UserCollection:    "users",
// 		RoomCollection:    "rooms",
// 		MessageCollection: "messages",
// 	}
// )

// type User struct {
// 	Id     bson.ObjectId `bson:"_id"`
// 	Nick   string
// 	Secret string
// 	Email  string
// }

// func NewUser() *User {
// 	user := &User{}
// 	user.Id = bson.NewObjectId()
// 	return user
// }

// type Room struct {
// 	Id          bson.ObjectId `bson:"_id"`
// 	Name        string
// 	Description string
// 	Users       []User
// }

// func NewRoom() *Room {
// 	room := &Room{}
// 	room.Id = bson.NewObjectId()
// 	return room
// }

// type Message struct {
// 	Id      bson.ObjectId `bson:"_id"`
// 	User    string
// 	Room    string
// 	Message string
// 	Time    time.Time
// }

// func main() {

// 	err := mongo.OpenSession()

// 	if err != nil {
// 		log.Println(err)
// 	}

// 	r := gin.Default()
// 	r.POST("/registeruser", registerUser)
// 	r.POST("/createroom", createRoom)
// 	r.POST("/sendmessage", sendMessage)
// 	r.POST("/getmessages", getMessagesForRoom)
// 	r.GET("/getrooms", getRooms)
// 	r.Run(":8080") // listen and serve on 0.0.0.0:8080
// }

// type DataStorage interface {
// 	OpenSession() error
// 	CloseSession()
// 	SaveUser(user *User) error
// 	FindUser(userId string) (*User, error)
// 	SaveRoom(room *Room) error
// 	FindRoom(roomId string) (*Room, error)
// 	SaveMessage(message *Message) error
// }

// type Mongo struct {
// 	ConnectionString  string
// 	Database          string `validate:"nonzero"`
// 	UserCollection    string `validate:"nonzero"`
// 	RoomCollection    string `validate:"nonzero"`
// 	MessageCollection string `validate:"nonzero"`
// 	mongoSession      *mgo.Session
// 	users             *mgo.Collection
// 	rooms             *mgo.Collection
// 	messages          *mgo.Collection
// 	db                *mgo.Database
// }

// //Opens mongo session for given url and
// //sets the session to global property
// func (m *Mongo) OpenSession() error {

// 	mongo, connError := mgo.Dial(m.ConnectionString)
// 	m.mongoSession = mongo
// 	if connError != nil {
// 		return connError
// 	}

// 	m.db = mongo.DB(m.Database)
// 	m.users = m.db.C(m.UserCollection)
// 	m.rooms = m.db.C(m.RoomCollection)
// 	m.messages = m.db.C(m.MessageCollection)

// 	//Indexes
// 	m.messages.EnsureIndexKey("room")
// 	m.messages.EnsureIndexKey("time")

// 	return nil
// }

// //Silently close mongo session
// func (m *Mongo) CloseSession() {
// 	m.mongoSession.Close()
// }

// func (m *Mongo) SaveUser(user *User) (*User, error) {
// 	err := m.users.Insert(user)
// 	return user, err
// }

// func (m *Mongo) FindUser(userId string) (*User, error) {
// 	user := &User{}
// 	err := m.users.FindId(userId).One(user)
// 	return user, err
// }

// func (m *Mongo) SaveRoom(room *Room) (*Room, error) {
// 	err := m.rooms.Insert(room)
// 	return room, err
// }

// func (m *Mongo) FindRoom(roomId string) (*Room, error) {
// 	room := &Room{}
// 	err := m.rooms.FindId(roomId).One(room)
// 	return room, err
// }

// func (m *Mongo) FindRoomsWithUsers() ([]Room, error) {
// 	query := []bson.M{
// 		{"$project": bson.M{"name": 1, "users": 1}},
// 	}
// 	qryRes := make([]Room, 0)
// 	err := m.rooms.Pipe(query).All(&qryRes)
// 	return qryRes, err
// }

// // func (m *Mongo)

// func (m *Mongo) SaveMessage(message *Message) error {
// 	message.Id = bson.NewObjectId()
// 	message.Time = bson.Now()
// 	err := m.messages.Insert(message)
// 	return err
// }

// func (m *Mongo) FindMessagesByRoom(roomName string, limit int) ([]Message, error) {
// 	messages := make([]Message, 0)
// 	qry := m.messages.Find(bson.M{"room": roomName})
// 	qry = qry.Sort("-time")
// 	qry = qry.Limit(limit)
// 	err := qry.All(&messages)
// 	return messages, err
// }

// func registerUser(ctx *gin.Context) {
// 	user := NewUser()
// 	var err error
// 	err = ctx.BindJSON(user)
// 	if err == nil {
// 		_, err = mongo.SaveUser(user)
// 	}

// 	if err != nil {
// 		log.Println(err)
// 		ctx.JSON(405, err)
// 	} else {
// 		ctx.JSON(200, user)
// 	}
// }

// func createRoom(ctx *gin.Context) {
// 	//TODO validate if no such room
// 	room := NewRoom()
// 	var err error
// 	err = ctx.BindJSON(room)
// 	if err == nil {
// 		_, err = mongo.SaveRoom(room)
// 	}

// 	if err != nil {
// 		log.Println(err)
// 		ctx.JSON(405, err)
// 	} else {
// 		ctx.JSON(200, room)
// 	}
// }

// func getRooms(ctx *gin.Context) {
// 	res, err := mongo.FindRoomsWithUsers()

// 	roomDtos := make([]RoomDTO, len(res))

// 	for i, room := range res {
// 		roomDtos[i] = RoomDTO{
// 			room.Name,
// 			len(room.Users),
// 		}
// 	}

// 	if err != nil {
// 		log.Println(err)
// 		ctx.JSON(405, err)
// 	} else {
// 		ctx.JSON(200, roomDtos)
// 	}
// }

// func sendMessage(ctx *gin.Context) {
// 	message := &Message{}
// 	ctx.BindJSON(message)

// 	err := mongo.SaveMessage(message)
// 	if err != nil {
// 		log.Println(err)
// 		ctx.JSON(405, err)
// 	} else {
// 		ctx.JSON(200, "")
// 	}

// }

// type RoomDTO struct {
// 	Name      string
// 	UserCount int
// }

// func getMessagesForRoom(ctx *gin.Context) {
// 	room := &RoomDTO{}
// 	ctx.BindJSON(room)
// 	result, err := mongo.FindMessagesByRoom(room.Name, 50)
// 	if err != nil {
// 		log.Println(err)
// 		ctx.JSON(405, err)
// 	} else {
// 		ctx.JSON(200, result)
// 	}
// }
