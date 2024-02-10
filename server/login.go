package rtfServer

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// login handler. used to upgrade connection and initiate websocket
// receives messages of type login or register
func Login(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade failed: ", err)
		return
	}
	defer conn.Close()

	for {
		var msgFromClient Message
		var respToClient LoginAnswer

		// read input from client
		err := conn.ReadJSON(&msgFromClient)

		// unkn own user disconnected
		if err != nil {
			return
		}

		var user_id int
		switch msgFromClient.Type {
		case "login":
			respToClient, user_id = CheckLogs(msgFromClient)
		case "register":
			respToClient = Register(msgFromClient, respToClient.Username, conn)
		default:
			fmt.Println("User sending invalid type message : " + msgFromClient.Type)
			respToClient.Type = "disconnected"
			respToClient.Txt = "Invalid type message."
			conn.WriteJSON(respToClient)
			return
		}

		conn.WriteJSON(respToClient)

		// login successful
		if respToClient.Type == "loginsuccess" {
			var username = respToClient.Username
			connectedUsers[username] = conn
			ShowUsers(username)
			HandleUser(conn, username, user_id)
			return
		}
	}
}

// checks email and password sent by client
func CheckLogs(msgFromClient Message) (LoginAnswer, int) {
	var storedPassword, username string
	var respToClient = LoginAnswer{Type: "loginfail"}

	userid := 0
	err := db.QueryRow("SELECT password, username, user_id FROM user WHERE email = ?", msgFromClient.Email).Scan(&storedPassword, &username, &userid)

	if err != nil {
		respToClient.Txt = "Invalid email"
	} else {
		if storedPassword == msgFromClient.Password {
			if _, logged := connectedUsers[username]; logged {
				respToClient.Txt = "User already logged in"
			} else {
				respToClient.Type = "loginsuccess"
				respToClient.Username = username
			}
		} else {
			respToClient.Txt = "Invalid password"
		}
	}
	return respToClient, userid
}

// register a new user
func Register(msgFromClient Message, username string, conn *websocket.Conn) LoginAnswer {
	fmt.Println(msgFromClient)
	var respToClient = LoginAnswer{Type: "registerfail"}
	count := 0

	err := db.QueryRow("SELECT COUNT(*) FROM user WHERE email = ?", msgFromClient.Email).Scan(&count)
	if err != nil {
		fmt.Println("Error in register (1)")
		respToClient.Txt = "Error from server"
		return respToClient
	}
	if count != 0 {
		respToClient.Txt = "Email already used"
		return respToClient
	}

	err = db.QueryRow("SELECT COUNT(*) FROM user WHERE username = ?", msgFromClient.Username).Scan(&count)
	if err != nil {
		fmt.Println("Error in register (2)")
		respToClient.Txt = "Error from server"
		return respToClient
	}
	if count != 0 {
		respToClient.Txt = "Username already used"
		return respToClient
	}

	_, err = db.Exec("INSERT INTO user (email, username, password) VALUES(?,?,?)", msgFromClient.Email, msgFromClient.Username, msgFromClient.Password)
	if err != nil {
		fmt.Println("Error in register (3)")
		respToClient.Txt = "Error from server"
		return respToClient
	}

	respToClient.Type = "registersuccess"
	ShowUsers(username)

	return respToClient
}
