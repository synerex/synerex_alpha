import React from "react";
import {
    Container,
    connectToHarmowareVis,
    HarmoVisLayers,
    MovesLayer,
    LineMapLayer,
    MovesInput,
    LoadingIcon,
    FpsDisplay,
    DepotsLayer,
    EventInfo,
    MovesbaseOperation,
    MovesBase,
    BasedProps
} from "harmoware-vis";

import Controller from "../components/controller";

import * as io from "socket.io-client";

const MAPBOX_TOKEN = process.env.MAPBOX_ACCESS_TOKEN; //Acquire Mapbox accesstoken
console.log("mapbox_token: ", MAPBOX_TOKEN);

class App extends Container {
    constructor(props) {
        super(props);
        const { setSecPerHour, setLeading, setTrailing } = props.actions;
        setSecPerHour(3600);
        setLeading(3);
        setTrailing(3);
        const socket = io();
        this.state = {
            moveDataVisible: true,
            moveOptionVisible: false,
            depotOptionVisible: false,
            heatmapVisible: false,
            optionChange: false,
            popup: [0, 0, ""],
            linemapData: [
                // area
                {
                    sourcePosition: [136.973172, 35.152476, 0],
                    targetPosition: [136.984031, 35.152476, 0],
                    strokeWidth: 6.0
                },
                {
                    sourcePosition: [136.973172, 35.160678, 0],
                    targetPosition: [136.984031, 35.160678, 0],
                    strokeWidth: 6.0
                },
                {
                    sourcePosition: [136.973172, 35.152476, 0],
                    targetPosition: [136.973172, 35.160678, 0],
                    strokeWidth: 6.0
                },
                {
                    sourcePosition: [136.984031, 35.152476, 0],
                    targetPosition: [136.984031, 35.160678, 0],
                    strokeWidth: 6.0
                },

                {
                    sourcePosition: [136.981014, 35.152476, 0],
                    targetPosition: [136.990047, 35.152476, 0],
                    strokeWidth: 6.0
                },
                {
                    sourcePosition: [136.981014, 35.160678, 0],
                    targetPosition: [136.990047, 35.160678, 0],
                    strokeWidth: 6.0
                },
                {
                    sourcePosition: [136.981014, 35.152476, 0],
                    targetPosition: [136.981014, 35.160678, 0],
                    strokeWidth: 6.0
                },
                {
                    sourcePosition: [136.990047, 35.152476, 0],
                    targetPosition: [136.990047, 35.160678, 0],
                    strokeWidth: 6.0
                },
                // controlled
                {
                    sourcePosition: [136.9825, 35.152476, 0],
                    targetPosition: [136.9825, 35.160678, 0],
                    color: [255, 0, 255],
                    strokeWidth: 6.0
                }

                // wall1　下
                /*{
                    sourcePosition: [136.9823, 35.155078, 0],
                    targetPosition: [136.9823, 35.152476, 0],
                    color: [55, 100, 200],
                    strokeWidth: 6.0
                },
                {
                    sourcePosition: [136.9827, 35.155078, 0],
                    targetPosition: [136.9827, 35.152476, 0],
                    color: [55, 100, 200],
                    strokeWidth: 6.0
                },
                {
                    sourcePosition: [136.9823, 35.155078, 0],
                    targetPosition: [136.9827, 35.155078, 0],
                    color: [55, 100, 200],
                    strokeWidth: 6.0
                },

                // wall2 上
                {
                    sourcePosition: [136.9823, 35.157576, 0],
                    targetPosition: [136.9823, 35.160678, 0],
                    color: [55, 100, 200],
                    strokeWidth: 6.0
                },
                {
                    sourcePosition: [136.9827, 35.157576, 0],
                    targetPosition: [136.9827, 35.160678, 0],
                    color: [55, 100, 200],
                    strokeWidth: 6.0
                },
                {
                    sourcePosition: [136.9823, 35.157576, 0],
                    targetPosition: [136.9827, 35.157576, 0],
                    color: [55, 100, 200],
                    strokeWidth: 6.0
                }*/
            ]
        };

        // for receiving event info.
        socket.on("connect", () => {
            console.log("Socket.IO connected!");
        });
        socket.on("event", this.getEvent.bind(this));
        socket.on("area", this.getArea.bind(this));
        socket.on("disconnect", () => {
            console.log("Socket.IO disconnected!");
        });
    }

