import {socket} from './socket.js'

let loaded = false,                 // page already received first posts
    loadedComments = false,         // page already received comments
    waitingPosts = Array(), // new posts waiting to be displayed
    displayedPost, // post currently displayed in page with comments
    waitingComments = Array(),  // new comments waiting to be displayed
    lastPostId = 0, // oldest post id, passed to server when requiring older posts
    lastCommentId = 0,  // oldest comment id, passed to server when requiring older comments
    lastRequest = 0,    // time of last request to server, used to avoid spamming
    spamDelay = 2000    // time between requests to avoid spamming

// initial build post part of the page
export function buildPostPage() {
    document.getElementById('posts').innerHTML = 
    '<div id="posts-container">' + 
        '<div id ="create-post"></div>'+
        '<div id="display-posts"></div>'+
    '</div>' + 
    '<div id="view-post"></div>'
    createNewPostDiv()
    const posts = document.getElementById('posts')

    posts.addEventListener('scroll', (e) => {
        // load more elements when scrolling
        if (posts.scrollTop >= (posts.scrollHeight - posts.offsetHeight) && Date.now() - lastRequest > spamDelay) {
            if (displayedPost == undefined) {
                socket.send(JSON.stringify({type:'request-posts', last_id:lastPostId}))
            } else {
                socket.send(JSON.stringify({type:'request-comments', last_id:lastCommentId}))
            }
            lastRequest = Date.now()
        }
    })
}

// triggered when receiving messages of type post
// adds posts on top of the page, updates buttons to display new posts
export function receivePostMsg(msg) {
    console.log(msg)

    const posts = document.getElementById('posts')
    const displayPosts = document.getElementById('display-posts')
    const postsContainer = document.getElementById('posts-container')
    let insertedLast 

    msg.posts.forEach(e => {
        const div = createPostDiv(e)
        div.style.cursor = 'pointer'

        if (e.post_id < lastPostId || lastPostId == 0) {
            lastPostId = e.post_id
        }

        // ugly stuff to make sure newer posts get on top of the page
        if (msg.new) {
            // new message
            displayPosts.after(div)
        } else {
            // older messages
            if (insertedLast == undefined) {
                postsContainer.lastChild.after(div)
            } else {
                postsContainer.insertBefore(div, insertedLast)
            }
            insertedLast = div
        }

        // put posts in waiting posts
        if (loaded && msg.new) {
            div.style.display = 'none'
            waitingPosts.push(div)
        }
        div.addEventListener('click', () => {
            viewPost(e)
        })
    })

    
    if (loaded && msg.new) {
        displayPosts.innerHTML = ""
        const button = document.createElement('button')
        displayPosts.append(button)
        button.innerText = 'Load new posts'
        button.addEventListener('click', () => {
            displayWaitingPosts()
        })
    } else {
        displayPosts.style.display = 'flex'
        displayPosts.innerHTML = "<div>No new post.</div>"
    }

    loaded = true
}

// display posts waiting to be loaded
function displayWaitingPosts() {
    waitingPosts.forEach(e => {
        e.style.display = 'block'
    })
    document.getElementById('display-posts').innerHTML = "<div>No new post.</div>"
    waitingPosts = Array()
}

// loads page for specific post and comments
// requests comments to server and creates buttons
function viewPost(postData) {
    socket.send(JSON.stringify({type:"request-comments", post_id:postData.post_id, last_id:0}))
    loadedComments = false
    lastCommentId = 0

    displayedPost = postData

    document.getElementById('posts-container').style.display = 'none'

    const vpDiv = document.getElementById('view-post')
    vpDiv.innerHTML = ""
    vpDiv.style.display = 'flex'

    const postDiv = createPostDiv(postData)
    const commentsDiv = document.createElement('div')
    commentsDiv.id = 'comments-container'

    const button = document.createElement('button')
    button.innerText = "Go back to forum"
    button.style.width = "50%"
    button.style.margin = "5px"
    button.style.alignSelf = "center"

    button.addEventListener('click', () => {
        vpDiv.innerHTML = ""
        vpDiv.style.display = 'none'
        document.getElementById('posts-container').style.display = 'flex'
        displayWaitingPosts()
        displayedPost = undefined
    })

    const loadNewComments = document.createElement('div')
    loadNewComments.id = 'display-comments'

    vpDiv.append(button)
    vpDiv.append(postDiv)
    vpDiv.append(createNewCommentDiv())
    vpDiv.append(loadNewComments)
    vpDiv.append(commentsDiv)
}

