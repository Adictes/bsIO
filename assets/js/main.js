var shs = new WebSocket("ws://localhost:8080/shs");

shs.onmessage = function (event) {
    var data = JSON.parse(event.data);

    document.getElementById("single-decker").innerText = data.SingleDecker
    document.getElementById("two-decker").innerText = data.TwoDecker
    document.getElementById("three-decker").innerText = data.ThreeDecker
    document.getElementById("four-decker").innerText = data.FourDecker
}

function send(cell) {
    shs.send(cell);
    changeCellColor(cell);
}

function changeCellColor(id) {
    if (document.getElementById(id).style.backgroundColor == "") {
        document.getElementById(id).style.backgroundColor = "#337ab7";
    } else {
        document.getElementById(id).style.backgroundColor = "";
    }
}

var hes = new WebSocket("ws://localhost:8080/hes");

function hit(cell) {
    hes.send(cell);
}

var clr = new WebSocket("ws://localhost:8080/clr");

clr.onmessage = function (event) {
    var data = JSON.parse(event.data);

    if (data.Leave != undefined) {
        cleanAll();
        alert('Твой противник вышел')
    }
}

hes.onmessage = function (event) {
    var data = JSON.parse(event.data);

    if (data.Turn != undefined) {
        if (data.Turn == true) {
            document.getElementById("turn").innerText = "Ходит";
            document.getElementById("eturn").innerText = "Ждет";
            document.getElementById("turn").classList.remove("badge-warning");
            document.getElementById("turn").classList.add("badge-success");
            document.getElementById("eturn").classList.remove("badge-success");
            document.getElementById("eturn").classList.add("badge-warning");
        } else {
            document.getElementById("turn").innerText = "Ждет";
            document.getElementById("eturn").innerText = "Ходит";
            document.getElementById("turn").classList.remove("badge-success");
            document.getElementById("turn").classList.add("badge-warning");
            document.getElementById("eturn").classList.remove("badge-warning");
            document.getElementById("eturn").classList.add("badge-success");
        }
    } else if (data.Hitted != undefined) {
        if (data.Hitted != "") {
            document.getElementById(data.Hitted).style.backgroundColor = "#990000";
            document.getElementById(data.Hitted).removeAttribute("onclick");
        }

        for (i in data.Ambient) {
            if (data.Ambient[i][1] == '-' || data.Ambient[i][3] == '-' ||
                (data.Ambient[i][1] == '1' && data.Ambient[i][2] == '0') ||
                (data.Ambient[i][3] == '1' && data.Ambient[i][4] == '0')) {
                continue;
            }
            document.getElementById(data.Ambient[i]).style.backgroundColor = "#e6ffff";
            document.getElementById(data.Ambient[i]).removeAttribute("onclick");
        }
    } else if (data.SingleDecker != undefined) {
        document.getElementById("esingle-decker").innerText = data.SingleDecker;
        document.getElementById("etwo-decker").innerText = data.TwoDecker;
        document.getElementById("ethree-decker").innerText = data.ThreeDecker;
        document.getElementById("efour-decker").innerText = data.FourDecker;
    } else if (data.Win != undefined) {
        clr.send("Do cleaning");
        cleanAll();
        if (data.Win == true) {
            alert('Ты победил!');
        } else {
            alert('Ты проиграл.');
        }
    }
}

function cleanAll() {
    document.getElementById("turn").innerText = "Ждет";
    document.getElementById("eturn").innerText = "Ждет";
    document.getElementById("ename").innerText = "Противник";

    if (document.getElementById("turn").classList.contains('badge-warning') !== true) {
        document.getElementById("turn").classList.remove('badge-success');
        document.getElementById("turn").classList.add('badge-warning');
    }
    if (document.getElementById("eturn").classList.contains('badge-warning') !== true) {
        document.getElementById("eturn").classList.remove('badge-success');
        document.getElementById("eturn").classList.add('badge-warning');
    }

    var cells = document.getElementsByClassName("hf");
    for (i = 0; i < cells.length; i++) {
        cells[i].style.backgroundColor = "";
        cells[i].setAttribute('onclick', 'send(\'' + cells[i].id + '\')');
    }

    var ecells = document.getElementsByClassName("ef");
    for (i = 0; i < ecells.length; i++) {
        ecells[i].style.backgroundColor = "";
        ecells[i].setAttribute('onclick', 'hit(\'' + ecells[i].id + '\')');
    }

    document.getElementById("single-decker").innerText = 4;
    document.getElementById("two-decker").innerText = 3;
    document.getElementById("three-decker").innerText = 2;
    document.getElementById("four-decker").innerText = 1;

    document.getElementById("esingle-decker").innerText = 0;
    document.getElementById("etwo-decker").innerText = 0;
    document.getElementById("ethree-decker").innerText = 0;
    document.getElementById("efour-decker").innerText = 0;

    document.getElementById("stg").disabled = false;
    document.getElementById("rff").disabled = false;

    document.getElementById("output").innerText = "";
}

var stg = new WebSocket("ws://localhost:8080/stg");

stg.onmessage = function (event) {
    var data = JSON.parse(event.data);

    if (data.Correctness != undefined) {
        if (data.Correctness == true) {
            document.getElementById("stg").disabled = true;
            document.getElementById("rff").disabled = true;

            var cells = document.getElementsByClassName("hf");
            for (i in cells) {
                cells[i].removeAttribute("onclick");
            }
        } else {
            alert("Корабли расставлены не верно");
        }
    } else if (data.Name != undefined) {
        document.getElementById("ename").innerText = data.Name;
    } else if (data.Turn != undefined) {
        if (data.Turn == true) {
            document.getElementById("turn").innerText = "Ходит";
            document.getElementById("turn").classList.remove("badge-warning");
            document.getElementById("turn").classList.add("badge-success");
        } else {
            document.getElementById("eturn").innerText = "Ходит";
            document.getElementById("eturn").classList.remove("badge-warning");
            document.getElementById("eturn").classList.add("badge-success");
        }
    }
}

var rff = new WebSocket("ws://localhost:8080/rff");

rff.onmessage = function (event) {
    var data = JSON.parse(event.data);

    if (data.Clear != undefined) {
        var cells = document.getElementsByClassName("hf");
        for (i = 0; i < cells.length; i++) {
            cells[i].style.backgroundColor = "";
        }
    } else if (data.Cell != undefined) {
        changeCellColor(data.Cell);
    } else if (data.SingleDecker != undefined) {
        document.getElementById("single-decker").innerText = data.SingleDecker;
        document.getElementById("two-decker").innerText = data.TwoDecker;
        document.getElementById("three-decker").innerText = data.ThreeDecker;
        document.getElementById("four-decker").innerText = data.FourDecker
    }
}

function randomFilling() {
    rff.send("Start the random filling");
}

function startTheGame() {
    stg.send("The game has begun");
}

//Chat functionality

var hm = new WebSocket("ws://localhost:8080/hm");

hm.onmessage = function (event) {
    var data = JSON.parse(event.data);

    document.getElementById("output").innerText += data.Name + ': ' + data.Message;
}

function sendMessage() {
    var message = document.getElementById("t_input").value + '\n';

    document.getElementById("output").innerHTML += `Me: ${message}<br>`;
    hm.send(message);

    document.getElementById("t_input").value = "";
    document.getElementById("t_input").focus();
}

function keyPressed(event) {
    var key = event.which;
    if (key == 13) {
        sendMessage();
    }
}

function clearMessage() {
    document.getElementById("t_input").value = "";
    document.getElementById("t_input").focus();
}