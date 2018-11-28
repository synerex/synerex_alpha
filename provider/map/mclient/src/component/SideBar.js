// SideBar.js

import React, {Component} from 'react';

export default class SideBar extends Component {
    constructor(props){
        super(props);
    }

    render(){
        return (
            <aside className="main-sidebar">
                <section className="sidebar">
                    {/*
                    <div className="user-panel">
                        <div className="pull-left image">
                            <img src="img/user2-160x160.jpg" className="img-circle" alt="User Image" />
                        </div>
                        <div className="pull-left info">
                            <p>Alexander Pierce</p>
                            <a href="#"><i className="fa fa-circle text-success"></i> Online</a>
                        </div>
                    </div>
                    */}
                    <form action="#" method="get" className="sidebar-form">
                        <div className="input-group">
                        <input type="text" name="q" className="form-control" placeholder="Filter..." />
                        <span className="input-group-btn">
                                <button type="submit" name="search" id="search-btn" className="btn btn-flat"><i className="fa fa-search"></i>
                                </button>
                            </span>
                        </div>
                    </form>
                    <ul className="sidebar-menu" data-widget="tree">
                        <li className="header">Controls</li>
                        <li>
                            <a href="#" onClick={()=>{this.props.resetView()}}>
                                <i className="fa fa-th"></i> <span>Reset View</span>
                                <span className="pull-right-container">
                            </span>
                            </a>
                        </li>
                        <li>
                            <a href="#" onClick={()=>{this.props.showBus()}}>
                                <i className="fa fa-th"></i> <span>Show Bus</span>
                                <span className="pull-right-container">
                            </span>
                            </a>
                        </li>
                        <li>
                            <a href="#" onClick={()=>{this.props.showBusTrace()}}>
                                <i className="fa fa-th"></i> <span>Show Bus Trace</span>
                                <span className="pull-right-container">
                            </span>
                            </a>
                        </li>
                        <li>
                            <a href="#" onClick={()=>{this.props.showTaxi()}}>
                                <i className="fa fa-th"></i> <span>Show Taxi</span>
                                <span className="pull-right-container">
                            </span>
                            </a>
                        </li>
                        <li>
                            <a href="#" onClick={()=>{this.props.showTrain()}}>
                                <i className="fa fa-th"></i> <span>Show Train</span>
                                <span className="pull-right-container">
                            </span>
                            </a>
                        </li>

                    </ul>
                </section>
            </aside> 
        )
    }
}
