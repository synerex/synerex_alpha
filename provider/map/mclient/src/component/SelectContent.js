import React, {Component} from 'react';

// Select Content selects the main component for Admin LTE
export default class SelectContent extends Component {

    constructor(props){
        super(props)
        this.state = { component: this.props.component,
                        args: this.props.args
        }
    }

    componentDidMount(){
    }

    componentWillUnmount(){
    }
    componentWillReceiveProps(nextProps){
        this.setState(nextProps);
    }

    render(){

        return (
            <div className="content-wrapper">
                {
                    React.createElement(this.state.component, this.state.args, "")
                }
            </div>
        )
    }
}
