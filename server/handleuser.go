package rtfServer

import (
	"fmt"

	"github.com/gorilla/websocket"
)

// listen to messages sent from connected users web sockets.
// sends messages to appropriate functions according to message type
func HandleUser(conn *websocket.Conn, username string, user_id int) {
	var msg Message

	connectedUsersData[username] = &conUsrData{0}
	connectedUsers[username] = conn
	ShowUsers(username)

	// send 10 most recent posts
	SendPosts(10, username, 0)

	for {
		err := conn.ReadJSON(&msg)
		if err != nil {
			// fmt.Println(username + " disconnected")
			delete(connectedUsers, username)
			delete(connectedUsersData, username)

			ShowUsers(username) // refresh users

			return
		}

		switch msg.Type {
		case "post":
			ReceivePost(msg, username, user_id)
		case "comment":
			ReceiveComment(msg, username, user_id)
		case "chat":
			SendMessageToChat(msg.Text, msg.To, username, msg.Time, conn)
		case "getoldmessage":
			GetOldMessage(msg.To, username, msg.Lastmessage_id, conn)
		case "request-comments":
			connectedUsersData[username].lastPostVisited = msg.Post_id
			SendComments(msg.Post_id, username, 10, msg.Last_id)
		case "request-posts":
			SendPosts(10, username, msg.Last_id)
		default:
			fmt.Println("something went wrong. unknown msg type from client : " + msg.Type)
		}
	}
}
