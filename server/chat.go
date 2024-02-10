package rtfServer

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

// GetOldMessage retrieves old messages from the database and sends them to the client.
func GetOldMessage(to, username string, lastmessage_id int, conn *websocket.Conn) {
	var rows *sql.Rows
	var err error

	// Check if lastmessage_id is 0, if yes, retrieve the latest 10 messages.
	if lastmessage_id == 0 {
		rows, err = db.Query(`
            SELECT message_id, content, date, sender, receiver
            FROM message
            WHERE (sender = ? AND receiver = ?)
            OR (receiver = ? AND sender = ?)
            ORDER BY message_id DESC
            LIMIT 10
        `, username, to, username, to)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// If lastmessage_id is provided, retrieve 10 messages with IDs less than lastmessage_id.
		rows, err = db.Query(`
            SELECT message_id, content, date, sender, receiver
            FROM message
            WHERE (sender = ? AND receiver = ? AND message_id < ?)
            OR (receiver = ? AND sender = ? AND message_id < ?)
            ORDER BY message_id DESC
            LIMIT 10
        `, username, to, lastmessage_id, username, to, lastmessage_id)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer rows.Close()

	var messages []OldMessages
	for rows.Next() {
		var message OldMessages
		// Scan each row into the OldMessages struct.
		if err := rows.Scan(&message.MessageID, &message.Content, &message.Date, &message.Sender, &message.Receiver); err != nil {
			log.Fatal(err)
		}
		messages = append(messages, message)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Create a SendOldMessages struct and send it to the client.
	message := SendOldMessages{
		Type: "oldmessages",
		Data: messages,
	}
	json, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
		return
	}
	jsonString := string(json)
	err = conn.WriteMessage(websocket.TextMessage, []byte(jsonString))
	if err != nil {
		fmt.Println(err)
	}
}

// ShowUsers sends the list of users and their status to the connected client.
func ShowUsers(username string) {
	// Iterate through connectedUsers to send user status information.
	for test, conn := range connectedUsers {
		// Get user status from the SQLite database.
		userswithstatus := GetUsersWithStatusFromSQLite(test)
		// Create a User struct and send it to the client.
		message := User{
			Type:          "users",
			Data:          userswithstatus,
			Userconnected: username,
		}
		json, err := json.Marshal(message)
		if err != nil {
			fmt.Println(err)
			return
		}
		jsonString := string(json)
		err = conn.WriteMessage(websocket.TextMessage, []byte(jsonString))
		if err != nil {
			fmt.Println(err)
		}
	}
}

// GetUsersWithStatusFromSQLite retrieves user status information from the SQLite database.
func GetUsersWithStatusFromSQLite(username string) []UserStatus {
	user1Value := username
	db, err := sql.Open("sqlite3", "./db/database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Retrieve user status from the lastmessage table.
	rows, err := db.Query(`
        SELECT user_2, lasttime
        FROM lastmessage
        WHERE user_1 = ?
        UNION
        SELECT user_1, lasttime
        FROM lastmessage
        WHERE user_2 = ?
        ORDER BY lasttime DESC
    `, user1Value, user1Value)
	if err != nil {
		log.Fatal(err)
	}

	var utilisateurs []UserStatus
	for rows.Next() {
		var user UserStatus
		// Scan each row into the UserStatus struct.
		if err := rows.Scan(&user.Username, &user.Lastmessage); err != nil {
			log.Fatal(err)
		}
		if conn, ok := connectedUsers[user.Username]; ok {
			user.Connected = conn != nil
		}
		utilisateurs = append(utilisateurs, user)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Retrieve usernames from the user table that are not present in the lastmessage table.
	rowsAutreTable, err := db.Query(`
        SELECT username
        FROM user
        WHERE username NOT IN (
            SELECT user_2
            FROM lastmessage
            WHERE user_1 = ?
            UNION
            SELECT user_1
            FROM lastmessage
            WHERE user_2 = ?
        )
        ORDER BY username ASC
    `, user1Value, user1Value)
	if err != nil {
		log.Fatal(err)
	}

	// Iterate through the rows and append UserStatus structs to the utilisateurs slice.
	for rowsAutreTable.Next() {
		var user UserStatus
		if err := rowsAutreTable.Scan(&user.Username); err != nil {
			log.Fatal(err)
		}
		if conn, ok := connectedUsers[user.Username]; ok {
			user.Connected = conn != nil
		}
		utilisateurs = append(utilisateurs, user)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return utilisateurs
}

// SendMessageToChat sends a chat message to the specified user and updates user status.
func SendMessageToChat(text, to, name, time string, conn *websocket.Conn) {
	// Create a Chat struct and send it to the client.
	message := Chat{
		Type: "chat",
		Text: text,
		From: name,
	}
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Save the message to the database.
	MessageToDB(text, to, name, time)
	// Update lastmessage table.
	LastMessageToDB(to, name)
	if a, ok := connectedUsers[to]; ok {
		// Send the chat message to the recipient if they are connected.
		err := a.WriteMessage(websocket.TextMessage, jsonMessage)
		if err != nil {
			fmt.Println(err)
		}
	}
	// Update user status for the sender and recipients.
	ShowUsers(name)
}

// MessageToDB saves the chat message to the database.
func MessageToDB(text, to, name, time string) {
	_, err := db.Exec(`
        INSERT INTO message (content, sender, receiver, date)
        VALUES (?, ?, ?, ?)
    `, text, name, to, time)
	if err != nil {
		log.Fatal(err)
	}
}

// LastMessageToDB updates the lastmessage table in the database.
func LastMessageToDB(user1, user2 string) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	var alreadyExists bool
	err = tx.QueryRow(`
        SELECT EXISTS (
            SELECT 1
            FROM lastmessage
            WHERE (user_1 = ? AND user_2 = ?) OR (user_1 = ? AND user_2 = ?)
        )
    `, user1, user2, user2, user1).Scan(&alreadyExists)
	if err != nil {
		log.Fatal(err)
	}

	if alreadyExists {
		// If the pair already exists in lastmessage, delete it.
		_, err := tx.Exec(`
            DELETE FROM lastmessage
            WHERE (user_1 = ? AND user_2 = ?) OR (user_1 = ? AND user_2 = ?)
        `, user1, user2, user2, user1)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Insert the new pair into lastmessage.
	_, err = tx.Exec(`
        INSERT INTO lastmessage (user_1, user_2)
        VALUES (?, ?)
    `, user1, user2)
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}
