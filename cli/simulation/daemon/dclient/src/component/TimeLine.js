import React, { Component } from 'react';

export default class TimeLine extends Component {
    constructor(props) {
        super(props)
    }

    render() {
        const { log } = this.props;
        return (
            <li>
                {/* <i className="fa fa-envelope bg-blue"></i> */}
                <div className="timeline-item">
                    {/* <span className="time"><i className="fa fa-clock-o"></i> {this.props.log.t}</span> */}
                    {/* <h3 className="timeline-header"><a href="#"> {this.props.log.m} </a> </h3> */}
                    <div className="box box-solid box-default">
                        <div className="box-header with-border">
                            <h3 className="box-title">{log.message}</h3>
                            <div className="box-tools pull-right">
                                <button className="btn btn-box-tool"><i class="fa fa-times"></i></button>
                            </div>
                        </div>
                        <div className="box-body">
                            <i className="fa fa-clock-o"></i> {log.time}
                        </div>
                        <div className="box-body">
                            {log.value}
                        </div>
                    </div>

                    {/*
                        <div className="timeline-body">
                            {this.props.log.v}
                        </div>
                        {/*
                        <div className="timeline-footer">
                            <a className="btn btn-primary btn-xs">...</a>
                        </div>
                    */}
                </div>
            </li>
        )
    }
}