// JSONを読み込む
$.getJSON('../questions.json', function (data) {
    let questions = data.questions;
    const div = [];

    // ForEach文
    Object.keys(questions).forEach(function (key) {

        // Typeで場合分け
        switch (questions[key].type) {

            case "select":

                div[key] = $('<div></div>', { addClass: "" });
                const select = $('<select></select>', {
                    name: questions[key].name,
                    id: questions[key].name,
                    addClass: ""
                });

                div[key].append('<label for="' + questions[key].name + '">' + questions[key].label + '</label>');
                for (let value of questions[key].option.options) {
                    select.append('<option value="' + value.value + '">' + value.text + '</option>');
                }

                div[key].append(select);
                break;

            case "checkbox":

                div[key] = $('<div></div>', { addClass: "" });
                div[key].append('<label for="' + questions[key].name + '">' + questions[key].label + '</label>');

                for (const value of questions[key].option.options) {
                    div[key].append('<input type="checkbox" name="' + questions[key].name + '" value="' + value.value + '">' + value.text);
                }
                break;

            case "range":

                div[key] = $('<div></div>', { addClass: "" });
                div[key].append('<label for="' + questions[key].name + '">' + questions[key].label + '</label>');
                div[key].append(questions[key].option.minText
                    + '<input type="range" name="' + questions[key].name
                    + '" name="' + questions[key].name
                    + '" max="' + questions[key].option.max
                    + '" min="' + questions[key].option.min
                    + '">'
                    + questions[key].option.maxText);
                break;

            case "textarea":

                div[key] = $('<div></div>', { addClass: "" });
                div[key].append('<label for="' + questions[key].name + '">' + questions[key].label + '</label>');
                div[key].append('<textarea name="' + questions[key].name
                    + '" id="' + questions[key].name
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
    $('form#questions').append('<button type=button onClick="alert(\'button was pushed!\');">送信する</button>');

});