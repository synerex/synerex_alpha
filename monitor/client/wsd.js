var wsd = require('websequencediagrams');
var fs = require('fs');

const text = fs.readFileSync('./logs.txt', 'utf8');
console.log(`file contents:\n${text}`);

wsd.diagram(text, "earth", "png", function (err, buf, typ) {
    if (err) {
        console.error(err);
    } else {
        console.log("Received MIME type:", typ);
        fs.writeFile("logs.png", buf);
    }
});