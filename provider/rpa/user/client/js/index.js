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
            var splits = datetime.split(' ');
            var date = splits[0];
            var time = splits[1];
            var ampm = splits[2];

            splits = date.split('/');
            var month = splits[0];
            var day = splits[1];
            var year = splits[2];

            splits = time.split(':');
            var hour = splits[0];
            var minute = splits[1];

            if (ampm == "PM") {
                hour = parseInt(hour, 10);
                hour += 12;
            }

            socket.emit("client_to_server", {
                Year: year,
                Month: month,
                Day: day,
                Hour: hour,
                Minute: minute
            });
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