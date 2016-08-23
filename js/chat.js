document.addEventListener("DOMContentLoaded", () => {
    var checkAuth = function () {
        let data = {}
        data["Token"] = localStorage.getItem("CTString")

        let req = new XMLHttpRequest() 

        let failure = (reqStatus) => {
            console.error(`Error: ${reqStatus}`)
        }

        req.onreadystatechange = () => {
            const DONE = 4
            const OK = 200
            if (req.readyState === DONE) {
                if (req.status === OK) {
                    console.log('checkedOut')
                    return true
                } else {
                    console.log("notCheckedOut")
                    return false
                }
            }
        }

        var checkLogin = function () {
            let isLoggedIn = false
            if (checkAuth(localStorage.getItem('CTString')) === true) {
                isLoggedIn = true
            } else if (checkAuth(localStorage.getItem('CTString')) === false){
                isLoggedIn = false
            }

            if (isLoggedIn === false) {
                //window.location = "/login"
            }
        }

        req.open("POST", "/api/v1/users/checkAuth")

        req.setRequestHeader("X-Requested-With", "XMLHttpRequest")
        req.setRequestHeader("Content-Type", "application/json;charset=UTF-8")

        req.send(JSON.stringify(data.Token))
    }

    checkLogin()
    window.setInterval(checkLogin, 240000)


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