    getArea(socketData) {
        socketData.forEach(areaData => {
            console.log("areaData: ", areaData);
        });
    }

    getEvent(socketsData) {
        //console.log("Get event4 socketsData!!", socketsData)
        const { actions, movesbase, movedData } = this.props;
        const time = Date.now() / 1000; // set time as now. (If data have time, ..)
        const setMovesbase = [];
        //const setMovedData = [];
        const movesbasedata = [...movesbase]; // why copy !?
        //const movedData2 = [...movedData]; // why copy !?

        console.log("socketData length", socketsData.length);
        //console.log("movesbasedata length", movesbasedata.length)

        socketsData.forEach(socketData => {
            const { mtype, id, lat, lon, angle, speed, area } = JSON.parse(
                socketData
            );

            let hit = false;
            movesbasedata.forEach(movedata => {
                if (mtype === movedata.mtype && id === movedata.id) {
                    let color = [0, 200, 0];
                    if (mtype == 0) {
                        // Ped
                        color = [0, 200, 120];
                    } else if (mtype == 1) {
                        // Car
                        color = [200, 0, 0];
                    }
                    hit = true;
                    movedata.arrivaltime = time;
                    movedata.operation.push({
                        elapsedtime: time,
                        position: [lon, lat, 0],
                        radius: 1,
                        angle,
                        speed,
                        color
                    });

                    setMovesbase.push(movedata);
                }
            });

            if (!hit) {
                let color = [0, 200, 0];
                if (mtype == 0) {
                    // Ped
                    color = [0, 200, 120];
                } else if (mtype == 1) {
                    // Car
                    color = [200, 0, 0];
                }
                setMovesbase.push({
                    mtype,
                    id,
                    departuretime: time,
                    arrivaltime: time,
                    operation: [
                        {
                            elapsedtime: time,
                            position: [lon, lat, 0],
                            radius: 1,
                            angle,
                            speed,
                            color
                        }
                    ]
                });
                /*setMovedData.push({
				sourceColor: [255, 0, 255]
	    });*/
            }
        });

        console.log("lenth before", setMovesbase.length);
        actions.updateMovesBase(setMovesbase);
        //actions.updateMovedData(setMovedData);
    }

    deleteMovebase(maxKeepSecond) {
        const { actions, animatePause, movesbase, settime } = this.props;
        const movesbasedata = [...movesbase];
        const setMovesbase = [];
        let dataModify = false;
        const compareTime = settime - maxKeepSecond;
        for (let i = 0, lengthi = movesbasedata.length; i < lengthi; i += 1) {
            const {
                departuretime: propsdeparturetime,
                operation: propsoperation
            } = movesbasedata[i];
            let departuretime = propsdeparturetime;
            let startIndex = propsoperation.length;
            for (
                let j = 0, lengthj = propsoperation.length;
                j < lengthj;
                j += 1
            ) {
                if (propsoperation[j].elapsedtime > compareTime) {
                    startIndex = j;
                    departuretime = propsoperation[j].elapsedtime;
                    break;
                }
            }
            if (startIndex === 0) {
                setMovesbase.push(Object.assign({}, movesbasedata[i]));
            } else if (startIndex < propsoperation.length) {
                setMovesbase.push(
                    Object.assign({}, movesbasedata[i], {
                        operation: propsoperation.slice(startIndex),
                        departuretime
                    })
                );
                dataModify = true;
            } else {
                dataModify = true;
            }
        }
        if (dataModify) {
            if (!animatePause) {
                actions.setAnimatePause(true);
            }
            actions.updateMovesBase(setMovesbase);
            if (!animatePause) {
                actions.setAnimatePause(false);
            }
        }
    }

    getMoveDataChecked(e) {
        this.setState({ moveDataVisible: e.target.checked });
    }

