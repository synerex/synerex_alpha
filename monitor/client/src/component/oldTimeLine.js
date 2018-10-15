import React, { Component } from 'react';

export default class TimeLine extends Component {
    constructor(props) {
        super(props)
    }

    judgePlayer(log) {
        if (log.msgType.indexOf("Subscribe") >= 0) {
            // log.srcをidに持つ<ul>を作成
            // 作成した<ul>にcount_liの分だけ<li>を追加する
            // 新たに<li>を1つ追加して、innerHTML=log.mstTypeとする
            // 全てのul.logs-ulに<li>を1つ追加する
        } else {
            // log.srcをidに持つ<ul>を探す
            // その<ul>に<li>を1つ追加して、innerHTML=log.mstTypeとする
            // 全てのul.logs-ulに<li>を1つ追加する
        }
    }

    render() {
        const { log } = this.props;
        // count_liを定義する
        return (
            <ul class="timeline">
                <li class="time-label">
                    <span class="bg-red">{log.src}</span>
                </li>
                <li>
                    <i class="fa fa-envelope bg-blue"></i>
                    <div class="timeline-item">
                        <span class="time"><i class="fa fa-clock-o"></i> {log.time}</span>
                        <h3 class="timeline-header"><a href="#">{log.msgType}</a></h3>
                        <div class="timeline-body">
                            {log.chType}, {log.dst}, {log.arg}
                        </div>
                    </div>
                </li>
            </ul>

            // <tr>
            //     <td>{log.msgType}</td>
            //     <td>{log.src}</td>
            //     <td>{log.time}</td>
            //     <td>{log.chType}</td>
            //     <td>{log.dst}, {log.arg}</td>
            // </tr>
        );
    }
}