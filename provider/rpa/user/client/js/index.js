$(() => {
    var socket = io();

    $("#disconnect").on("click", () => {
        socket.close()
        console.log("Socket.IO Disconnected.")
    });

});