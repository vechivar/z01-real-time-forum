package rtfServer

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// user(user_id, email, username, password)
// category(category_id, name)
// post(post_id, created, title, content, user_id, category_id)
// comment(comment_id, created, content, user_id, post_id)

func InitiateDatabase() {
	var err error

	if _, err := os.Stat(dbPath); err == nil {
		db, _ = sql.Open("sqlite3", dbPath)
		return
	}

	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	CreateTables()
	AddContent()
}

func CreateTables() {
	_, err := db.Exec(`
		CREATE TABLE user (
			user_id INTEGER PRIMARY KEY AUTOINCREMENT,
			email    TEXT,
			username TEXT,
			password TEXT
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`
	CREATE TABLE message (
		message_id INTEGER PRIMARY KEY AUTOINCREMENT,
	    content TEXT,
		date TEXT,
		sender TEXT,
		receiver TEXT
	);
	
`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE category (
			category_id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE post (
			post_id INTEGER PRIMARY KEY AUTOINCREMENT,
			created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			title TEXT,
			content TEXT,
			user_id INTEGER REFERENCES user(user_id),
			category_id REFERENCES category(category_id)
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE comment (
			comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
			created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			content TEXT,
			user_id REFERENCES user (user_id),
			post_id REFERENCES post (post_id)
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`
	CREATE TABLE lastmessage (
		lastmessage_id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_1 TEXT,
		user_2 TEXT,
		lasttime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP

	)
`)
	if err != nil {
		log.Fatal(err)
	}
}

func AddContent() {
	_, err := db.Exec(`
	INSERT INTO user (email, username, password) VALUES
	('jean@jean.com', 'jean', 'jeanjean'),
	('michel@michel.com', 'michel', 'michelmichel'),
	('alice@alice.com', 'alice', 'alicealice'),
	('emma@emma.com', 'emma', 'emmaemma'),
	('pierre@pierre.com', 'pierre', 'pierrepierre'),
	('lucie@lucie.com', 'lucie', 'lucielucie'),
	('olivier@olivier.com', 'olivier', 'olivierolivier'),
	('claire@claire.com', 'claire', 'claireclaire'),
	('antoine@antoine.com', 'antoine', 'antoineantoine'),
	('sophie@sophie.com', 'sophie', 'sophiesophie')
	`)
	if err != nil {
		fmt.Println("Error adding users")
		log.Fatal(err)
	}

	_, err = db.Exec(`
	INSERT INTO category (name) VALUES
	('none'),
	('random')
	`)
	if err != nil {
		fmt.Println("Error adding categories")
		log.Fatal(err)
	}

	_, err = db.Exec(`
	INSERT INTO post (title, content, user_id, category_id) VALUES
	('Post by jean (category none)', 'incroyable contenu de jean (none)', 1,1)`)

	if err != nil {
		fmt.Println("Error adding posts")
		log.Fatal(err)
	}

	_, err = db.Exec(`
	INSERT INTO post (title, content, user_id, category_id) VALUES
	('Post by jean (category random)', 'incroyable contenu de jean (random)', 1,2)`)

	if err != nil {
		fmt.Println("Error adding posts")
		log.Fatal(err)
	}

	_, err = db.Exec(`
	INSERT INTO post (title, content, user_id, category_id) VALUES
	('Post by alice (category none)', 'incroyable contenu de alice (none)', 3,1)`)

	if err != nil {
		fmt.Println("Error adding posts")
		log.Fatal(err)
	}

	// comment(comment_id, created, content, user_id, post_id)
	_, err = db.Exec(`
	INSERT INTO comment (content, user_id, post_id) VALUES
	('Comment by alice to post by jean', 3,1)`)

	if err != nil {
		fmt.Println("Error adding comments")
		log.Fatal(err)
	}

	_, err = db.Exec(`
    INSERT INTO message (content, date, sender, receiver) VALUES
	('Bonjour', '15:25', 'jean', 'alice'),
	('Bonjour', '15:28', 'alice', 'jean'),
	('Bonjour', '15:29', 'jean', 'alice'),
	('Salut', '15:35', 'alice', 'jean'),
	('Salut', '15:40', 'jean', 'alice'),
	('Comment ça va?', '15:45', 'alice', 'jean'),
	('Ça va bien, merci!', '15:50', 'jean', 'alice'),
	('Et toi?', '15:55', 'alice', 'jean'),
	('Moi aussi ça va bien', '16:00', 'jean', 'alice'),
	('Qu''as-tu fait aujourd''hui?', '16:05', 'alice', 'jean'),
	('Rien de spécial, et toi?', '16:10', 'jean', 'alice'),
	('J''ai travaillé un peu et ensuite j''ai regardé un film', '16:15', 'alice', 'jean'),
	('C''était un bon film?', '16:20', 'jean', 'alice'),
	('Oui, j''ai beaucoup aimé', '16:25', 'alice', 'jean'),
	('Quel genre de film c''était?', '16:30', 'jean', 'alice'),
	('C''était une comédie romantique', '16:35', 'alice', 'jean'),
	('J''aime aussi les comédies romantiques', '16:40', 'jean', 'alice'),
	('On devrait en regarder un ensemble un jour', '16:45', 'alice', 'jean'),
	('Ça serait génial!', '16:50', 'jean', 'alice'),
	('Quel est ton acteur préféré?', '16:55', 'alice', 'jean'),
	('J''adore Ryan Gosling', '17:00', 'jean', 'alice'),
	('Il est vraiment talentueux', '17:05', 'alice', 'jean'),
	('Et toi?', '17:10', 'jean', 'alice'),
	('Je suis fan de Leonardo DiCaprio', '17:15', 'alice', 'jean'),
	('Il a remporté enfin un Oscar il y a quelques années', '17:20', 'jean', 'alice'),
	('Oui, enfin! Il le méritait depuis longtemps', '17:25', 'alice', 'jean'),
	('On devrait organiser une soirée cinéma chez moi', '17:30', 'jean', 'alice'),
	('Ça me semble être une excellente idée', '17:35', 'alice', 'jean'),
	('Quand es-tu disponible?', '17:40', 'jean', 'alice'),
	('Je suis libre ce week-end', '17:45', 'alice', 'jean'),
	('Parfait! On peut se retrouver samedi soir', '17:50', 'jean', 'alice'),
	('D''accord, samedi soir ça me convient', '17:55', 'alice', 'jean'),
	('J''ai hâte!', '18:00', 'jean', 'alice'),
	('Moi aussi, ça va être amusant', '18:05', 'alice', 'jean'),
	('As-tu des suggestions de films?', '18:10', 'jean', 'alice'),
	('Qu''en penses-tu de "La La Land"?', '18:15', 'alice', 'jean'),
	('C''est une excellente idée! J''adore ce film', '18:20', 'jean', 'alice'),
	('Parfait! On a trouvé notre film', '18:25', 'alice', 'jean'),
	('Tu veux qu''on apporte quelque chose?', '18:30', 'jean', 'alice'),
	('Peut-être des snacks et des boissons?', '18:35', 'alice', 'jean'),
	('D''accord, je m''occuperai des snacks', '18:40', 'jean', 'alice'),
	('Et moi des boissons!', '18:45', 'alice', 'jean'),
	('Super! Ça va être une soirée géniale', '18:50', 'jean', 'alice'),
	('Je suis d''accord, j''ai hâte d''être à samedi', '18:55', 'alice', 'jean'),
	('On peut inviter d''autres amis aussi', '19:00', 'jean', 'alice'),
	('Bonne idée! Je connais quelques personnes qui seraient intéressées', '19:05', 'alice', 'jean'),
	('Génial! Plus on est de fous, plus on rit', '19:10', 'jean', 'alice'),
	('Exactement! Ça va être une soirée mémorable', '19:15', 'alice', 'jean'),
	('J''ai hâte de voir tout le monde', '19:20', 'jean', 'alice'),
	('Ça va être une soirée cinéma inoubliable', '19:25', 'alice', 'jean'),
	('Tu as raison, on devrait faire ça plus souvent', '19:30', 'jean', 'alice'),
	('Absolument! On peut faire des soirées cinéma régulières', '19:35', 'alice', 'jean'),
	('On peut même choisir un thème différent à chaque fois', '19:40', 'jean', 'alice'),
	('C''est une excellente idée! On pourrait faire une nuit d''horreur', '19:45', 'alice', 'jean'),
	('Oui, ça pourrait être amusant! On pourrait regarder des classiques de l''horreur', '19:50', 'jean', 'alice')`)
	if err != nil {
		fmt.Println("Error adding users")
		log.Fatal(err)
	}

	_, err = db.Exec(`
INSERT INTO lastmessage (user_1, user_2) VALUES
('jean', 'alice')`)

	if err != nil {
		fmt.Println("Error adding comments")
		log.Fatal(err)
	}

}
