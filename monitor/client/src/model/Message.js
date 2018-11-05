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

    extractBitsFrom (n, begin, end) {
        return (n % Math.pow(2, end) - n % Math.pow(2, begin)) / Math.pow(2, begin);
    };

    //;
    snowflakeToNodeID(sid){
        return this.extractBitsFrom(sid,12,22)
    }

    getSrcNodeID(){
        let s = this.snowflakeToNodeID(this.obj.src);
//        console.log("sid",this.obj.src, s);
        return s
    }

    getDstNodeID(){
        let s = this.snowflakeToNodeID(this.obj.dst);
        return s;
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