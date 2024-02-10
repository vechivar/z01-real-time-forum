import * as Auth from './modules/auth.js'
import * as Post from './modules/posts.js'
import * as Chat from './modules/chat.js'
import {socket} from './modules/socket.js'

export var logged = false, username

Auth.buildLogger()

socket.onmessage = (e) => {
    const msg = JSON.parse(e.data)

    // handles messages related to register and login
    if (!logged) {
        switch (msg.type) {
            case "loginsuccess":
                logged = true
                username = msg.username
                createPage()
            break
            case "loginfail":
                document.getElementById('logmsg').innerText = msg.txt
            break
            case "registersuccess":
                document.getElementById('logmsg').innerText = "Register successful. You can now log in."
                document.getElementById('posts').innerHTML = ""
            break
            case "registerfail":
                const log = document.getElementById('regLog')
                if (log) {
                    log.innerText = msg.txt
                }
            break
            default:
            console.log("Msg type not supported : " + msg.type)
            break
        }
        return
    }

    // sends messages to appropriate functions
    switch (msg.type) {
        case "users":
        Chat.showUsers(msg)
        break
        case "oldmessages":
        Chat.oldmessagesofserv(msg)
        break
        case "chat":
        Chat.receiveChatMsg(msg)
        break
        case "post":
        Post.receivePostMsg(msg)
        break
        case "comment":
        Post.receiveCommentMsg(msg)
        break
        default:
        console.log("Msg type not supported : " + msg.type)
        break
    }
}

function createPage() {
    Post.buildPostPage()
    Auth.buildHeader(username)
    Chat.buildChat()

    document.getElementById('disconnect').addEventListener('click', () => {
        location.reload()
    })
}