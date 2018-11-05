// Message.js

// Model for Messages

// Provider -> Register
// Provider <- Server


// Msg Json format
// {
//    src: srcNodeID
//    dst: destNodeID
//
// }

export default class Message {


    constructor(json){
        this.js = json;
//        console.log("Parse:"+json);
        this.obj = JSON.parse(json);
    }

    getSrcNodeID(){
        return this.obj.src;
    }

    getDstNodeID(){
        return this.obj.dst;
    }

    getMsgType(){
        return this.obj.msgType;
    }

    getChType(){
        return this.obj.chType;
    }

    getArgs(){
        return this.obj.arg;
    }
}