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
                <p>Nickname : ${nickname} </p>
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
        })
            .then(r => {
                if (!r.ok) {
                    return r.json().then(errorData => {
                        throw new Error(errorData.Text || `HTTP error! Status: ${r.status}`);
                    });
                }
                return r.json();
            })
            .then(r => {
                const existingPopup = document.querySelector(".content");
                if (existingPopup) {
                    existingPopup.remove();
                }

                const popupDiv = document.createElement('div');
                popupDiv.className = 'popup-container';
                popupDiv.innerHTML = `
                        <div class="content">
                        Your information has been updated
                        </div>`
                document.getElementById('content').append(popupDiv)

                const editProfileForm = document.getElementById('edit');
                if (editProfileForm) {
                    editProfileForm.reset();
                }
            }).catch(err => {

                const existingPopup = document.querySelector(".content");
                if (existingPopup) {
                    existingPopup.remove();
                }
                const ErrorDiv = document.createElement('div');
                ErrorDiv.className = 'error-container';
                ErrorDiv.innerHTML = `<div class="content">${err.message}</div>`;
                document.querySelector('body').append(ErrorDiv);
                setTimeout(() => {
                    ErrorDiv.remove()
                }, 1000)
            }

            )
    })
    document.getElementById('Back').addEventListener('click', () => {
        login()
    })
}