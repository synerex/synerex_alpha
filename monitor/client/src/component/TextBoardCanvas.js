//平面オブジェクト＋テクスチャマッピングによる
// sample http://www.natural-science.or.jp/article/20161205233604.php

import * as  THREE from 'three';

//
export default class  TextBoardCanvas {

    constructor(parameter) {
        parameter = parameter || {};

        //背景色（RGBA値を0から１で指定）
        this.backgroundColor = parameter.backgroundColor || {r: 1, g: 1, b: 1, a: 1};
        //文字色（RGBA値を0から１で指定）
        this.textColor = parameter.textColor || {r: 0, g: 0, b: 0, a: 1};

        //マッピング対象オブジェクトのサイズ（縦横比は２のべき乗を推奨）
        this.boardWidth = parameter.boardWidth || 100;
        this.boardHeight = parameter.boardHeight || 100;

        //フォントサイズと行間（canvas要素の横幅に対する[%]で指定）
        this.fontSize = parameter.fontSize || 10;      //フォントサイズ
        this.lineHeight = parameter.lineHeight || 1.1; //行間

        //フォント名（CSSで指定可能な形式）
        this.fontName = parameter.fontName || "serif"; //フォント名
        //解像度
        this.resolution = parameter.resolution || 4;

        this._lineHeight = 0;
        this.textLines = [];

        this.plane = null;
        this.sprite = null;

        this.init();
    }

    init() {
        //canvas要素の生成
        this.canvas = document.createElement('canvas');
        //canvas要素のサイズ
        this.canvas.width = Math.pow(2, Math.floor(Math.log2(this.boardWidth)) + this.resolution);  //横幅
        this.canvas.height = Math.pow(2, Math.floor(Math.log2(this.boardHeight)) + this.resolution); //縦幅

//        console.log("canvas要素のサイズ：", this.canvas.width, "×", this.canvas.height);

        //コンテキストの取得
        this.canvas.context = this.canvas.getContext('2d');

        this.setBackGroundColor(this.backgroundColor.r, this.backgroundColor.g, this.backgroundColor.b, this.backgroundColor.a);
        this.setTextColor(this.textColor.r, this.textColor.g, this.textColor.b, this.textColor.a);
        this.setFontSize(this.fontSize);
        this.setFontName(this.fontName);
        this.setLineHeight(this.lineHeight)
    }

//背景色の設定
    setBackGroundColor(r, g, b, a) {

        this.backgroundColor.r = r || 0;
        this.backgroundColor.g = g || 0;
        this.backgroundColor.b = b || 0;
        this.backgroundColor.a = a || 0;

        this.canvas.context.fillStyle = "rgba(" + 255 * this.backgroundColor.r + " ," + 255 * this.backgroundColor.g + " ," + 255 * this.backgroundColor.b + " ," + this.backgroundColor.a + ")";
        this.canvas.context.fillRect(0, 0, this.canvas.width, this.canvas.height);

    }

//全消し
    clear() {

        this.canvas.context.clearRect(0, 0, this.canvas.width, this.canvas.height);
        this.canvas.context.fillStyle = "rgba(" + 255 * this.backgroundColor.r + " ," + 255 * this.backgroundColor.g + " ," + 255 * this.backgroundColor.b + " ," + this.backgroundColor.a + ")";
        this.canvas.context.fillRect(0, 0, this.canvas.width, this.canvas.height);
        this.canvas.context.fillStyle = "rgba(" + 255 * this.textColor.r + " ," + 255 * this.textColor.g + " ," + 255 * this.textColor.b + " ," + this.textColor.a + ")";
        this._lineHeight = 0;
    }

//文字色の設定
    setTextColor(r, g, b, a) {
        this.textColor.r = r || 0;
        this.textColor.g = g || 0;
        this.textColor.b = b || 0;
        this.textColor.a = a || 0;

        this.canvas.context.fillStyle = "rgba(" + 255 * this.textColor.r + " ," + 255 * this.textColor.g + " ," + 255 * this.textColor.b + " ," + this.textColor.a + ")";
    }

//文字サイズの設定
    setFontSize(size) {
        this.fontSize = size || 10;
        this.canvas.context.font = this.fontSize / 100 * this.canvas.width + "px " + this.fontName;
    }

//フォントの設定
    setFontName(name) {
        this.fontName = name || "serif";
        this.canvas.context.font = this.fontSize / 100 * this.canvas.width + "px " + this.fontName;
    }

//行間の設定
    setLineHeight = function (height) {
        this.lineHeight = height || 1.1;
    }

//文字列の追加
    addTextLine(text, indent, lineHeight) {
        text = text || "";
        indent = indent || 0;
        lineHeight = lineHeight || this.lineHeight;

        this.textLines.push({text: text, indent: indent, lineHeight: lineHeight});
        this._lineHeight += lineHeight * this.fontSize / 100 * this.canvas.width;

        this.canvas.context.fillText(
            text,
            indent / 100 * this.canvas.width,
            this._lineHeight
        );

    }

//canvas要素を取得
    getTextCanvas() {
        return this.canvas;
    }

    createPlaneObject() {

        //テクスチャ画像用のcanvas要素の取得
        var canvas = this.getTextCanvas();
        //テクスチャオブジェクトの生成
        this.texture = new THREE.Texture(canvas);
        //テクスチャ画像の更新
        this.texture.needsUpdate = true;

        //形状オブジェクトの宣言と生成
        var geometry = new THREE.PlaneGeometry(this.boardWidth, this.boardHeight);
        //材質オブジェクトの宣言と生成
        var material = new THREE.MeshBasicMaterial({map: this.texture, transparent: true});
        //平面オブジェクトの生成
        this.plane = new THREE.Mesh(geometry, material);

        return this.plane;
    }

    createSpriteObject() {

        //テクスチャ画像用のcanvas要素の取得
        var canvas = this.getTextCanvas();
        //テクスチャオブジェクトの生成
        this.texture = new THREE.Texture(canvas);
        //テクスチャ画像の更新
        this.texture.needsUpdate = true;

        //材質オブジェクトの宣言と生成
        var material = new THREE.SpriteMaterial({map: this.texture});
        //スプライトオブジェクトの生成
        this.sprite = new THREE.Sprite(material);

        this.sprite.scale.set(this.boardWidth, this.boardHeight, 1);

        return this.sprite;
    }

    createTextScreen() {
        this.textScene = new THREE.Scene();
        if (this.sprite) {
            this.textScene.add(this.sprite);
        } else {
            this.textScene.add(this.createSpriteObject());
        }

        this.textCamera = new THREE.OrthographicCamera(-this.boardWidth / 2, this.boardWidth / 2, this.boardHeight / 2, -this.boardHeight / 2, -10, 10);

    }

    update() {

        if (this.plane) this.plane.material.map.needsUpdate = true;
        if (this.sprite) this.sprite.material.map.needsUpdate = true;

    }

    getPlaneObject() {
        return this.plane;
    }

    getSpriteObject() {
        return this.sprite;
    }
}
