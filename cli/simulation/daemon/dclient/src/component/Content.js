import React, { Component } from 'react';
import TimeLine from './TimeLine';

export default class Content extends Component {
    constructor(props) {
        super(props)
        this.state = { logs: this.props.logs }
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

    render(){
        const messages =  this.state.logs.map((log, index) => {
            return <TimeLine log ={log} key={index}/>
        });

        return (
            <section className="content">
                <div className="row">

                    <div className="col-md-12">

                            <ul className="timeline">
                                <li className="time-label">{
//<!--                                    <span className="bg-red">
//                                        Synergic Market Server
//                                    </span>
//                                    -->
                                }
                                </li>
                                {messages}
                            </ul>
                    </div>
                </div>
            </section>
        )
    }
}