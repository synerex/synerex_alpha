import React, { Component } from 'react';
import Header from './component/Header';
import SideBar from './component/SideBar';
import Content from './component/Content';
import SelectContent from './component/SelectContent';
import WorldView from './component/WorldView';
import MsgStore from './model/MsgStore';
import io from 'socket.io-client';


class App extends Component {

    constructor(props) {
        super(props);
        const socket = io();
        this.mstore = new MsgStore();
        this.state = {
            logs:[],
            store: this.mstore,
            socket: socket,
            turn: true
        }
        socket.on('connect', () => { console.log("Socket.IO Connected!") });
        socket.on('event', this.getEvent.bind(this));
        socket.on('disconnect', () => { console.log("Socket.IO Disconnected!") });

//        this.selComp =Content;
        this.selComp =WorldView;
        this.selArg =  {logs:this.state.logs,
            store:this.state.store,
            turn: this.state.turn
        };
    }

    addLog(message, value){
        let lg = this.state.logs;
        const obj = JSON.parse(message);
        const now = new Date();
        if (lg.length === 999) console.log("Log length is now maximum 1000");
        if(lg.length > 999 ) lg.shift(); // if log array length larger than 1000, then remove head.
        lg.push({
            msgType: obj.msgType,
            chType: obj.chType,
            src: obj.src,
            dst: obj.dst,
            arg: obj.arg,
            value: value,
            time: now.toLocaleString()
        })

//        console.log("lg:", lg);
        this.setState({
            logs:lg,
        });
    }

    getEvent(data){
       console.log("GetEvent:", data);
        // Parse Message

        this.addLog(data,"")
        this.mstore.addMsg(data)
        this.setState({
            store:this.mstore
        })
    }

    clearLogs(){
      console.log("App:Clearlog");
        this.selArg =  {logs:[]};
        this.mstore.clear();
        this.setState({
            logs:this.selArg.logs,
            store: this.mstore
        });
    }

    // change Views
    viewChange(){
        if (this.selComp !== WorldView) {

            this.selComp = WorldView;
            this.selArg = {logs: this.state.logs, store:this.state.store};
            this.setState({
                logs:this.state.logs,
                turn:this.state.turn,
                store:this.state.store
            });
        }else{
            this.selComp = Content;
            this.selArg = {logs: this.state.logs};
            this.setState({
                logs:this.state.logs
            });
        }
    }

    addMsg(){
        let ix = Math.floor(Math.random()*8);
        let lg = '{"src":'+ix+'}'
        this.getEvent(lg)
        this.addLog(lg)
    }

    queryNodeName(nid,names){
        this.socket.emit("node",nid,
            function(data){
               console.log("Got! nodeInfo",data);
               names[nid] = data;
            });
    }

    addDemo() {
        this.addLog("Hello", "abc");
    }

    toggleTurn(){
        this.setState({
            turn:!this.state.turn
        });
        this.selArg.turn = this.state.turn;
    }

    componentDidMount(){
//        this.interval = setInterval(() => this.addDemo(), 5000);
    }


    render() {
      const header =  <Header />;
      const content = <SelectContent component={this.selComp}  args={this.selArg}/>;
      const sidebar = <SideBar clearLogs={() => this.clearLogs()}
                               addLog={()=>this.addMsg()}
                               viewChange={()=>this.viewChange()}
                               toggleTurn={()=>this.toggleTurn()}
                    />;
        return (
            <div>
                {header}
                {sidebar}
                {content}
            </div>
        );
    }
}

export default App;