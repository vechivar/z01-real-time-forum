// Store messages in the client
export let message = [];
export let messagesByUser = {};
export let user_selected = null;
let alertuser = {};

import { username } from '../main.js';
import { socket } from './socket.js';

let lastRefresh = 0, spamDelay = 2000;

// Function to handle old messages received from the server
export function oldmessagesofserv(msg) {
    // Check if the message data is null
    if (msg.data == null) {
        return;
    }

    let i = 0;
    const messages = document.getElementById('messages');
    const firstMess = messages.firstChild;

    while (i < msg.data.length) {
        let nouveauMessage = {
            content: msg.data[i].Content,
            date: msg.data[i].Date,
            by: msg.data[i].Sender,
            id: msg.data[i].MessageID
        };
        messagesByUser[user_selected].unshift(nouveauMessage);
        messages.prepend(buildMessageDiv(nouveauMessage));

        i++;
    }
    if (firstMess) {
        firstMess.scrollIntoView(true);
    } else {
        messages.lastChild.scrollIntoView(false);
    }
}

// Function to handle receiving chat messages
export function receiveChatMsg(msg) {
    let sender = msg.from;

    let nouveauMessage = {
        content: msg.text,
        date: getTime(),
        by: sender,
        id: 999
    };

    // Check if the sender is not the currently selected user
    if (sender !== user_selected) {
        // console.log(sender, "+", user_selected);
        // console.log(messagesByUser);
        var Element = document.getElementById(sender);

        Element.innerText = sender + " " + "!!!";
        alertuser[sender] = [];

        if (messagesByUser[sender]) {
            messagesByUser[sender].push(nouveauMessage);
        }
    } else {
        // If sender is the currently selected user
        if (!messagesByUser[sender]) {
            messagesByUser[sender] = [];
            getmessageofserv(sender, 0);
        }

        messagesByUser[sender].push(nouveauMessage);
        const div = buildMessageDiv(nouveauMessage);
        document.getElementById('messages').append(div);
        div.scrollIntoView(false);
    }
}

// Function to display connected/disconnected users
export function showUsers(msg) {
    let messagesDiv = document.getElementById("messages");

    const contactList = document.getElementById('contact-list');
    contactList.innerHTML = "";

    msg.data.forEach(function (user) {
        if (username != user.username) {
            const contact = document.createElement("div");
            if (alertuser[user.username]) {
                contact.textContent = user.username + "  !!!";
            } else {
                contact.textContent = user.username;
            }
            contact.classList.add(user.connected ? 'connected' : 'disconnected');

            contact.setAttribute('id', user.username);
            contact.classList.add(user.username);
            contact.classList.add("contact");

            contact.addEventListener('click', function () {
                const chatName = document.getElementById('chat-name');
                chatName.innerHTML = "";
                const roomName = document.createElement('div');
                roomName.innerText = "User : " + user.username;
                chatName.append(roomName);

                const loadMoreButton = document.createElement('button');
                loadMoreButton.innerText = "Load more";
                chatName.append(loadMoreButton);
                loadMoreButton.addEventListener('click', () => {
                    if (!messagesByUser[user_selected] || !messagesByUser[user_selected][0]) {
                        return;
                    }
                    if (Date.now() - lastRefresh > spamDelay) {
                        lastRefresh = Date.now();
                        getmessageofserv(user_selected, messagesByUser[user_selected][0].id);
                    }
                });

                user_selected = user.username;
                if (!messagesByUser[user_selected]) {
                    messagesByUser[user_selected] = [];
                    getmessageofserv(user_selected);
                }

                showMessages(user_selected);

                var Element = document.getElementById(user_selected);

                Element.innerText = user_selected;
                delete alertuser[user_selected];

                messagesDiv.scrollTop = messagesDiv.scrollHeight;

            });

            contactList.append(contact);
        }
    });
}

// Function to display message content in chat
export function showMessages(username) {
    let messagesDiv = document.getElementById("messages");
    messagesDiv.innerHTML = '';

    if (messagesByUser.hasOwnProperty(username)) {
        let messages = messagesByUser[username];

        for (let i = 0; i < messages.length; ++i) {

            messagesDiv.appendChild(buildMessageDiv(messages[i]));
        }
    } else {
        // console.log("Aucun message");
    }
}

// Function to build message div
function buildMessageDiv(msgData) {
    let div = document.createElement('div');
    div.classList = 'message';
    div.innerText = msgData.by + " (" + msgData.date + "): " + msgData.content;
    return div;
}

// Function to get old messages from the server
export function getmessageofserv(osuser, lastmess) {
    var getoldmessage = JSON.stringify({ type: "getoldmessage", to: osuser, time: getTime(), lastmessage_id: lastmess });
    socket.send(getoldmessage);
}

// Function to get current time
function getTime() {
    // Get current date and time
    const now = new Date();

    // Get hours and minutes
    const hours = now.getHours();
    const minutes = now.getMinutes();

    const formattedMinutes = minutes < 10 ? '0' + minutes : minutes;
    var time = (`${hours}:${formattedMinutes}`);
    return time;
}

// Function to build the chat interface
export function buildChat() {
    document.getElementById('chat').innerHTML = '' +
        '<div id="contact-list"></div>' +
        '<div id="private-chat">' +
        '<div id="chat-name"></div>' +
        '<div id="messages"></div>' +
        '<div id="msg-input-container">' +
        '<textarea id="msg-input"></textarea>' +
        '<button id="send-button">Send</button>' +
        '</div>';

    document.addEventListener('keyup', (e) => {
        const msgInput = document.getElementById('msg-input');
        if (e.key == 'Enter' && document.activeElement === msgInput) {
            sendMessageFunc(true);
        }
    });

    document.getElementById("send-button").addEventListener('click', () => {
        sendMessageFunc(false);
    });

    function sendMessageFunc(removeLast) {
        let contenu = document.getElementById("msg-input").value;

        if (removeLast) {
            contenu = contenu.substring(0, contenu.length - 1);
        }

        document.getElementById("msg-input").value = "";

        if (user_selected == undefined) {
            return
        }

        var sendmessage = JSON.stringify({ type: "chat", text: contenu, to: user_selected });
        let sender = user_selected;

        if (!messagesByUser[sender]) {
            messagesByUser[sender] = [];
            getmessageofserv(sender, 0);

        }

        let nouveauMessage = {
            content: contenu,
            date: getTime(),
            by: username,
            id: 777
        };
        messagesByUser[sender].push(nouveauMessage);

        const div = buildMessageDiv(nouveauMessage);
        document.getElementById('messages').append(div);
        div.scrollIntoView(false);

        socket.send(sendmessage);
    }

    const chat = document.getElementById('messages');

    chat.addEventListener('scroll', (e) => {
        if (!messagesByUser[user_selected] || !messagesByUser[user_selected][0]) {
            return;
        }

        if (chat.scrollTop === 0 && Date.now() - lastRefresh > spamDelay) {
            lastRefresh = Date.now();
            getmessageofserv(user_selected, messagesByUser[user_selected][0].id);
        }
    });
}
