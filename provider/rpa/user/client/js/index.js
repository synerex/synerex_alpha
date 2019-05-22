$(() => {
    var socket = io();

    socket.on("server_to_client", (data) => {
        console.log("server_to_client:", data);
    });

    $("#send").on("click", () => {

        var datetime = $("#datetimepicker").val();

        $("#send").prop("disabled", true);
        $("#datetimepicker").prop("disabled", true);

        if (datetime == "") {
            $("#send").prop("disabled", false);
            $("#datetimepicker").prop("disabled", false);
            console.log("You must set the date and time.")
        } else {
            socket.emit("client_to_server", { value: datetime });
        }
    });

    // datetime-picker
    $("#datetimepicker").datetimepicker({
        icons: {
            time: 'far fa-clock',
            date: 'far fa-calendar',
            up: 'fas fa-arrow-up',
            down: 'fas fa-arrow-down',
            previous: 'fas fa-chevron-left',
            next: 'fas fa-chevron-right',
            today: 'fas fa-calendar-check',
            clear: 'far fa-trash-alt',
            close: 'far fa-times-circle'
        },
    });
});