document.addEventListener("DOMContentLoaded", () => {
    let form = document.getElementById('loginForm')
    form.addEventListener('submit', (e) => {
        e.preventDefault()

        let data = {}
        data["username"] = form.username.value
        data["password"] = form.password.value

        let req = new XMLHttpRequest()

        let failure = (reqStatus) => {
            console.error("Error: " + reqStatus)
        }

        req.onreadystatechange = () => {
            const DONE = 4
            const OK = 200
            if (req.readyState === DONE) {
                if (req.status === OK) {
                    localStorage.removeItem('CTString')
                    localStorage.setItem('CTString', req.response)
                    //window.location = "/"
                } else {
                    failure(req.status)
                }
            }
        }

        req.open("POST", "/api/v1/users/login")

        req.setRequestHeader("X-Requested-With", "XMLHttpRequest")
        req.setRequestHeader("Content-Type", "application/json;charset=UTF-8")

        req.send(JSON.stringify(data))
    })
})
