// WorldView

import React from 'react';
import * as  THREE from 'three';
import TrackballControls from './TrackballControls';
import {ResizeListener} from 'react-resize-listener';
import TextBoardCanvas from './TextBoardCanvas';

const MAX_MES_NUM = 500;

// メッセージ の可視化ベース
export default class WorldView extends React.Component {


    constructor(props, context) {
        console.log("World3DView!");
        super(props, context);
        this.cpw = 10; // core line location
        this.lpw = 50; // node location
        this.coreCount = 3;
        this.yaw =0;

        this.text2d = []; // 2d text

        this.cameraPosition = new THREE.Vector3(0, 0, 100);
        this.element = null;
        this.mcount = 0; //  we should fix the number of messages (around 500?)
        this.visible = false;
//        const width = this.element.clientWidth;
//        const height = this.element.clientHeight;
        this.log = this.props.log

        this.store = this.props.store

        this.state = {
            cameraRotation: new THREE.Euler(-1.2, 0, 0),
            width: window.innerWidth,
            height: window.innerHeight-50, // we have to fetch screen size.
        };
        this.vertexShaderStr = "";
        this.fragmentShaderStr = "";

    }

    getCoreVector(type, loc){
       let tv = new THREE.Vector3(this.cpw * Math.cos(Math.PI*2/this.coreCount*type),loc,  this.cpw * Math.sin(Math.PI*2/this.coreCount*type));
       return tv;
    }

    // rotate all world
    rotateCamera(angle){
        let m = new THREE.Matrix4();
        let am = new THREE.Matrix4();
        this.yaw = this.camera.rotation.y;
        this.pitch = this.camera.rotation.x;
        this.roll = this.camera.rotation.z;

        this.yaw += angle;
        m.multiplyMatrices(am.makeRotationY(this.yaw),m);
        m.multiplyMatrices(am.makeRotationX(this.pitch),m);
        m.multiplyMatrices(am.makeRotationZ(this.roll),m);
//        m.multiplyMatrices(am.makeTranslation(0,0,0),m);
        this.camera.position.copy(new THREE.Vector3(0,0,110).applyMatrix4(m));
//        this.camera.rotation.order="YXZ";
        this.camera.rotation.y = this.yaw;
    }

    onAnimate(){
        if(!this.visible) return;
        if( this.controls){
            this.controls.update();
        }else{
            console.log("no control");
        }

        //  Update Scene using current MsgObject
        const lpw = this.lpw;
        const tmspan = 1.8;

        // we need to check non-displayed messages.

        const mct = this.store.getMsgCount();
        // if this.mcount larger then max_message,, 

        if(this.mcount > MAX_MES_NUM){ // we need to update all vectors...
/*            for(let i = 0; i < MAX_MES_NUM ; i ++){
                let ms = this.store.getMsg(mct-MAX_MES_NUM+i);
                const srcIx =  this.store.getNodeIndex(ms.getSrcNodeID())
                let gm = new THREE.Geometry();
                let tm = 20-i*tmspan;
                let cv = this.getCoreVector(ms.getChType(),tm-tmspan/2);
                let fromNodept = new THREE.Vector3(lpw * Math.cos(Math.PI*2/8*srcIx), tm,lpw * Math.sin(Math.PI*2/8*srcIx));
                gm.vertices.push(fromNodept,  cv);
                let str = ms.getMsgType()+":"+ms.getArgs();
                this.add2dRenderText(str,fromNodept, true);
                let mline = new THREE.Line(gm,this.mesMaterial);
                this.mgroup.add(mline)
            }*/
        }else if(mct > this.mcount){
            for(let i = this.mcount; i < mct ; i ++){
                let ms = this.store.getMsg(i);
                const srcIx =  this.store.getNodeIndex(ms.getSrcNodeID())
//                const dstIx =  this.store.getNodeIndex(ms.getDstNodeID())
                let gm = new THREE.Geometry();
                let tm = 20-i*tmspan;
                let cv = this.getCoreVector(ms.getChType(),tm-tmspan/2);
//                console.log("Link Message", ms.getMsgType(), ms.obj);
                let fromNodept = new THREE.Vector3(lpw * Math.cos(Math.PI*2/8*srcIx), tm,lpw * Math.sin(Math.PI*2/8*srcIx));
                gm.vertices.push(fromNodept,  cv);

                let str = ms.getMsgType()+":"+ms.getArgs();

                this.add2dRenderText(str,fromNodept, true);


               let mline = new THREE.Line(gm,this.mesMaterial);
                this.mgroup.add(mline)
            }
            this.mcount = mct;
        }

        // turning ::
        if(this.props.turn){
//            this.scene.rotation.y += 0.003;
//            this.camera.rotation.y += 0.003;
//            this.rotateCamera(0.003);
        }


        this.renderer.render(this.scene,this.camera);

        this.render2d();


        requestAnimationFrame(this.onAnimate.bind(this));
    }

