package rtfServer

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// main function
func LaunchServer() {
	InitiateDatabase()

	connectedUsers = make(map[string]*websocket.Conn)
	connectedUsersData = map[string]*conUsrData{}

	css := http.FileServer(http.Dir("./static/css"))
	http.Handle("/static/css/", http.StripPrefix("/static/css", css))
	js := http.FileServer(http.Dir("./static/js"))
	http.Handle("/static/js/", http.StripPrefix("/static/js", js))

	// handler for websocket initiation
	http.HandleFunc("/login", Login)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "home.html")
	})

	fmt.Println("Lisening on port : " + PORT)

	http.ListenAndServe(":"+PORT, nil)
}
