import React, {useState} from 'react';
import logo from './logo.svg';
import './App.css';

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

const MAPBOX_TOKEN = process.env.MAPBOX_ACCESS_TOKEN; //Acquire Mapbox accesstoken



const App: React.FC = () => {

	[areaData, setAreaData] = useState()
	[geoData, setGeoData] = useState()

	const socket = io();
	socket.on("geojson", (data)=>getGeoJson(data));
    socket.on("areajson", (data)=>getAreaJson(data));
	socket.on("agents", (data)=>getAgents(data));
	
	getGeoJson(data) {
		const geoData = Json.parse(data)
        console.log("GeoData: ", geoData);
        setGeoData(geoData);
    }

    getAreaJson(data) {
		const areaData = Json.parse(data)
        console.log("AreaData: ", areaData);
        setAreaData(areaData);
	}
	
	getAgents(data){
		const { actions, movesbase, movedData } = this.props;
        const time = Date.now() / 1000; // set time as now. (If data have time, ..)
        const setMovesbase = [];
        //const setMovedData = [];
        const movesbasedata = [...movesbase];

        console.log("socketData length", data.length);
        //console.log("movesbasedata length", movesbasedata.length)

        data.forEach(value => {
            const { mtype, id, lat, lon, angle, speed, area } = JSON.parse(
                value
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
                        radius: 20,
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


	var layers = [];

  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          Edit <code>src/App.tsx</code> and save to reload.
        </p>
        <a
          className="App-link"
          href="https://reactjs.org"
          target="_blank"
          rel="noopener noreferrer"
        >
          Learn React
        </a>
      </header>
    </div>
  );
}

export default connectToHarmowareVis(App);