    convertWorldToScreenXY(vec3d){
//        console.log("Before ",vec3d.x,vec3d.y,vec3d.z)
        let vv = new THREE.Vector3(vec3d.x,vec3d.y, vec3d.z);
        vv.project(this.camera)

        const xw = this.state.width /2, yw = this.state.height /2;
        let x = (vv.x +1 )* xw;
        let y = -(vv.y -1)* yw;
//        console.log("Convert ",vec3d.x,vec3d.y, "->" ,x, y);
        return {x:x,y:y}
    }

    // add 2d text on the vector3d.
    add2dRenderText(text,vec3d, withline){
        let vvec = {text:text,vec:vec3d , withline:withline}
        this.text2d.push(vvec);
    }

    // rendering 2d text over the 3d models.
    render2d(){
        const ctx = this.canvas2dctx;

        ctx.clearRect(0,0,this.state.width,this.state.height);

        ctx.strokeStyle = "white";
        ctx.lineWidth = 2.0;
        ctx.font = "16px sans-serif";
        ctx.fillStyle = "white";

        for (let i = 0; i < this.text2d.length; i ++){
            const {text,vec, withline } = this.text2d[i];
            let v2 = this.convertWorldToScreenXY(vec);
            if( withline){
                ctx.beginPath();
                ctx.moveTo(v2.x,v2.y); ctx.lineTo(v2.x+30, v2.y-30); ctx.lineTo(v2.x+30+ctx.measureText(text).width, v2.y-30);
                ctx.stroke();
                ctx.fillText(text, v2.x+30, v2.y-30-3);
            }else{
                ctx.fillText(text, v2.x-ctx.measureText(text).width/2, v2.y-3);
            }
        }
    }

    componentWillUnmount() {
        this.visible = false;
    }

    componentDidMount() {
        console.log("WorldView did mount");
        this.visible = true;

        const {width, height} = this.state;

// Prepare Three Scene
        this.scene = new THREE.Scene();
        this.camera = new THREE.PerspectiveCamera(55, width/ height, 0.00001, 10000);
        this.camera.position.set(0,0,110);

        const threeDOM = window.document.querySelector('#threejs') ;
        this.canvas2d = window.document.querySelector('#canvas2d');
        this.canvas2dctx = this.canvas2d.getContext('2d');
        this.renderer =  new THREE.WebGLRenderer({ antialias:true, canvas : threeDOM });


        const lineMaterial = new THREE.MeshBasicMaterial({color:  0x0066aa});
        const nlineMaterial = new THREE.MeshBasicMaterial({color:  0x003388});
        const coreLineMaterial = new THREE.MeshBasicMaterial({color:  0x003300});
        const clineMaterial = new THREE.MeshBasicMaterial({color:  0x336633});
        this.mesMaterial = new THREE.MeshBasicMaterial({color:  0x660033});
        lineMaterial.transparent = true;

        // make a circle from lines
        let loopGeom = new THREE.Geometry();
        const lpw = this.lpw;
        const cpw = this.cpw;
        const names = ["Lib","RideShare","Ad"];

        // core lines
        for (let i = 0; i < this.coreCount; i++){
            let corelGeom = new THREE.Geometry();
            corelGeom.vertices.push( new THREE.Vector3(cpw * Math.cos(Math.PI*2/this.coreCount*i),20,cpw * Math.sin(Math.PI*2/this.coreCount*i)),
                new THREE.Vector3(cpw * Math.cos(Math.PI*2/this.coreCount*i),-1000,cpw * Math.sin(Math.PI*2/this.coreCount*i)))
            this.coreLine = new THREE.Line(corelGeom,coreLineMaterial)
            this.scene.add(this.coreLine)

            let mesh= new THREE.Mesh(
                new THREE.SphereGeometry(1.2),
                new THREE.MeshBasicMaterial({color: 0x883300})
            );
            mesh.position.x = cpw*Math.cos(Math.PI*2/this.coreCount*i);
            mesh.position.y = 20;
            mesh.position.z = cpw*Math.sin(Math.PI*2/this.coreCount*i);
            this.scene.add(mesh);

            /*
                        let tbc = new TextBoardCanvas({
                            boardWidth: 20,
                            boardHeight: 4,
                            fontSize: 15,
                            textColor: {r:1,g:1,b:1,a:1},
                            backgroundColor: {r:1,g:1,b:1,a:0.1},
                            fontName: "Times New Roman"
                        });
                        tbc.clear();
                        tbc.addTextLine(names[i],25,1);
            //            console.log("Tobj:",i);
                        tbc.update();
                        let tobj = tbc.createPlaneObject();
                        tobj.position.set(cpw * Math.cos(Math.PI*2/this.coreCount*i), 23, cpw * Math.sin(Math.PI*2/this.coreCount*i));
                        this.scene.add(tobj);
            */
            //canvas2d

            this.add2dRenderText(names[i], new THREE.Vector3(cpw * Math.cos(Math.PI*2/this.coreCount*i), 23, cpw * Math.sin(Math.PI*2/this.coreCount*i)), false);

        }

        for (let i = 0; i < 180; i++){
            loopGeom.vertices.push(
                new THREE.Vector3(lpw * Math.cos(Math.PI*i/90), 20,lpw * Math.sin(Math.PI*i/90))
            )
        }

        this.loop = new THREE.LineLoop(loopGeom,lineMaterial)

        let clGeom = new THREE.Geometry();
        clGeom.vertices.push( new THREE.Vector3(0,50,0), new THREE.Vector3(0,-1000,0))
        this.cLine = new THREE.Line(clGeom,clineMaterial)
        this.scene.add(this.cLine)

        const ct = this.store.getNodeNum();

        this.scene.add(this.loop);

        this.nodes = [];
        this.nodeLines = [];
        this.nodeTexts = [];

        console.log("Now count",ct);
        // add each node
        for(let i = 0; i < ct ; i++){
            let mesh= new THREE.Mesh(
                new THREE.SphereGeometry(2.2),
                new THREE.MeshNormalMaterial()
            );
            mesh.position.x = lpw*Math.cos(Math.PI*2/ct*i);
            mesh.position.y = 20;
            mesh.position.z = lpw*Math.sin(Math.PI*2/ct*i);
            this.scene.add(mesh);
            this.nodes.push(mesh);

            let llgeom = new THREE.Geometry();
            llgeom.vertices.push( new THREE.Vector3(mesh.position.x,20,mesh.position.z),
                new THREE.Vector3(mesh.position.x,-1000,mesh.position.z));
            let cl =  new THREE.Line(llgeom,nlineMaterial);
            this.scene.add(cl);
            this.nodeLines.push(cl);
/*
            let tbc = new TextBoardCanvas({
                boardWidth: 20,
                boardHeight: 4,
                fontSize: 15,
                textColor: {r:1,g:1,b:1,a:1},
                backgroundColor: {r:1,g:1,b:1,a:0.1},
                fontName: "Times New Roman"
            });
            tbc.clear();
            // get node info.


            tbc.addTextLine("Node:"+i,25,1);
//            console.log("Tobj:",i);
            tbc.update();
            let tobj = tbc.createPlaneObject();
            tobj.position.set(mesh.position.x, 23, mesh.position.z);
            this.nodeTexts.push(tobj);
            this.scene.add(tobj);
*/

            this.add2dRenderText("Node"+i, new THREE.Vector3(mesh.position.x, 23, mesh.position.z), false);


        }

        this.mgroup = new THREE.Group();
        this.scene.add(this.mgroup);

        this.element = window.document.getElementById("viewer"); // this might be problem with others

//        const controls = new TrackballControls(this.camera,  threeDOM);
        const controls = new TrackballControls(this.camera,  this.canvas2d);

        controls.rotateSpeed = 2.0;
        controls.zoomSpeed = 3.8;
        controls.panSpeed = 1.5;

        controls.noZoom = false;
        controls.noPan = false;
        this.onResizeInner();

        controls.dynamicDampingFactor = 0.3;
        this.controls = controls;

        this.base2dTextSize = this.text2d.length;

       this.done();
    }




