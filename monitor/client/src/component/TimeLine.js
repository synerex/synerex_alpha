import React, { Component } from 'react';

export default class TimeLine extends Component {
    constructor(props) {
        super(props)
    }

    render() {
        const { log } = this.props;
        return (
            <li>
                <i className="fa fa-envelope bg-blue"></i>
                <div className="timeline-item">
                    <span className="time"><i className="fa fa-clock-o"></i> {log.time}</span>
                    <h3 className="timeline-header"><a href="#">{log.msgType}</a></h3>
                    <div className="timeline-body">
                        {log.chType}, {log.dst}, {log.arg}
                    </div>
                </div>
            </li>
        );
    }
}