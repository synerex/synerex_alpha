import React, { Component } from 'react';
import Header from './component/Header';
import SideBar from './component/SideBar';
import Content from './component/Content';
import ProcessLog from './component/ProcessLog';
import SelectContent from './component/SelectContent';
import MsgStore from './model/MsgStore';
import io from 'socket.io-client';


class App extends Component {

    constructor(props) {
        super(props);
        const socket = io();
        this.socket = socket;
        this.mstore = new MsgStore();
        this.state = {
            logs:[],
            providers:[],
            store: this.mstore,
            socket: socket,
            turn: true
        }
        socket.on('connect', () => { console.log("Socket.IO Connected!") });
        socket.on('event', this.getEvent.bind(this));
        socket.on('log', this.getLog.bind(this));
        socket.on('disconnect', () => { console.log("Socket.IO Disconnected!") });
        socket.on('providers', this.setProviders.bind(this));
        socket.on('nodeserv', this.getLog.bind(this));
        this.selComp =ProcessLog;
        this.selArg =  {logs:this.state.logs,
            store:this.state.store,
            turn: this.state.turn
        };
    }

    setProviders(data){
        console.log("Set Providers",data);
//        let providers = JSON.parse(data);
//        console.log("Parsed Data",providers);
        this.setState({
            providers:data
        });

    }

    addLog(message, value){
        let lg = this.state.logs;
        const now = new Date();
        // lg.push({ m: message, v: value, a: message_array, t: now.toLocaleString() });
        lg.push({
            message: message,
            value: value,
            time: now.toLocaleString()
        })
        this.setState({
            logs:lg,
        });
    }

    getEvent(data){
//        console.log("GetEvent:",data,":");
        // Parse Message

        this.addLog(data,"")
        this.mstore.addMsg(data)
        this.setState({
            store:this.mstore
        })
    }

    // we have to fix which event...
    getLog(data){
        console.log("GetLog:",data,":");
        let lg = this.state.logs;
        lg.push({
            message: data,
            value: "",
            time: ""
        })
        this.setState({
            logs:lg,
        });
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

    startNodeServ() {
        const sock = this.socket;
        console.log("start NodeServ2!")
        console.log(this.socket)
//        this.socket.emit("run","nodeserv", (data) =>{
//            console.log(data);
//        });
        this.socket.emit("run", "nodeserv");
    }

    runTarget(target){
        const sock = this.socket;
        console.log("run ",target)
        this.socket.emit("run", target);
    }



    render() {
        const header =  <Header />;
        const content = <SelectContent component={this.selComp}  args={this.selArg}/>;
        const sidebar = <SideBar clearLogs={() => this.clearLogs()}
                                 addLog={()=>this.addMsg()}
                                 viewChange={()=>this.viewChange()}
                                 toggleTurn={()=>this.toggleTurn()}
                                 startNodeServ={()=>this.startNodeServ()}
                                 start={(target)=>this.runTarget(target)}
                                 providers ={this.state.providers}
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