    done(){
        this.onAnimate();
    }
    onResizeInner(){
        const comp = document.getElementById('threediv');
        if (comp) {
            const width = comp.clientWidth;
            const height = window.innerHeight-50;
//            console.log("Resize!",width,height);
            this.setState({
                width: width,
                height: height
            });

            if (this.camera) {
                this.camera.aspect = width / height;
                this.camera.updateProjectionMatrix();
                this.renderer.setSize(width, height);
            }else{
                console.log("Not yet camera");
            }
        }

    }


    render() {
        const { data } = this.state; // eslint-disable-line no-unused-vars
        const comp = document.getElementById('threediv');
        let width,height;
        if (comp ) {
            width = comp.clientWidth;
            height = comp.clientHeight;
            if( width !== this.state.width || height !== this.state.height){
                console.log("Diff with state");
                if (this.camera) {
                    this.camera.aspect = width / height;
                    this.camera.updateProjectionMatrix();
                    this.renderer.setSize(width, height);
                }

            }else{
            }
        }else {
            width = this.state.width;
            height = this.state.height;
        }

        if(this.store.getMsgCount() === 0 && this.mgroup !== undefined){ // clear button
            this.mgroup.children=[];
            this.mcount = 0;
            this.text2d.splice(this.base2dTextSize);
//            console.log("Clear",this.base2dTextSize,this.text2d.length);
        }
        return(
            <div id="threediv" style={{position: 'relative'}}>
                <ResizeListener  onResize={this.onResizeInner.bind(this)} />
                <canvas id="threejs" style={{display: "block"}} width={width} height={height} />
                <canvas id="canvas2d" style={
                    {
                        position: "absolute",
                        left: 0,
                        top: 0,
                        backgroundColor: "transparent"
                    }
                }
                width={width} height={height} />
            </div>
        );
    }

}
