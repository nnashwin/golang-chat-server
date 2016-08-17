var checkIfMatch = function (string1, string2) {
    if (string1 === string2) {
        return true
    }
    return false
} 

document.addEventListener("DOMContentLoaded", () => {
    let form = document.getElementById('newUserForm')
    form.addEventListener('submit', (e) => {
        e.preventDefault()
        console.log(checkIfMatch(form.pass, form.confirmPass))

        if (checkIfMatch(form.pass, form.confirmPass) === false) {
            return console.log('no good!!')
        }

        let data = {}



        let req = new XMLHttpRequest()
    })


})
