import React, { useState, useEffect } from "react";
import logo from "./logo.svg";
import "./App.css";
import io from "socket.io-client";
import { Header, Sidebar, Content } from "./components";
import {
    Provider,
    Command,
    Log,
    CommandType,
    ProviderType,
    RunningInfo,
    Option,
    AgentType
} from "./types";

const mockProviders: Provider[] = [
    new Provider("1", ProviderType.PEDESTRIAN),
    new Provider("2", ProviderType.PEDESTRIAN)
];

const mockCommands: Command[] = [
    {
        Type: CommandType.SET_AREA,
        Name: "Set Area",
        Option: {
            AreaCoord: []
        }
    },
    {
        Type: CommandType.SET_CLOCK,
        Name: "Set Clock",
        Option: {
            Time: ""
        }
    },
    {
        Type: CommandType.SET_AGENTS,
        Name: "Set Agents",
        Option: {
            Type: AgentType.PEDESTRIAN,
            Num: 0
        }
    },
    {
        Type: CommandType.START_CLOCK,
        Name: "Start Clock",
        Option: {}
    },
    {
        Type: CommandType.STOP_CLOCK,
        Name: "Stop Clock",
        Option: {}
    }
];

const socket: SocketIOClient.Socket = io();
const App: React.FC = () => {
    //const [logs, setLogs] = useState<Log[]>([]);
    //const [runningInfos, setRunningInfos] = useState<RunningInfo[]>([]);
    const [providers, setProviders] = useState<Provider[]>(mockProviders);
    const [commands, setCommands] = useState<Command[]>(mockCommands);
    useEffect(() => {
        socket.on("connect", () => {
            console.log("Socket.IO Connected!");
        });
        socket.on("log", (jsonArray: string) => addLog(jsonArray));
        socket.on("running", (jsonArray: string) => addRunnningInfo(jsonArray));
        socket.on("providers", (jsonStr: string[]) => addProvider(jsonStr));
        //socket.on("commands", (jsonArray: string[]) => getCommands(jsonArray));
        socket.on("disconnect", () => {
            console.log("Socket.IO Disconnected!");
        });
    }, []);

    const getProvider = (id: Provider["ID"]) => {
        providers.forEach((provider: Provider) => {
            if (provider.ID === id) {
                return provider;
            }
        });
        return null;
    };

    // logがdaemonから送られる
    const addLog = (jsonStr: string) => {
        let newProviders: Provider[] = [...providers];
        const id: Provider["ID"] = "0";
        const log: Log = {
            ID: id,
            Type: ProviderType.NODEID_SERVER,
            Description: jsonStr
        };
        providers.forEach((provider: Provider, index: number) => {
            if (provider.ID === id) {
                provider.Logs.push(log);
                newProviders.splice(index, 1, provider);
            }
        });
        setProviders(newProviders);
    };

    // 1サイクル毎にプロバイダ毎の所要時間などの情報がdaemonから送られる
    const addRunnningInfo = (jsonStr: string) => {
        let newProviders: Provider[] = [...providers];
        const id: Provider["ID"] = "0";
        const runningInfo: RunningInfo = {
            Duration: 0,
            AgentsNum: 0
        };
        providers.forEach((provider: Provider, index: number) => {
            if (provider.ID === id) {
                provider.RunningInfos.push(runningInfo);
                newProviders.splice(index, 1, provider);
            }
        });
        setProviders(newProviders);
    };

    // providerに変化があった場合にdaemonから送られる
    const addProvider = (jsonArray: string[]) => {
        console.log("Get Providers!", jsonArray);
        let newProviders: Provider[] = [...providers];
        jsonArray.forEach((pName: string) => {
            const id: Provider["ID"] = "0";
            const type: ProviderType = ProviderType.PEDESTRIAN;
            var hasID: boolean = false;
            providers.forEach((provider: Provider) => {
                if (provider.ID === id) {
                    hasID = true;
                    newProviders.push(provider);
                }
            });
            if (!hasID) {
                newProviders.push(new Provider(id, type));
            }
        });

        setProviders(newProviders);
    };

    /*const getCommands = (jsonArray: string[]) => {
        console.log("Get Commands!", jsonArray);
        let commands: Command[] = [];
        jsonArray.forEach((pName: string) => {
            commands.push({
                Name: pName,
                Options: []
            });
        });
        setCommands(commands);
    };*/

    /*const runProvider = (provider: Provider) => {
        console.log("Click Provider!", provider);
        socket.emit("run", provider);
    };*/

    const runCommand = (command: Command) => {
        console.log("Click Command!", command);
        socket.emit("command", command);
    };

    return (
        <div className="App">
            <Header />
            <Sidebar
                providers={providers}
                commands={commands}
                runCommand={runCommand}
            />
            <Content providers={providers} />
        </div>
    );
};

export default App;
