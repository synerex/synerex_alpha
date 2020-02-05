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
    AgentType,
    ContentType
} from "./types";

const mockProviders: Provider[] = [
    new Provider(1, ProviderType.PEDESTRIAN),
    new Provider(2, ProviderType.CAR)
];

const mockCommands: Command[] = [
    {
        Type: CommandType.SET_AREA,
        Name: "SetArea",
        Option: {
            AreaCoord: []
        }
    },
    {
        Type: CommandType.SET_CLOCK,
        Name: "SetClock",
        Option: {
            Time: ""
        }
    },
    {
        Type: CommandType.SET_AGENTS,
        Name: "SetAgents",
        Option: {
            Type: AgentType.PEDESTRIAN,
            Num: 0
        }
    },
    {
        Type: CommandType.START_CLOCK,
        Name: "StartClock",
        Option: {}
    },
    {
        Type: CommandType.STOP_CLOCK,
        Name: "StopClock",
        Option: {}
    }
];

const socket: SocketIOClient.Socket = io();
const App: React.FC = () => {
    const [providers, setProviders] = useState<Provider[]>([]);
    const [commands, setCommands] = useState<Command[]>(mockCommands);
    const [contentType, setContentType] = useState<ContentType>(
        ContentType.LOG_MONITOR
    );
    useEffect(() => {
        socket.on("connect", () => {
            console.log("Socket.IO Connected!");
        });
        socket.on("log", (jsonArray: string) => addLog(jsonArray));
        socket.on("running", (jsonArray: string) => addRunnningInfo(jsonArray));
        socket.on("providers", (jsonStr: string[]) => getProviders(jsonStr));
        //socket.on("commands", (jsonArray: string[]) => getCommands(jsonArray));
        socket.on("disconnect", () => {
            console.log("Socket.IO Disconnected!");
        });
    }, []);

    // logがdaemonから送られる
    const addLog = (jsonStr: string) => {
        setProviders(prevProviders => {
            let newProviders: Provider[] = [...prevProviders];
            //console.log("setLog11", prevProviders);
            //console.log("setLog1", jsonStr);
            const prelog = JSON.parse(jsonStr);
            //console.log("setLog2", prelog);
            const log: Log = {
                ID: prelog.ID,
                Description: prelog.Description
            };
            //console.log("Log: ", log);
            prevProviders.forEach((provider: Provider, index: number) => {
                //console.log("pro same: ", provider.ID, log.ID);
                if (provider.ID === log.ID) {
                    //console.log("same id");
                    provider.addLog(log);
                    newProviders.splice(index, 1, provider);
                }
            });
            //console.log("setLog", newProviders);
            return newProviders;
        });
    };

    // 1サイクル毎にプロバイダ毎の所要時間などの情報がdaemonから送られる
    const addRunnningInfo = (jsonStr: string) => {
        setProviders(prevProviders => {
            let newProviders: Provider[] = [...prevProviders];
            const id: Provider["ID"] = 0;
            const runningInfo: RunningInfo = {
                Duration: 0,
                AgentsNum: 0
            };
            prevProviders.forEach((provider: Provider, index: number) => {
                if (provider.ID === id) {
                    provider.RunningInfos.push(runningInfo);
                    newProviders.splice(index, 1, provider);
                }
            });
            return newProviders;
        });
        //console.log("setRunningInof");
        //setProviders(newProviders);
    };

    // providerに変化があった場合にdaemonから送られる
    const getProviders = (jsonArray: string[]) => {
        setProviders(prevProviders => {
            let newProviders: Provider[] = [];
            jsonArray.forEach((pjson: string) => {
                const provider = JSON.parse(pjson);
                let existProvider = null;
                let newProvider = null;
                prevProviders.forEach((prevProvider: Provider) => {
                    // 新規プロバイダの場合追加
                    if (provider.ID === prevProvider.ID) {
                        // すでに存在する場合、元のプロバイダを追加
                        existProvider = prevProvider;
                        /*newProvider = new Provider(
                            provider.ID,
                            checkProviderType(provider.Name)
                        );*/
                    }
                });
                if (existProvider) {
                    //console.log("exist!", existProvider);
                    newProviders.push(existProvider);
                } else {
                    console.log("new!", provider);
                    newProviders.push(
                        new Provider(
                            provider.ID,
                            checkProviderType(provider.Name)
                        )
                    );
                }
            });
            return newProviders;
        });
    };

    const checkProviderType = (name: string) => {
        switch (name) {
            case "NodeIDServer":
                return ProviderType.NODEID_SERVER;
            case "MonitorServer":
                return ProviderType.MONITOR_SERVER;
            case "SynerexServer":
                return ProviderType.SYNEREX_SERVER;
            case "Scenario":
                return ProviderType.SCENARIO;
            case "Pedestrian":
                return ProviderType.PEDESTRIAN;
            case "Visualization":
                return ProviderType.VISUALIZATION;
            case "Car":
                return ProviderType.CAR;
            case "Clock":
                return ProviderType.CLOCK;
            default:
                return ProviderType.PEDESTRIAN;
        }
    };

    // providerに変化があった場合にdaemonから送られる
    const addProvider = (jsonArray: string[]) => {
        console.log("Get Providers!", jsonArray);
        let newProviders: Provider[] = [...providers];
        jsonArray.forEach((pName: string) => {
            const id: Provider["ID"] = 0;
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
        socket.emit("order", command);
        const err = true;
        return err;
    };

    const changeContent = (type: ContentType) => {
        setContentType(type);
    };

    return (
        <div className="App">
            <Header changeContent={changeContent} />
            <Sidebar
                providers={providers}
                commands={commands}
                runCommand={runCommand}
            />
            <Content providers={providers} contentType={contentType} />
        </div>
    );
};

export default App;