    getMoveOptionChecked(e) {
        this.setState({ moveOptionVisible: e.target.checked });
    }

    getDepotOptionChecked(e) {
        this.setState({ depotOptionVisible: e.target.checked });
    }

    getOptionChangeChecked(e) {
        this.setState({ optionChange: e.target.checked });
    }

    render() {
        const props = this.props;
        const {
            actions,
            clickedObject,
            inputFileName,
            viewport,
            deoptsData,
            loading,
            routePaths,
            lightSettings,
            movesbase,
            movedData
        } = props;
        //console.log("viewPort: ", viewport);
        /*var pedMovesbase = []
	var pedMovedData = []
	var carMovesbase = []
	var carMovedData = []
	var movedDataInfo = [...movedData]
	movesbase.forEach((movebase, index)=>{
		let mtype = movebase.mtype
		if (mtype === 0){	// ped
			pedMovesbase.push(movebase)
			pedMovedData.push(movedDataInfo[index])
		}else if(mtype === 1){
			carMovesbase.push(movebase)
			carMovedData.push(movedDataInfo[index])
		}
	})
	console.log("pedmobes", pedMovesbase)
		console.log("pedmovedData", pedMovedData)
		console.log("carmobes", carMovesbase)
		console.log("carmovedData", carMovedData)*/

        const optionVisible = false;
        const onHover = el => {
            if (el && el.object) {
                let disptext = "";
                const objctlist = Object.entries(el.object);
                for (
                    let i = 0, lengthi = objctlist.length;
                    i < lengthi;
                    i += 1
                ) {
                    const strvalue = objctlist[i][1].toString();
                    disptext += i > 0 ? "\n" : "";
                    disptext += `${objctlist[i][0]}: ${strvalue}`;
                }
                this.setState({ popup: [el.x, el.y, disptext] });
            } else {
                this.setState({ popup: [0, 0, ""] });
            }
        };

        return (
            <div>
                <Controller
                    {...props}
                    deleteMovebase={this.deleteMovebase.bind(this)}
                    getMoveDataChecked={this.getMoveDataChecked.bind(this)}
                    getMoveOptionChecked={this.getMoveOptionChecked.bind(this)}
                    getDepotOptionChecked={this.getDepotOptionChecked.bind(
                        this
                    )}
                    getOptionChangeChecked={this.getOptionChangeChecked.bind(
                        this
                    )}
                />
                <div className="harmovis_area">
                    <HarmoVisLayers
                        viewport={viewport}
                        actions={actions}
                        mapboxApiAccessToken={MAPBOX_TOKEN}
                        layers={
                            this.state.moveDataVisible && movedData.length > 0
                                ? [
                                      new LineMapLayer({
                                          viewport,
                                          linemapData: this.state.linemapData
                                      }),
                                      new MovesLayer({
                                          viewport,
                                          routePaths,
                                          movesbase,
                                          movedData,
                                          clickedObject,
                                          actions,
                                          lightSettings,
                                          visible: this.state.moveDataVisible,
                                          optionVisible: this.state
                                              .moveOptionVisible,
                                          optionChange: this.state.optionChange,
                                          iconChange: false,
                                          onHover
                                      })
                                  ]
                                : [
                                      new LineMapLayer({
                                          viewport,
                                          linemapData: this.state.linemapData
                                      })
                                  ]
                        }
                    />
                </div>
                <svg
                    width={viewport.width}
                    height={viewport.height}
                    className="harmovis_overlay"
                >
                    <g fill="white" fontSize="12">
                        {this.state.popup[2].length > 0
                            ? this.state.popup[2]
                                  .split("\n")
                                  .map((value, index) => (
                                      <text
                                          x={this.state.popup[0] + 10}
                                          y={this.state.popup[1] + index * 12}
                                          key={index.toString()}
                                      >
                                          {value}
                                      </text>
                                  ))
                            : null}
                    </g>
                </svg>
                <LoadingIcon loading={loading} />
                <FpsDisplay />
            </div>
        );
    }
}
export default connectToHarmowareVis(App);
