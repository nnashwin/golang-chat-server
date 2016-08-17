var checkStringMatch = function (string1, string2) {
    if (typeof string1 !== "string" || typeof string2 !== "string") {
        return false
    }

    if (string1 === string2) {
        return true
    }

    return false
} 

document.addEventListener("DOMContentLoaded", () => {
    let form = document.getElementById('newUserForm')
    form.addEventListener('submit', (e) => {
        e.preventDefault()

        if (checkStringMatch(form.pass.value, form.confirmPass.value) === false) {
            return console.error('doesnt match')
        }

        let data = {}
        data["username"] = form.username.value
        data["password"] = form.pass.value

        let req = new XMLHttpRequest()

        let failure = (reqStatus) => {
            console.error("Error: " + reqStatus)
        }

        req.onreadystatechange = () => {
            const DONE = 4
            const OK = 201
            if (req.readyState === DONE) {
                if (req.status === OK) {
                    console.log(req) 
                    console.log(req.status) 
                } else {
                    failure(req.status)
                }
            }
        }

        req.open("POST", "/api/v1/users/signup")

        req.setRequestHeader("X-Requested-With", "XMLHttpRequest")
        req.setRequestHeader("Content-Type", "application/json;charset=UTF-8")

        req.send(JSON.stringify(data))
    })
})
