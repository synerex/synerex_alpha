import React, { Component } from 'react';
import L from 'leaflet';
import {Map, Marker, Popup, TileLayer, Polyline} from 'react-leaflet';
import RMarker from './RMarker';

var smallVehicleIcon = L.icon({
    iconUrl: "img/car_above.png", iconSize: [16, 32], iconAnchor: [8, 16]
});
var midVehicleIcon = L.icon({
    iconUrl: "img/car_above.png", iconSize: [32, 64], iconAnchor: [16, 32]
});
var largeVehicleIcon = L.icon({
    iconUrl: "img/car_above.png", iconSize: [64, 128], iconAnchor: [32, 64]
});

var smallBusIcon = L.icon({
    iconUrl: "img/bus_above_active.png", iconSize: [16, 32], iconAnchor: [8, 16]
});
var midBusIcon = L.icon({
    iconUrl: "img/bus_above_active.png", iconSize: [32, 64], iconAnchor: [16, 32]
});
var largeBusIcon = L.icon({
    iconUrl: "img/bus_above_active.png", iconSize: [64, 128], iconAnchor: [32, 64]
});

var smallTrainIcon = L.icon({
    iconUrl: "img/train_above.png", iconSize: [16, 32], iconAnchor: [8, 16]
});
var midTrainIcon = L.icon({
    iconUrl: "img/train_above.png", iconSize: [32, 64], iconAnchor: [16, 32]
});
var largeTrainIcon = L.icon({
    iconUrl: "img/train_above.png", iconSize: [64, 128], iconAnchor: [32, 64]
});


export default class LeafletMap extends Component {
    constructor(props) {
        super(props)
        this.vehicleIcon = midVehicleIcon;
        this.busIcon = midBusIcon;
        this.trainIcon = midTrainIcon;
    }


    componentDidMount() {
        // this.interval = setInterval(() => this.addDemo(), 1000);
    }

    componentWillUnmount() {
        // clearInterval(this.interval)
    }

    componentWillReceiveProps(nextProps) {
        // console.log("Content:willUpdate");
        this.setState(nextProps);
    }


    render() {
        const position = [34.8594, 137.1720];
        let ms = []
        if( this.props.taxi){
            let vs = this.props.store.getVehicle(0); // Car should be ..0

            Object.keys(vs).forEach(function (key) {
                ms.push(
                    <RMarker
                        position={[vs[key][0][0],vs[key][0][1]]}
                        icon={midVehicleIcon}
                        rotationOrigin={(midVehicleIcon.options.iconAnchor[0] + 'px ' + midVehicleIcon.options.iconAnchor[1] + 'px')}
                        rotationAngle ={[vs[key][0][2]]}
                    />
                )
            });
        }
        if( this.props.bus){
            let vs = this.props.store.getVehicle(3);
            Object.keys(vs).forEach(function (key) {
                ms.push(
                    <RMarker
                        position={[vs[key][0][0],vs[key][0][1]]}
                        icon={midBusIcon}
                        rotationOrigin={(midBusIcon.options.iconAnchor[0] + 'px ' + midBusIcon.options.iconAnchor[1] + 'px')}
                        rotationAngle ={[vs[key][0][2]]}
                    />
                )
            });
        }
        if(this.props.busTrace){
            let vs = this.props.store.getVehicle(3);
            Object.keys(vs).forEach(function (key) {

                let vec = vs[key]
                let arr = []
                vec.forEach(function (arg){
                    arr.push([arg[0],arg[1]])
                })
                ms.push(
                    <Polyline
                        positions={arr}
                    />
                )
            });
        }
        if( this.props.train){
            let vs = this.props.store.getVehicle(2);
            Object.keys(vs).forEach(function (key) {
                ms.push(
                    <RMarker
                        position={[vs[key][0][0],vs[key][0][1]]}
                        icon={midTrainIcon}
                        rotationOrigin={(midTrainIcon.options.iconAnchor[0] + 'px ' + midTrainIcon.options.iconAnchor[1] + 'px')}
                        rotationAngle ={[vs[key][0][2]]}
                    />
                )
            });
        }


        let markers = <div></div>;
        if(ms.length > 0){
            markers = ms;
        }

//              url = "https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"

        const map = (
            <Map center={position} zoom={13}>
                <TileLayer
                    url = "https://tiles.wmflabs.org/bw-mapnik/{z}/{x}/{y}.png"
                    attribution ="&copy; <a href=&quot;http://osm.org/copyright&quot;>OpenStreetMap</a> contributors"
                />
                {markers}
            </Map>
        )

//        console.log("zoom:"+map.status.zoom);

        return (
            <section className="content">
                {map}
            </section>
        );
    }
}