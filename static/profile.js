import { login } from "./login.js"

export function profile(age, email, nickname, photo) {
    const body = document.querySelector('body')
    body.innerHTML = `
    <div id="content">
        <header>
            <button id="Back" style="z-index: 10;"><i class="fa-solid fa-arrow-right"></i></button>
        </header>
        <div class="profile-container">
            <form id="edit">
                <img src="${photo}" alt="Profile Picture">
                <br>
                <input type="file" name="photo" accept="image/*">
                <h2>${nickname} </h2>
                <input type="text" name="nickname" id="title">
                <p>Email: ${email} </p>
                <input type="email" name="email" id="title">
                <p>Age: ${age} </p>
                <input type="number" name="age" id="title">
                <div class="profile-actions">
                    <button type="submit" id="editProfile">Edit Profile</button>
                </div>
            </form>
        </div>
    </div>
        <script type="module" src="static/script.js"></script>
    `
    document.getElementById('edit').addEventListener('submit', (e) => {
        e.preventDefault();

        const form = e.target;
        const formData = new FormData(form);

        fetch('/editProfile', {
            method: 'POST',
            body: formData,
        }).then(r => {
            console.log(r);

        })
    })
    document.getElementById('Back').addEventListener('click', () => {
        login()
    })
}