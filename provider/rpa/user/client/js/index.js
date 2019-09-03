$(() => {
    var socket = io();

    socket.on("server_to_client", (data) => {
        console.log(data.arg_json)
        var booking_results = $("#booking_results");
        var list = `<li>Success booking: ${data.arg_json}</li>`;
        booking_results.append(list);

        $("#booking_options").empty();

        resetValues();
        changeAllProps();
    });

    socket.on("check_booking", (json) => {
        var booking_options = $("#booking_options");

        var obj = JSON.parse(json)

        var contents = `
            <p>${obj.room}<span><button type="button" id="yes" data-id="${obj.id}" class="btn btn-primary">Select</button></span></p>
        `;
        booking_options.append(contents);
    });

    $(document).on("click", "#yes", (e) => {
        $(e.target).prop("disabled", true);
        $("button").prop("disabled", true);
        socket.emit("confirm_booking", e.target.dataset.id);
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
