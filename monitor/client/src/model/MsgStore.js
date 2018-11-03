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
        this.store = [];
        this.nodes = {};
        this.nodeList = [];
    }

    getNodeNum(){
        const size =  Object.keys(this.nodes).length;
        if (size < 8)  return 8;
        return size;
    }

    addNode(nodeID){
        if( this.nodes[nodeID] === undefined){
            // there is no nodeID info
            const nd = new Node(nodeID);
            const n = this.nodeList.length;
            this.nodeList.push(nd);
            this.nodes[nodeID] = {node:nd, idx:n};
        }
    }

    // message store also should have the maximum number of messages
    addMsg(mes){ // get JSON string
        const ms = new Message(mes)
        this.addNode(ms.getSrcNodeID())
        if (this.store.length > MAX_NUM){
            this.store.shift(); // remove top data
        }
//        this.addNode(ms.getDestNodeID())
        this.store.push(ms)
    }

    getMsgCount(){
        return this.store.length;
    }

    getMsg(i){
        return this.store[i];
    }

    getNodeIndex(nodeID){
        return this.nodes[nodeID].idx;
    }
}

