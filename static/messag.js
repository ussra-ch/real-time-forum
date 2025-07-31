import { webSockert } from "./websocket.js"
export function mesaageDiv(user, userId, receiverId) {
    const body = document.querySelector('body')
    const div = document.createElement('div')
    div.id = 'messag'
    div.innerHTML = `
        <div class="head">
            <h3>${user}</h3>
        </div>
        <div class="body" id="chat-body"></div>
        <form class="input-area">
            <input type="text" placeholder="..." class="chat-input" required>
            <button type="submit" class="send-btn"><i class="fa-solid fa-paper-plane"></i></button>
        </form>
    `
    body.append(div)

    div.querySelector('.input-area').addEventListener('submit', (e) => {
        e.preventDefault()
        const input = div.querySelector('.chat-input')
        const message = input.value.trim()
        if (message !== "") {
            // console.log(1);
            
            webSockert(userId, receiverId, input.value)
            input.value = ''
        }
    })
}
