import { main } from "./main.js";

export function Error(err) {
    const style = document.createElement("style");
    style.textContent = `
    @import url("https://fonts.googleapis.com/css?family=Montserrat:400,400i,700");

*,
*:after,
*:before {
    box-sizing: border-box;
}

body {
    background-color: #313942;
    font-family: 'Montserrat', sans-serif;
}

body {
    align-items: center;
    display: flex;
    flex-direction: column;
    height: 100vh;
    justify-content: center;
    text-align: center;
}

h1 {
    color: #e7ebf2;
    font-size: 12.5rem;
    letter-spacing: .10em;
    margin: .025em 0;
    text-shadow: .05em .05em 0 rgba(0, 0, 0, .25);
    white-space: nowrap;
}

h1>span {
    animation: spooky 2s alternate infinite linear;
    color: #528cce;
    display: inline-block;
}

@media (max-width: 30rem) {
    h1 {
        font-size: 8.5rem;
    }
}

h2 {
    color: #e7ebf2;
    margin-bottom: .40em;
}

p {
    color: #ccc;
    margin-top: 0;
}

@keyframes spooky {
    from {
        transform: translatey(.15em) scaley(.95);
    }

    to {
        transform: translatey(-.15em);
    }
}

    `
    const error = err.split('')
    document.head.appendChild(style);
    document.body.innerHTML = `
                     <button id="Back" style="z-index: 10;">Home</button>
                      <h1>${error[0]}<span><i class="fas fa-ghost"></i></span>${error[2]}</h1>
                      <h2>Error: ${err} page not found</h2>
                      <p>Sorry, the page you're looking for cannot be accessed</p>
                      
                    </main >
        <script type="module" src="static/main.js"></script>
    `
    const button = document.getElementById('Back')
    button.style.position = 'absolute'
    button.style.top = '70%'
    button.addEventListener('click', (e) => {
        e.preventDefault()
        history.pushState({}, '', 'http://localhost:8080')
        main()

    })

}