// creates container and content for new post fields
function createNewPostDiv() {
    const div = document.getElementById('create-post')
    let shown = false

    const showButton = document.createElement('button')
    showButton.innerText = "Create new post"
    showButton.style.margin = "8px"
    showButton.addEventListener('click',() => {
        if (shown) {
            showButton.innerText = "Create new post"
            container.style.display = 'none'
        } else {
            showButton.innerText = "Close"
            container.style.display = 'flex'
        }
        shown = !shown
    })

    const container = document.createElement('div')
    container.classList = 'create-post-container'

    div.append(showButton)
    div.append(container)
    container.style.display = 'none'

    const title = document.createElement('input')
    title.type = "text"
    title.classList = 'title-input'
    title.placeholder = "Title"
    container.append(title)

    const content = document.createElement('textarea')
    content.classList = 'content-input'
    container.append(content)

    const button = document.createElement('button')
    button.innerText = "Send"
    button.style.width = "50%"
    button.style.margin = "5px"
    button.style.alignSelf = "center"
    button.addEventListener('click', () => {
        if (content.value == '' || title.value == '') {
            return
        }
        socket.send(JSON.stringify({type:"post", title:title.value, content:content.value}))
        title.value = ""
        content.value = ""
    })
    container.append(button)
}

// creates container and content for new comment fields
function createNewCommentDiv() {
    const div = document.createElement('div')
    div.id = 'create-comment'
    let shown = false

    const showButton = document.createElement('button')
    showButton.innerText = "Create new comment"
    showButton.style.margin = "8px"
    showButton.addEventListener('click',() => {
        if (shown) {
            showButton.innerText = "Create new comment"
            container.style.display = 'none'
        } else {
            showButton.innerText = "Close"
            container.style.display = 'flex'
        }
        shown = !shown
    })

    const container = document.createElement('div')
    container.classList = 'create-comment-container'

    div.append(showButton)
    div.append(container)
    container.style.display = 'none'

    const content = document.createElement('textarea')
    content.classList = 'content-input'
    container.append(content)

    const button = document.createElement('button')
    button.innerText = "Send"
    button.style.width = "50%"
    button.style.margin = "5px"
    button.style.alignSelf = "center"
    button.addEventListener('click', () => {
        if (content.value == '') {
            return
        }
        socket.send(JSON.stringify({type:"comment", post_id:displayedPost.post_id, content:content.value}))
        content.value = ""
    })
    container.append(button)

    return div
}

// triggered when receiving msg of type comment
// inserts all comments received in page 
export function receiveCommentMsg(msg) {
    const commentsDiv = document.getElementById('comments-container')
    if (commentsDiv == undefined || msg.post_id != displayedPost.post_id) {
        // user changed page, comments are now useless
        return
    }

    const loadNewComments = document.getElementById('display-comments')

    if (loadedComments && msg.new) {
        loadNewComments.innerHTML = ''
        const button = document.createElement('button')
        button.innerText = "Load new comments."

        button.addEventListener('click', () => {
            loadNewComments.innerHTML = '<div>No new comments.</div>'
            waitingComments.forEach(e => {
                e.style.display = 'block'
            })
            waitingComments = Array()
        })

        loadNewComments.append(button)
    } else {
        loadNewComments.innerHTML = '<div>No new comments.</div>'
    }

    let lastInserted

    msg.comments.forEach(e => {
        if (lastCommentId == 0 || e.comment_id < lastCommentId) {
            lastCommentId = e.comment_id
        }
        const firstComment = commentsDiv.firstChild
        const div = createCommentDiv(e)
        if (firstComment == undefined) {
            commentsDiv.append(div)
        } else {
            if (msg.new) {
                commentsDiv.insertBefore(div, firstComment)
            } else {
                if (lastInserted == undefined) {
                    commentsDiv.lastChild.after(div)
                } else {
                    commentsDiv.insertBefore(div, lastInserted)
                }
                lastInserted = div
            }
        }
        if (loadedComments && msg.new) {
            waitingComments.push(div)
            div.style.display = 'none'
        }
    })

    loadedComments = true
}

// creates div for post datas
function createPostDiv(postData) {
    const post = createDiv('post', "")
    const title = createDiv('title', "Title : " + postData.title)
    const username = createDiv('username', "Author : " + postData.username)
    // const category = createDiv('category', "Category : " + postData.category)

    const time = new Date(postData.created)
    const date = createDiv('date', "Created : " + formatDate(time))

    const content = createDiv('content', postData.content)
    post.append(title)
    post.append(username)
    // post.append(category)
    post.append(date)
    post.append(content)

    return post
}

// creates dive for comment datas
function createCommentDiv(commentData) {
    const comment = createDiv('comment', "")
    const username = createDiv('username', "Author : " + commentData.username)
    const date = createDiv('date', "Created : " + formatDate(new Date(commentData.created)))
    const content = createDiv('content', commentData.content)
    comment.append(username)
    comment.append(date)
    comment.append(content)

    return comment
}

// creates a div, intiated with class and innertext
function createDiv(classname, txt) {
    const div = document.createElement('div')
    div.classList = classname
    div.innerText = txt
    return div
}

// proper date formating
function formatDate(date) {
    let h = date.getHours(), m = date.getMinutes()

    return date.toDateString() + " ("  + (h<10?'0'+h:h) + ':'  + (m<10?'0'+m:m) + ')'
}
