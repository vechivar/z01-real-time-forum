# z01-real-time-forum

This 2 people project is an introduction to web sockets. The goal is to build a single page site where identified users can communicate with each others with a general forum and a private-chat system.

## How to use

`go run .` to launch server, then visit http://localhost:8001

At first launch, a database will be created and populated with a few users and posts.

Create a new account and login to access the site. Several pages can be opened with different users to test the project.

## Features

- Login/register buttons to create a new account or login as an already existing user. Extra information during register are required by subject, but not currently stored anywhere.
- Forum part :
    - Post display with title, author, date and content.
    - Button to create a new post.
    - Clic on a post to display its comments and create a new comment.
    - Real time warning that a new post or comment was created on the page, and a button to display it.
    - More posts can be loaded when scrolling to the bottom of the section (with a delay to avoid spamming server)
- Chat part :
    - Display connected / disconnected users, sorted by last message exchange date, and alphabetical order if no message exchanged.
    - Clic on a user to open private chat window.
    - Get your 10 last messages when opening a chat window. Load more messages with "load more" button or by scrolling to top of the section. (with a delay to avoid spamming server)
    - Real time message exchange. New message by users will be signaled in users section if chat window with this user is not currently displayed.

## Project status

This project has been validated by Zone01 Rouen Community. It's quite rudimentary and rough, but most of the features that could have been implemented to improve it will be developed in a future project called social-network (not added to this github yet). Therefore, we chose to handle this project as a draft for social-network, and a training in websocket management. We felt like this goal has been reached.