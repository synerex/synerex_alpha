// MsgStore.js

import Message from './Message';
import Node from './Node';

const MAX_NUM=10000;

export default class MsgStore {
    constructor(qfunc){
        console.log("New Message Store");
        this.queryFunc = qfunc;
        this.clear();
    }

    clear() {
        this.store = [];
        this.nodes = {};
        this.nodeList = [];
    }


    setNodeUpdateFunc(fn){
        this.nodeUpdateFunc = fn;
    }

    getNodeNum(){
        const size =  Object.keys(this.nodes).length;
        if (size < 6)  return 6; //Todo: fix: just limit node number to 6 for demo  20181105
        return size;
    }
    // called from App.js socket.io
    setNodeName(nodeID, name){
        this.nodes[nodeID].name = name;
        //
        console.log("Got Name for Visual:",nodeID,name);
        this.nodeUpdateFunc(this.getNodeIndex(nodeID), name);
    }
    addNode(nodeID){
//        console.log("AddNode",nodeID);
        if( this.nodes[nodeID] === undefined){
            // there is no nodeID info
            const nd = new Node(nodeID);
//            console.log("Query nodename:",nodeID);
            this.queryFunc(nodeID,this);
            const n = this.nodeList.length; // Todo: fix to adapt node removal.
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
//        if(ms.getDstNodeID() != 0) {// may not appear unknown nodes for destination.
//            this.addNode(ms.getDstNodeID());
//        }
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

