$(() => {
    var socket = io();

    socket.on("server_to_client", (data) => {
        var booking_results = $("#booking_results");
        var list = `<li>${data}</li>`;
        booking_results.append(list);

        $("#confirm_message").remove();
        $("#datetimepicker").val("");
        $("#send").prop("disabled", false);
        $("#datetimepicker").prop("disabled", false);
    });

    socket.on("check_booking", (msg) => {
        var confirm_area = $("#confirm_area");
        var contents = `<div id="confirm_message" class="row"><div class="col-sm-6"><p>${msg}</p></div><div class="col-sm-6"><button type="button" id="yes" class="btn btn-primary">YES</button><button type="button" id="stop" class="btn btn-secondary">STOP</button></div></div>`;
        confirm_area.append(contents);
    });

    $(document).on("click", "#yes", () => {
        socket.emit("confirm_booking", "yes");
    });
    $(document).on("click", "#stop", () => {
        socket.emit("confirm_booking", "stop");
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

    // pickadate.js
    $(".datepicker").pickadate();
    $(".timepicker").pickatime({
        format: 'H:i',
    });

});