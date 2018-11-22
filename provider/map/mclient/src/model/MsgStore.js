// MsgStore.js

import Message from './Message';
import Node from './Node';

const MAX_NUM=10000;

export default class MsgStore {
    constructor(){
        console.log("New Message Store");
        this.clear();
    }

    clear(){
        this.bus = {};
        this.taxi = {};
        this.busTrace = {};
    }

    addVehicle(store, ms){
        if( store[ms.id] === undefined) store[ms.id]=[]
        if(store[ms.id].length > 1000) store[ms.id].pop()
        store[ms.id].unshift([ms.lat, ms.lon,ms.angle, ms.speed])
    }


    // message store also should have the maximum number of messages
    addPosition(mes){ // get JSON string
        const ms =  JSON.parse(mes);
        if( ms.mtype == 1){
            this.addVehicle(this.taxi,ms);
        }else{
            this.addVehicle(this.bus,ms);
        }
    }

    getVehicle(mtype){
        if (mtype===1){
            return this.taxi
        }else{
            return this.bus
        }
    }

}

