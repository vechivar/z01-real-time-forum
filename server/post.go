package rtfServer

import (
	"database/sql"
	"fmt"
	"time"
)

// sends n posts from database to a connected user.
// last_id is used to pass last post_id from user and get older posts
// if last_id == 0 : send n most recent posts
func SendPosts(n int, connectedUsername string, last_id int) {
	var rows *sql.Rows
	var err error
	if last_id == 0 {
		rows, err = db.Query("SELECT * FROM post ORDER BY created DESC limit ?", n)
	} else {
		rows, err = db.Query("SELECT * FROM post WHERE post_id < ? ORDER BY created DESC limit ?", last_id, n)
	}

	if err != nil {
		fmt.Print("Error getting posts from database : query failed\n)")
		fmt.Println(err)
		return
	}

	var posts []Post

	// post(post_id, created, title, content, user_id, category_id)
	for rows.Next() {
		var (
			post_id     int
			created     time.Time
			title       string
			content     string
			user_id     int
			username    string
			category_id int
			category    string
		)

		rows.Scan(&post_id, &created, &title, &content, &user_id, &category_id)

		username, err = GetUsernameFromId(user_id)
		if err != nil {
			fmt.Printf("Error getting posts from database : unkwown user_id %d\n", user_id)
			fmt.Println(err)
			return
		}

		category, err = GetCategoryFromId(category_id)
		if err != nil {
			fmt.Printf("Error getting posts from database : unkwown category_id %d\n", category_id)
			fmt.Println(err)
			return
		}

		posts = append([]Post{{post_id, title, content, username, created, category}}, posts...)
	}

	if len(posts) == 0 {
		return
	}

	err = connectedUsers[connectedUsername].WriteJSON(Posts{"post", last_id == 0, posts})

	if err != nil {
		if _, logged := connectedUsers[connectedUsername]; logged {
			fmt.Printf("Error sending posts sending posts to user %s\n", connectedUsername)
			fmt.Println(err)
		}
	}
}

// sends n comments from database to a connected user.
// last_id is used to pass last post_id from user and get older posts
// if last_id == 0 : send n most recent posts
func SendComments(post_id int, connectedUsername string, n int, last_id int) {
	var rows *sql.Rows
	var err error

	if last_id == 0 {
		rows, err = db.Query("SELECT comment_id, created, content, user_id FROM comment WHERE post_id = ? ORDER BY created DESC limit ?", post_id, n)
	} else {
		rows, err = db.Query("SELECT comment_id, created, content, user_id FROM comment WHERE post_id = ? AND comment_id < ? ORDER BY created DESC limit ?", post_id, last_id, n)
	}

	if err != nil {
		fmt.Print("Error getting comments from database : query failed\n)")
		fmt.Println(err)
		return
	}

	var comments []Comment

	// comment(comment_id, created, content, user_id, post_id)
	for rows.Next() {
		var (
			comment_id int
			created    time.Time
			content    string
			user_id    int
			username   string
		)

		rows.Scan(&comment_id, &created, &content, &user_id)

		username, err = GetUsernameFromId(user_id)
		if err != nil {
			fmt.Printf("Error getting posts from database : unkwown user_id %d\n", user_id)
			fmt.Println(err)
			return
		}

		comments = append([]Comment{{comment_id, content, username, created}}, comments...)
	}

	if len(comments) == 0 {
		return
	}

	err = connectedUsers[connectedUsername].WriteJSON(Comments{"comment", last_id == 0, post_id, comments})

	if err != nil {
		if _, logged := connectedUsers[connectedUsername]; logged {
			fmt.Printf("Error sending posts sending posts to user %s\n", connectedUsername)
			fmt.Println(err)
		}
	}
}

// handles new post data from client
// inserts into database and sends notification to other connected users
func ReceivePost(msg Message, username string, user_id int) {
	res, err := db.Exec("INSERT INTO post(title, content, user_id, category_id) VALUES(?,?,?,1)", msg.Title, msg.Content, user_id)
	if err != nil {
		fmt.Printf("Something went wrong inserting new post from user %s (1)\n", username)
		fmt.Println(err)
		return
	}

	post_id, err := res.LastInsertId()
	if err != nil {
		fmt.Printf("Something went wrong inserting new post from user %s (2)\n", username)
		fmt.Println(err)
		return
	}

	var created time.Time
	err = db.QueryRow("SELECT created FROM post WHERE post_id = ?", post_id).Scan(&created)
	if err != nil {
		fmt.Printf("Something went wrong inserting new post from user %s (3)\n", username)
		fmt.Println(err)
		return
	}

	var post []Post
	post = append(post, Post{int(post_id), msg.Title, msg.Content, username, created, "none"})

	for _, conn := range connectedUsers {
		conn.WriteJSON(Posts{"post", true, post})
	}
}

// handles new comment data from client
// inserts into database and sends notification to other connected users
func ReceiveComment(msg Message, username string, user_id int) {
	res, err := db.Exec("INSERT INTO comment(content, user_id, post_id) VALUES(?,?,?)", msg.Content, user_id, msg.Post_id)
	if err != nil {
		fmt.Printf("Something went wrong inserting new comment from user %s on post %d (1)\n", username, msg.Post_id)
		fmt.Println(err)
		return
	}

	comment_id, err := res.LastInsertId()
	if err != nil {
		fmt.Printf("Something went wrong inserting new comment from user %s on post %d (2)\n", username, msg.Post_id)
		fmt.Println(err)
		return
	}

	var created time.Time
	err = db.QueryRow("SELECT created FROM comment WHERE comment_id = ?", comment_id).Scan(&created)
	if err != nil {
		fmt.Printf("Something went wrong inserting new comment from user %s on post %d (3)\n", username, msg.Post_id)
		fmt.Println(err)
		return
	}

	var comment []Comment
	comment = append(comment, Comment{int(comment_id), msg.Content, username, created})

	for usr, conn := range connectedUsers {
		if connectedUsersData[usr].lastPostVisited == msg.Post_id {
			conn.WriteJSON(Comments{"comment", true, msg.Post_id, comment})
		}
	}
}
