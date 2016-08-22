document.addEventListener("DOMContentLoaded", () => {
    var isLoggedIn = false
    var checkLogin = function (loginBool) {
        if (checkAuth(localStorage.getItem('CTString')) === true) {
            loginBool = true
        } else if (checkAuth(localStorage.getItem('CTString')) === false){
            loginBool = false
        }
        if (loginBool === false) {
            window.location = "/login"
        }
    }

    var checkAuth = function (jwtStr) {
        let data = {}
        data["Token"] = jwtStr

        let req = new XMLHttpRequest() 

        let failure = (reqStatus) => {
            console.error(`Error: ${reqStatus}`)
        }

        req.onreadystatechange = () => {
            const DONE = 4
            const OK = 200
            if (req.readyState === DONE) {
                if (req.status === OK) {
                    return true
                } else {
                    return false
                }
            }
        }

        req.open("POST", "/api/v1/users/checkAuth")

        req.setRequestHeader("X-Requested-With", "XMLHttpRequest")
        req.setRequestHeader("Content-Type", "application/json;charset=UTF-8")

        req.send(JSON.stringify(data))
    }

    checkLogin(isLoggedIn)
    window.setInterval(checkLogin, 240000, isLoggedIn)


    var url = "ws://" + window.location.host + "/auth/ws";
    var ws = new WebSocket(url);
    var name = "Guest" + Math.floor(Math.random() * 1000);
    var chat = document.getElementById("chat");
    var text = document.getElementById("text");
    var now = function () {
        var iso = new Date().toISOString();
        return iso.split("T")[1].split(".")[0];
    };
    ws.onmessage = function (msg) {
        var line =  now() + " " + msg.data + "\n";
        chat.innerText += line;
    };
    text.onkeydown = function (e) {
        if (e.keyCode === 13 && text.value !== "") {
            ws.send("<" + name + "> " + text.value);
            text.value = "";
        }
    };
})
