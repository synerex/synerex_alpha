// クエリ文字列の取得
function q(name, url) {
    if (!url) url = window.location.href;
    name = name.replace(/[\[\]]/g, "\\$&");
    var results = new RegExp("[?&]" + name + "(=([^&#]*)|&|#|$)").exec(url);
    if (!results) return null;
    if (!results[2]) return '';
    return decodeURIComponent(results[2].replace(/\+/g, " "));
}

// Socket接続
$(() => {
    var socket = null;

    // イベント送信
    function emit(name, data) {
        if (socket && socket.connected) {
            socket.emit(name, data);
            console.log("送信メッセージ(" + name + "): " + data);
        }
    }

    // 接続・切断
    $("#connect").click(function () {
        socket = io({ transports: ["websocket"] });

        socket.on("disp_start", function (data) {
            console.log("受信メッセージ: " + data)

            // StringからJSONにパースする
            data = JSON.parse(data);
            console.dir(data);

            // 広告 or アンケート で場合分け
            const contents = data.contents[0];
            switch (contents.type) {
                case 'AD':
                    console.log('case "AD" is called');
                    const ad = contents;
                    console.dir(ad);
                    $('#ad-area').children('img').attr('src', ad.data);
                    emit("disp_complete", { command: "RESULTS", results: null });

                    break;

                case 'ENQ':
                    console.log('case "ENQ" is called');
                    const questions = contents;

                    const div = [];

                    // ForEach文
                    Object.keys(questions).forEach((key) => {

                        // Typeで場合分け
                        switch (questions[key].type) {

                            case "select":

                                div[key] = $('<div></div>', { addClass: "form-group" });
                                div[key].append('<label for="' + questions[key].name + '">' + questions[key].label + '</label>');

                                const select = $('<select></select>', {
                                    name: questions[key].name,
                                    id: questions[key].name,
                                    addClass: "form-control"
                                });

                                for (let value of questions[key].option.options) {
                                    select.append('<option value="' + value.value + '">' + value.text + '</option>');
                                }

                                div[key].append(select);
                                break;

                            case "checkbox":

                                div[key] = $('<div></div>', { addClass: "form-group" });
                                div[key].append('<p>' + questions[key].label + '</p>');

                                questions[key].option.options.forEach((value, index) => {
                                    const divCheck = $('<div></div>', { addClass: "form-check" });
                                    divCheck.append('<input class="form-check-input" type="checkbox" name="' + questions[key].name
                                        + '" id="' + questions[key].name + index
                                        + '" value="' + value.value
                                        + '">');
                                    divCheck.append('<label class="form-check-label" for="' + questions[key].name + index + '">' + value.text + '</label>');
                                    div[key].append(divCheck);
                                });

                                break;

                            case "range":

                                div[key] = $('<div></div>', { addClass: "form-group" });
                                div[key].append('<label>' + questions[key].label + '</label>');

                                const container = $('<div></div>', { addClass: "container" });
                                const row = $('<div></div>', { addClass: "row" });

                                const col = [];
                                col[0] = $('<div></div>', { addClass: "col" });
                                col[1] = $('<div></div>', { addClass: "col" });
                                col[2] = $('<div></div>', { addClass: "col" });

                                col[0].append('<p class="text-right">' + questions[key].option.minText + '</p>');

                                col[1].append('<input class="custom-range" type="range" name="' + questions[key].name
                                    + '" name="' + questions[key].name
                                    + '" max="' + questions[key].option.max
                                    + '" min="' + questions[key].option.min
                                    + '">');

                                col[2].append('<p>' + questions[key].option.maxText + '</p>');

                                for (let value of col) {
                                    row.append(value);
                                }

                                container.append(row);
                                div[key].append(container);

                                break;

                            case "textarea":

                                div[key] = $('<div></div>', { addClass: "form-group" });
                                div[key].append('<label for="' + questions[key].name + '">' + questions[key].label + '</label>');
                                div[key].append('<textarea name="' + questions[key].name
                                    + '" id="' + questions[key].name
                                    + '" class="form-control'
                                    + '" placeholder="' + questions[key].option.placeholder
                                    + '"></textarea>');

                                break;

                            default:
                                console.log(`switch文に case"${questions[key].type}" を追記してください。`);
                                break;
                        }

                    });

                    // <form>に<div>を追加
                    for (let value of div) {
                        $('form#questions').append(value);
                    }

                    // <form>の最下部にbuttonを追加
                    $('form#questions').append('<button type="button" id="submitJson" class="btn btn-primary btn-lg btn-block">送信</button>');

                    // JSON形式で送信
                    $('button#submitJson').on('click', () => {
                        const json = $('form#questions').serializeArray();
                        Object.keys(json).forEach((key) => {
                            console.log(json[key]);
                        });
                        alert('ありがとうございました！');
                        $('form#questions')[0].reset();
                    });

                    break;

                default:
                    console.log('case "default" is called');
                    break;
            }
            // ここまで広告・アンケートの場合分け


        });
    });
    $("#disconnect").click(function () {
        socket.close()
    });

    // 搭載車両登録
    $("#register").click(function () {
        emit("disp_register", { taxi: $("#taxi").val(), disp: $("#disp").val() });
    });
    // 出発
    $("#depart").click(function () {
        emit("depart", { taxi: $("#taxi").val() });
    });
    // 到着
    $("#arrive").click(function () {
        emit("arrive", { taxi: $("#taxi").val() });
    })
    // 完了
    $("#complete").click(function () {
        emit("disp_complete", { command: "RESULTS", results: null });
    });

    // タクシー・ディスプレイ設定 (あれば)
    var taxi = q("taxi"), disp = q("disp");
    if (taxi) $("#taxi").val(taxi);
    if (disp) $("#disp").val(disp);

});