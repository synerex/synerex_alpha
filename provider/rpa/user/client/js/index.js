$(() => {
    var socket = io();

    socket.on("server_to_client", (data) => {
        var booking_results = $("#booking_results");
        var list = `<li>${data}</li>`;
        booking_results.append(list);

        $("#booking_options").empty();

        resetValues();
        changeAllProps();
    });

    socket.on("check_booking", (msg) => {
        var booking_options = $("#booking_options");

        var contents = `
            <p>${msg}<span><button type="button" id="yes" class="btn btn-primary">Select</button></span></p>
        `;
        booking_options.append(contents);
    });

    $(document).on("click", "#yes", (e) => {
        $(e.target).prop("disabled", true);
        $("button").prop("disabled", true);
        socket.emit("confirm_booking", "yes");
    });

    $("#send").on("click", (e) => {
        e.preventDefault();
        changeAllProps();

        var date = $("#booking-date").val();
        var start = $("#booking-start").val();
        var end = $("#booking-end").val();
        var people = $("#number-people").val();
        var title = $("#booking-title").val();

        if (date == "" || start == "" || end == "") {
            changeAllProps();
            console.log("You must set the date and time.")
        } else {
            var splits = date.split(' ');
            var week = splits[1];

            splits = splits[0].split('/');
            var year = splits[0];
            var month = splits[1];
            var day = splits[2];

            socket.emit("client_to_server", {
                Year: year,
                Month: month,
                Day: day,
                Week: week,
                Start: start,
                End: end,
                People: people,
                Title: title,
            });
        }
    });

    // pickadate.js
    $(".datepicker").pickadate({
        format: 'yyyy/m/d (ddd)',
    });
    $(".timepicker").pickatime({
        format: 'H:i',
    });

});

const changeAllProps = () => {
    var flag = $("#send").prop("disabled");

    if (flag) {
        $("#booking-date").prop("disabled", false);
        $("#booking-start").prop("disabled", false);
        $("#booking-end").prop("disabled", false);
        $("#number-people").prop("disabled", false);
        $("#booking-title").prop("disabled", false);
        $("#send").prop("disabled", false);
    } else {
        $("#booking-date").prop("disabled", true);
        $("#booking-start").prop("disabled", true);
        $("#booking-end").prop("disabled", true);
        $("#number-people").prop("disabled", true);
        $("#booking-title").prop("disabled", true);
        $("#send").prop("disabled", true);
    }
}

const resetValues = () => {
    $("#booking-date").val("");
    $("#booking-start").val("");
    $("#booking-end").val("");
    $("#number-people").val("");
    $("#booking-title").val("");
}
