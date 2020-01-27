import React from "react";
import {
    Container,
    connectToHarmowareVis,
    HarmoVisLayers,
    MovesLayer,
    MovesInput,
    LoadingIcon,
    FpsDisplay,
    DepotsLayer,
    EventInfo,
    MovesbaseOperation,
    MovesBase,
    BasedProps
} from "harmoware-vis";

//import { StaticMap,  } from 'react-map-gl';
import { Layer } from "@deck.gl/core";
import DeckGL from "@deck.gl/react";
import { GeoJsonLayer, LineLayer } from "@deck.gl/layers";

import {
    _MapContext as MapContext,
    InteractiveMap,
    NavigationControl
} from "react-map-gl";

import Controller from "../components/controller";

import * as io from "socket.io-client";

const MAPBOX_TOKEN = process.env.MAPBOX_ACCESS_TOKEN; //Acquire Mapbox accesstoken

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
            mapbox_token: MAPBOX_TOKEN,
            geojson: null,
            areajson: null,
            lines: [
                {
                    sourcePosition: [136.97285, 35.159431, 0],
                    targetPosition: [136.97705, 35.159431, 0]
                }
            ],
            linecolor: [0, 255, 255],
            popup: [0, 0, ""]
        };

        // for receiving event info.
        socket.on("connect", () => {
            console.log("Socket.IO connected2!");
        });
        socket.on("event", this.getEvent.bind(this));
        socket.on("geojson", this.getGeoJson.bind(this));
        socket.on("areajson", this.getAreaJson.bind(this));
        socket.on("agents", this.getAgents.bind(this));

        /*socket.on('mapbox_token', (token) => {
			console.log("Token Got:" + MAPBOX_TOKEN)
			this.setState({ mapbox_token: token })
		});*/

        socket.on("disconnect", () => {
            console.log("Socket.IO disconnected!");
        });
    }

    getGeoJson(data) {
        console.log("Geojson:" + data.length);
        console.log("jsonData", data);
        console.log(JSON.parse(data));
        this.setState({ geojson: JSON.parse(data) });
    }

    getAreaJson(data) {
        console.log("Areajson:" + data.length);
        console.log("jsonData", data);
        console.log(JSON.parse(data));
        this.setState({ areajson: JSON.parse(data) });
    }

    getLines(data) {
        console.log("getLines!:" + data.length);
        //		console.log(data)
        if (this.state.lines.length > 0) {
            const ladd = JSON.parse(data);
            const lbase = this.state.lines;
            const lists = lbase.concat(ladd);
            this.setState({ lines: lists });
        } else {
            this.setState({ lines: JSON.parse(data) });
        }
    }

    getAgents(data) {
        const { actions, movesbase } = this.props;
        const agents = JSON.parse(data).agents;
        //		console.log(data)
        //		console.log(agents)

        const time = Date.now() / 1000; // set time as now. (If data have time, ..)
        //		let hit = false;
        //		const movesbasedata = [...movesbase]; // why copy !?
        let setMovesbase = [];

        if (movesbase.length == 0) {
            //			console.log("Initial!:" + agents.length)
            for (let i = 0, len = agents.length; i < len; i++) {
                setMovesbase.push({
                    mtype: 0,
                    id: i,
                    departuretime: time,
                    arrivaltime: time,
                    operation: [
                        {
                            elapsedtime: time,
                            position: [
                                agents[i].point[0],
                                agents[i].point[1],
                                0
                            ],
                            angle: 0,
                            speed: 1
                        }
                    ]
                });
            }
            // we may refresh viewport
        } else {
            //			console.log("Aget Update!" + data.length+":"+ agents[0])
            for (let i = 0, lengthi = movesbase.length; i < lengthi; i++) {
                movesbase[i].arrivaltime = time;
                movesbase[i].operation.push({
                    elapsedtime: time,
                    position: [agents[i].point[0], agents[i].point[1], 0],
                    angle: 0,
                    speed: 1
                });
                //				setMovesbase.push(movesbase[i]);
            }
            setMovesbase = movesbase;
        }

        actions.updateMovesBase(setMovesbase);
    }

    getEvent(socketsData) {
        const { actions, movesbase, movedData } = this.props;
        const time = Date.now() / 1000; // set time as now. (If data have time, ..)
        const setMovesbase = [];
        //const setMovedData = [];
        const movesbasedata = [...movesbase];

        console.log("socketData length", socketsData.length);
        //console.log("movesbasedata length", movesbasedata.length)

        socketsData.forEach(socketData => {
            const { mtype, id, lat, lon, angle, speed, area } = JSON.parse(
                socketData
            );

            let hit = false;
            movesbasedata.forEach(movedata => {
                if (mtype === movedata.mtype && id === movedata.id) {
                    let color = [0, 255, 0];
                    /*if (mtype == 0) {
                        // Ped
                        color = [0, 200, 120];
                    } else if (mtype == 1) {
                        // Car
                        color = [200, 0, 0];
                    }*/
                    hit = true;
                    movedata.arrivaltime = time;
                    movedata.operation.push({
                        elapsedtime: time,
                        position: [lon, lat, 0],
                        radius: 20,
                        angle,
                        speed,
                        color
                    });

                    setMovesbase.push(movedata);
                }
            });

            if (!hit) {
                let color = [0, 255, 0];
                /*if (mtype == 0) {
                    // Ped
                    color = [0, 200, 120];
                } else if (mtype == 1) {
                    // Car
                    color = [200, 0, 0];
                }*/
                setMovesbase.push({
                    mtype,
                    id,
                    departuretime: time,
                    arrivaltime: time,
                    operation: [
                        {
                            elapsedtime: time,
                            position: [lon, lat, 0],
                            radius: 20,
                            angle,
                            speed,
                            color
                        }
                    ]
                });
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

        /*
		for (let i = 0, lengthi = movesbasedata.length; i < lengthi; i += 1) {
			const { departuretime: propsdeparturetime, operation: propsoperation } = movesbasedata[i];
			let departuretime = propsdeparturetime;
			let startIndex = propsoperation.length;
			for (let j = 0, lengthj = propsoperation.length; j < lengthj; j += 1) {
				if (propsoperation[j].elapsedtime > compareTime) {
					startIndex = j;
					departuretime = propsoperation[j].elapsedtime;
					break;
				}
			}
			if (startIndex === 0) {
				setMovesbase.push(Object.assign({}, movesbasedata[i]));
			} else
				if (startIndex < propsoperation.length) {
					setMovesbase.push(Object.assign({}, movesbasedata[i], {
						operation: propsoperation.slice(startIndex), departuretime
					}));
					dataModify = true;
				} else {
					dataModify = true;
				}
		}*/
        if (!animatePause) {
            actions.setAnimatePause(true);
        }
        actions.updateMovesBase(setMovesbase);
        if (!animatePause) {
            actions.setAnimatePause(false);
        }
        console.log(this.props.viewport);
        console.log(MapContext.viewport);
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

    initialize(gl) {
        gl.enable(gl.DEPTH_TEST);
        gl.depthFunc(gl.LEQUAL);
        console.log("GL Initialized!");
    }

    logViewPort(state, view) {
        console.log("Viewport changed!", state, view);
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
            movedData,
            mapStyle
        } = props;
        //	const { movesFileName } = inputFileName;
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
        var layers = [];

        if (this.state.geojson != null) {
            console.log("push layer geojson");
            layers.push(
                new GeoJsonLayer({
                    id: "geojson-layer",
                    data: this.state.geojson,
                    pickable: true,
                    stroked: false,
                    filled: true,
                    extruded: true,
                    lineWidthScale: 2,
                    lineWidthMinPixels: 2,
                    getFillColor: [160, 160, 180, 200],
                    //				getLineColor: d => colorToRGBArray(d.properties.color),
                    getLineColor: [255, 255, 255],
                    getRadius: 1,
                    getLineWidth: 1,
                    getElevation: 10
                    //				onHover: ({object, x, y}) => {
                    //				  const tooltip = object.properties.name || object.properties.station;
                    //				}
                })
            );
        }

        if (this.state.areajson != null) {
            console.log("push layer areajson");
            this.state.areajson.forEach(value => {
                value.duplicate_area.forEach((area, index) => {
                    let sourcePosition = [];
                    let targetPosition = [];
                    if (index != value.duplicate_area.length - 1) {
                        sourcePosition = [area.longitude, area.latitude];
                        targetPosition = [
                            value.duplicate_area[index + 1].longitude,
                            value.duplicate_area[index + 1].latitude
                        ];
                    } else {
                        sourcePosition = [area.longitude, area.latitude];
                        targetPosition = [
                            value.duplicate_area[0].longitude,
                            value.duplicate_area[0].latitude
                        ];
                    }
                    layers.push(
                        new LineLayer({
                            visible: true,
                            data: value.duplicate_area,
                            getSourcePosition: d => sourcePosition,
                            getTargetPosition: d => targetPosition,
                            getColor: this.state.linecolor,
                            getWidth: 1,
                            widthMinPixels: 0.1
                        })
                    );
                });

                value.control_area.forEach((area, index) => {
                    let sourcePosition = [];
                    let targetPosition = [];
                    if (index != value.control_area.length - 1) {
                        sourcePosition = [area.longitude, area.latitude];
                        targetPosition = [
                            value.control_area[index + 1].longitude,
                            value.control_area[index + 1].latitude
                        ];
                    } else {
                        sourcePosition = [area.longitude, area.latitude];
                        targetPosition = [
                            value.control_area[0].longitude,
                            value.control_area[0].latitude
                        ];
                    }
                    layers.push(
                        new LineLayer({
                            visible: true,
                            data: value.control_area,
                            getSourcePosition: d => sourcePosition,
                            getTargetPosition: d => targetPosition,
                            getColor: this.state.linecolor,
                            getWidth: 1,
                            widthMinPixels: 0.1
                        })
                    );
                });
            });

            /*controlAreaLine = layers.push(
                new GeoJsonLayer({
                    id: "areajson-layer",
                    data: this.state.areajson,
                    pickable: true,
                    stroked: false,
                    filled: true,
                    extruded: true,
                    lineWidthScale: 2,
                    lineWidthMinPixels: 2,
                    getFillColor: [160, 160, 180, 200],
                    //				getLineColor: d => colorToRGBArray(d.properties.color),
                    getLineColor: [255, 255, 255],
                    getRadius: 1,
                    getLineWidth: 1,
                    getElevation: 10
                    //				onHover: ({object, x, y}) => {
                    //				  const tooltip = object.properties.name || object.properties.station;
                    //				}
                })
            );*/
        }

        if (this.state.moveDataVisible && movedData.length > 0) {
            layers.push(
                new MovesLayer({
                    viewport,
                    routePaths,
                    movesbase,
                    movedData,
                    clickedObject,
                    actions,
                    lightSettings,
                    visible: this.state.moveDataVisible,
                    optionVisible: this.state.moveOptionVisible,
                    layerRadiusScale: 0.01,
                    getRaduis: x => 0.02,
                    getStrokeWidth: 0.01,
                    optionCellSize: 2,
                    sizeScale: 1,
                    iconChange: false,
                    optionChange: false, // this.state.optionChange,
                    onHover
                })
            );
        }

        const onViewportChange =
            this.props.onViewportChange || actions.setViewport;
        //		viewState={viewport}

        const visLayer =
            this.state.mapbox_token.length > 0 ? (
                <DeckGL
                    layers={layers}
                    onWebGLInitialized={this.initialize}
                    initialViewState={{
                        longitude: 136.974572,
                        latitude: 35.158625,
                        zoom: 17
                    }}
                    controller={true}
                    ContextProvider={MapContext.Provider}
                >
                    <InteractiveMap
                        viewport={viewport}
                        mapStyle={"mapbox://styles/mapbox/dark-v8"}
                        onViewportChange={onViewportChange}
                        mapboxApiAccessToken={this.state.mapbox_token}
                        visible={true}
                    ></InteractiveMap>
                </DeckGL>
            ) : (
                <LoadingIcon loading={true} />
            );

        /*					<div style={{ position: "absolute", left: 30, top: 120, zIndex: 1 }}>
						<NavigationControl />
					</div>
				*/

        /*				<InteractiveMap
					viewport={viewport} 
					mapStyle={'mapbox://styles/mapbox/dark-v8'}
					onViewportChange={onViewportChange}
					mapboxApiAccessToken={this.state.mapbox_token}
					visible={true}>

					<DeckGL viewState={viewport} layers={layers} onWebGLInitialized={this.initialize} />

				</InteractiveMap>
				: <LoadingIcon loading={true} />;
*/
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
                <div className="harmovis_area">{visLayer}</div>
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

                <FpsDisplay />
            </div>
        );
    }
}
export default connectToHarmowareVis(App);
