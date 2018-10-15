import React, { Component } from 'react';

export default class ProcessLog extends Component {
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
            return (<span>{log.message}<br /></span>)
        });

        return (
            <section className="content">
                {messages}
            </section>
        )
    }
}