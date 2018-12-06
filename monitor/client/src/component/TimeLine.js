import React, { Component } from 'react';

export default class TimeLine extends Component {
    constructor(props) {
        super(props)
    }
/*
<!--                <i className="fa fa-envelope bg-blue"></i> -->

<!--                    <span className="time"><i className="fa fa-clock-o"></i> {log.time}</span>
                    <h3 className="timeline-header"><a href="#">{log.msgType}</a></h3> -->
 */

    render() {
        const { log } = this.props;
        return (
            <li>
                <div className="timeline-item">
                    <div className="timeline-body">
                        {log.msgType}: {log.chType}, id:{log.id}, dst:{log.dst}, target:{log.tgt}, {log.arg}
                    </div>
                </div>
            </li>
        );
    }
}