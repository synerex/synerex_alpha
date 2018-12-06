// クエリ文字列の取得
function q(name, url) {
    if (!url) url = window.location.href;
    name = name.replace(/[\[\]]/g, "\\$&");
    var results = new RegExp("[?&]" + name + "(=([^&#]*)|&|#|$)").exec(url);
    if (!results) return null;
    if (!results[2]) return '';
    return decodeURIComponent(results[2].replace(/\+/g, " "));
}

// モーダルダイアログ
$(() => {
    $('#connection-area').dialog({
        modal: true,
        title: "初期設定",
        buttons: {
            "閉じる": function () {
                $(this).dialog("close");
            }
        }
    });
});

$(() => {
    var socket = null;

    // イベント送信
    function emit(name, data) {
        if (socket && socket.connected) {
            socket.emit(name, data);
            console.log("送信メッセージ(" + name + "): ", data);
        }
    }

    // 接続・切断
    $("#connect").click(function () {
        socket = io({ transports: ["websocket"] });

        socket.on("disp_start", function (data) {
            console.log("受信メッセージ: ", data)

            // StringからJSONにパースする
            data = JSON.parse(data);

            // 広告 or アンケート で場合分け
            const contents = data.contents[0];
            switch (contents.type) {
                case 'AD':
                    const ad = data.contents;

                    // アンケートの表示を終了する
                    $('form#questions').empty();

                    waiting();

                    // ループ処理を指定の秒数だけ待つための処理
                    function waiting() {

                        // 全ての広告を表示したらループを終わる
                        if (ad.length == 0) {
                            emit("disp_complete", { command: "RESULTS", results: null });
                            return;
                        }

                        // 配列の先頭の広告を読み込む
                        const param = ad[0];

                        // 広告を表示する
                        // $('#ad-area').children('img').attr('src', param.data);
                        $('.box').css('background-image', `url("${param.data}")`);
                        $('#ad-area').addClass('box');
                        console.log(`広告表示中: ${param.data}`);

                        // 表示し終わった広告を配列から取り除き、次に読み込む広告を先頭にする
                        ad.shift();

                        // 広告をperiod秒間表示してから次のループに移る
                        setTimeout(() => {
                            waiting();
                        }, param.period * 1000);
                    }

                    break;

                case 'ENQ':
                    const questions = data.contents[0].data.questions;
                    const div = [];

                    // 広告の表示を終了する
                    // $('#ad-area').children('img').attr('src', '');
                    $('#ad-area').removeClass('box');

                    // ForEach文
                    Object.keys(questions).forEach((i) => {

                        // Typeで場合分け
                        switch (questions[i].type) {

                            case "select":

                                div[i] = $('<div></div>', { addClass: "form-group" });
                                div[i].append('<label for="' + questions[i].name + '">' + questions[i].label + '</label>');

                                const select = $('<select></select>', {
                                    name: questions[i].name,
                                    id: questions[i].name,
                                    addClass: "form-control form-control-lg"
                                });

                                for (let value of questions[i].option.options) {
                                    select.append('<option value="' + value.value + '">' + value.text + '</option>');
                                }

                                div[i].append(select);
                                break;

                            case "checkbox":

                                div[i] = $('<div></div>', { addClass: "form-group" });
                                div[i].append('<p>' + questions[i].label + '</p>');

                                // checkedされない場合のhidden checkbox
                                div[i].append(`<input type="hidden" name="${questions[i].name}" id="${questions[i].name}_hidden" value="">`);

                                questions[i].option.options.forEach((value, index) => {
                                    const divCheck = $('<div></div>', { addClass: "form-check" });
                                    divCheck.append('<input class="form-check-input" type="checkbox" name="' + questions[i].name
                                        + '" id="' + questions[i].name + index
                                        + '" value="' + value.value
                                        + '">');
                                    divCheck.append('<label class="form-check-label" for="' + questions[i].name + index + '">' + value.text + '</label>');
                                    div[i].append(divCheck);
                                });

                                break;

                            case "range":

                                div[i] = $('<div></div>', { addClass: "form-group" });
                                div[i].append('<label>' + questions[i].label + '</label>');

                                const container = $('<div></div>', { addClass: "container" });
                                const row = $('<div></div>', { addClass: "row" });

                                const col = [];
                                col[0] = $('<div></div>', { addClass: "col" });
                                col[1] = $('<div></div>', { addClass: "col" });
                                col[2] = $('<div></div>', { addClass: "col" });

                                col[0].append('<p class="text-right">' + questions[i].option.minText + '</p>');

                                col[1].append('<input class="custom-range" type="range" name="' + questions[i].name
                                    + '" name="' + questions[i].name
                                    + '" max="' + questions[i].option.max
                                    + '" min="' + questions[i].option.min
                                    + '" step="' + 0.1
                                    + '">');

                                col[2].append('<p>' + questions[i].option.maxText + '</p>');

                                for (let value of col) {
                                    row.append(value);
                                }

                                container.append(row);
                                div[i].append(container);

                                break;

                            case "textarea":

                                div[i] = $('<div></div>', { addClass: "form-group" });
                                div[i].append('<label for="' + questions[i].name + '">' + questions[i].label + '</label>');
                                div[i].append('<textarea name="' + questions[i].name
                                    + '" id="' + questions[i].name
                                    + '" class="form-control'
                                    + '" placeholder="' + questions[i].option.placeholder
                                    + '"></textarea>');

                                break;

                            default:
                                console.log(`switch文に case"${questions[i].type}" を追記してください。`);
                                break;
                        }

                    });

                    // 配列divに格納したinputを<form>に追加
                    for (let value of div) {
                        $('form#questions').append(value);
                    }

                    // <form>の最下部にbuttonを追加
                    $('form#questions').append('<button type="button" id="submitJson" class="btn btn-primary btn-lg btn-block">送信</button>');

                    // <form>の上下にpaddingを追加
                    $('form#questions').css('padding', '30px 0');

                    // JSON形式で送信
                    $('button#submitJson').on('click', () => {
                        const serialized = $('form#questions').serializeArray();
                        const hash = {};

                        Object.keys(serialized).forEach((i) => {
                            const key = serialized[i].name;
                            const value = serialized[i].value;
                            const array = [];

                            // keyが重複するか判定する
                            if (key in hash) {

                                // 既に存在するkeyのvalueを配列に保存する
                                for (let val of hash[key]) {
                                    array.push(val);
                                }

                                // 新たに追加するvalueを配列に格納する
                                array.push(value);

                                // 全てのvalueを格納する
                                hash[key] = array;

                            } else {
                                // keyが重複しないので普通にvalueを格納
                                hash[key] = value;
                            }
                        });

                        // JSON形式に整形する
                        let answers = [];
                        for (let key in hash) {
                            let obj = {};
                            if (typeof hash[key] == 'object') {
                                obj = {
                                    name: key,
                                    value: hash[key]
                                };
                            } else {
                                obj = {
                                    name: key,
                                    value: [
                                        hash[key]
                                    ]
                                };
                            }
                            answers.push(obj);
                        }
                        answers = JSON.stringify(answers);

                        emit("disp_complete", { command: "RESULTS", results: { answers } });
                        alert('ありがとうございました！');
                        $('form#questions')[0].reset();
                    });

                    break;

                default:
                    console.log('case "default" is called');
                    break;
            } // ここまで広告・アンケートの場合分け

        });
    });
    $("#disconnect").click(function () {
        socket.close()
    });

    // 搭載車両登録
    $("#register").click(function () {
        emit("disp_register", { taxi: $("#taxi").val(), disp: $("#disp").val() });
    });
    // 完了
    $("#complete").click(function () {
        emit("disp_complete", { command: "RESULTS", results: null });
    });

    // 出発
    $("#depart").click(function () {
        emit("depart", { taxi: $("#taxi").val() });
    });
    // 到着
    $("#arrive").click(function () {
        emit("arrive", { taxi: $("#taxi").val() });
    })

    // タクシー・ディスプレイ設定 (あれば)
    var taxi = q("taxi"), disp = q("disp");
    if (taxi) $("#taxi").val(taxi);
    if (disp) $("#disp").val(disp);
});