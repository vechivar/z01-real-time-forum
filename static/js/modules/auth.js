import { socket } from "./socket.js"

// builds divs for login header and disconnect button
export function buildLogger() {
    const header = document.getElementById('header')
    header.innerHTML = ''

    const logContainer = document.createElement('div')
    header.append(logContainer)
    logContainer.classList = 'log-container'

    const email = document.createElement('input')
    email.type = "text"
    email.id = 'email-log'
    email.placeholder = 'email'

    const password = document.createElement('input')
    password.type = "password"
    password.id = 'password-log'
    password.placeholder = "password"

    const button = document.createElement('button')
    button.innerText = 'Login'
    button.addEventListener('click', () => {
        const txt = JSON.stringify({type:"login",email:document.getElementById("email-log").value, password:document.getElementById("password-log").value})
        socket.send(txt)
    })

    const registerButton = document.createElement('button')
    registerButton.innerText = "Register"
    registerButton.addEventListener('click', () => {
        buildRegister()
    })

    const logMsg = document.createElement('div')
    logMsg.id = 'logmsg'

    logContainer.append(email)
    logContainer.append(password)
    logContainer.append(button)
    logContainer.append(registerButton)
    logContainer.append(logMsg)
}

// builds header after login
export function buildHeader(username) {
    const header = document.getElementById('header')
    header.innerHTML = ""

    const headerContainer = document.createElement('div')
    headerContainer.id = "header-container"

    header.append(headerContainer)

    const welcome = document.createElement('div')
    welcome.innerText = "Welcome " + username + "!"

    const disconnect = document.createElement('button')
    disconnect.innerText = "Disconnect"
    disconnect.id = "disconnect"

    headerContainer.append(welcome)
    headerContainer.append(disconnect)
}

// build register inputs after clicking on register button
function buildRegister() {
    const posts = document.getElementById('posts')
    posts.innerHTML = ""

    const container = document.createElement('div')
    container.style.display = "flex"
    container.style.flexDirection = "column"
    posts.append(container)

    buildRegisterInput("username", "text", container)
    buildRegisterInput("first-name", "text", container)
    buildRegisterInput("last-name", "text", container)
    buildRegisterInput("age", "number", container)
    buildRegisterInput("gender", "text", container)
    buildRegisterInput("email", "email", container)
    buildRegisterInput("password", "password", container)
    buildRegisterInput("confirm", "password", container)

    const button = document.createElement('button')
    button.style.width = "30%";
    button.style.alignSelf = "center"
    button.innerText = "Submit"
    container.append(button)

    const msg = document.createElement('div')
    msg.id = "regLog"
    container.append(msg)

    button.addEventListener('click', () => {
        // send register datas to server
        msg.innerText = ""
        let flag = true
        if (document.getElementById('confirm').value != document.getElementById('password').value) {
            msg.innerText = "Error : password and confirm are different.\n"
            flag = false
        }
        const fields = ['username', 'first-name', 'last-name', 'age', 'gender', 'email', 'password', 'confirm']
        fields.forEach(e => {
            if (document.getElementById(e).value == "") {
                flag = false
                msg.innerText += "Error : field " + e + " is empty.\n"
            }
        })
        if (!flag) {
            return
        }
        const msgToServer = {type:"register"}
        for (let i = 0; i < fields.length - 1; i++) {
            msgToServer[fields[i]] = document.getElementById(fields[i]).value
        }
        socket.send(JSON.stringify(msgToServer))
    })

}

// subfunction for buildRegister()
function buildRegisterInput(name, type, parent) {
    const div = document.createElement('div')
    div.style.alignSelf = "center"
    div.style.margin = "10px"
    const input = document.createElement('input')
    div.append(input)
    input.type = type
    input.id = name
    input.placeholder = name
    parent.append(div)
}