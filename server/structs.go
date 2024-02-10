package rtfServer

import (
	"database/sql"
	"time"

	"github.com/gorilla/websocket"
)

var PORT string = "8001"

var (
	dbPath = "db/database.db"
	db     *sql.DB
)

// stores connected user and their corresponding web sockets
var connectedUsers map[string]*websocket.Conn

// stores additional datas about connected user
// should have included websocket from previous map
var connectedUsersData map[string]*conUsrData

type conUsrData struct {
	lastPostVisited int
}

// used to parse messages sent from client.
// contains all possible fields used
type Message struct {
	Type           string `json:"type"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	Text           string `json:"text"`
	To             string `json:"to"`
	Me             string `json:"me"`
	Recipient      string `json:"recipient"`
	Title          string `json:"title"`
	Content        string `json:"content"`
	Post_id        int    `json:"post_id"`
	Last_id        int    `json:"last_id"`
	Time           string `json:"time"`
	Email          string `json:"email"`
	Lastmessage_id int    `json:"lastmessage_id"`
}

// all different structures used to send data to client
// structure names same as types

type User struct {
	Type          string      `json:"type"`
	Data          interface{} `json:"data"`
	Userconnected string      `json:"userconnected"`
}

type UserStatus struct {
	Username    string `json:"username"`
	Connected   bool   `json:"connected"`
	Lastmessage string `json:"lastmessage"`
}

type Chat struct {
	Type string `json:"type"`
	Text string `json:"text"`
	From string `json:"from"`
}

type OldMessages struct {
	MessageID int
	Content   string
	Sender    string
	Receiver  string
	Date      string
}

type SendOldMessages struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type Post struct {
	Post_id  int       `json:"post_id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Username string    `json:"username"`
	Created  time.Time `json:"created"`
	Category string    `json:"category"`
}

type Posts struct {
	Type  string `json:"type"`
	New   bool   `json:"new"` // true if posts are not old content
	Posts []Post `json:"posts"`
}

type Comment struct {
	Comment_id int       `json:"comment_id"`
	Content    string    `json:"content"`
	Username   string    `json:"username"`
	Created    time.Time `json:"created"`
}

type Comments struct {
	Type     string    `json:"type"`
	New      bool      `json:"new"` // true if comments are not old content
	Post_id  int       `json:"post_id"`
	Comments []Comment `json:"comments"`
}

type LoginAnswer struct {
	Type     string `json:"type"`
	Username string `json:"username"` // retrieve username after login
	Txt      string `json:"txt"`      // message to display to client
}
