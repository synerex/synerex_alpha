// SideBar.js

import React, {Component} from 'react';

export default class SideBar extends Component {
    constructor(props){
        super(props);
    }

    render(){
        let addLists=[];
        for(let i = 0; i < this.props.providers.length; i++){
            addLists[i]=
                <li><a href="#" onClick={()=>{this.props.start(this.props.providers[i])}}><i className="fa fa-circle-o"></i>{this.props.providers[i]}</a></li>
        }


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
                            <a href="#" onClick={()=>{this.props.startNodeServ()}}>
                                <i className="fa fa-th"></i> <span>NodeServ</span>
                                <span className="pull-right-container">
                            </span>
                            </a>
                        </li>
                        <li>
                            <a href="#" onClick={()=>{this.props.start("monitor")}}>
                                <i className="fa fa-th"></i> <span>Monitor</span>
                                <span className="pull-right-container">
                               </span>
                            </a>
                        </li>
                        <li>
                            <a href="#" onClick={()=>{this.props.start("smarket")}}>
                                <i className="fa fa-th"></i> <span>SynerexServer</span>
                                <span className="pull-right-container">
                            </span>
                            </a>
                        </li>
                        <li>
                        <a href="#" onClick={()=>{this.props.clearLogs()}}>
                            <i className="fa fa-th"></i> <span>ClearLog</span>
                            <span className="pull-right-container">
                            </span>
                        </a>
                        </li>



                    </ul>

                    <ul className="sidebar-menu" data-widget="tree">
                        <li className="header">Providers</li>
                        <li className="treeview">
                            <a href="#">
                                <i className="fa fa-th"></i>
                                <span>Kota Providers</span>
                                <span className="pull-right-container">
                                <span className="label label-primary pull-right">{this.props.providers.length}</span>
                                </span>
                            </a>
                            <ul className="treeview-menu">
                                {
                                    addLists
                                }
                            </ul>
                        </li>

                        {/*
                        <a href="pages/widgets.html">
                            <i className="fa fa-th"></i> <span>Widgets</span>
                            <span className="pull-right-container">
                            <small className="label pull-right bg-green">new</small>
                            </span>
                        </a>
                        </li>
                        <li className="treeview">
                        <a href="#">
                            <i className="fa fa-pie-chart"></i>
                            <span>Charts</span>
                            <span className="pull-right-container">
                            <i className="fa fa-angle-left pull-right"></i>
                            </span>
                        </a>
                        <ul className="treeview-menu">
                            <li><a href="pages/charts/chartjs.html"><i className="fa fa-circle-o"></i> ChartJS</a></li>
                            <li><a href="pages/charts/morris.html"><i className="fa fa-circle-o"></i> Morris</a></li>
                            <li><a href="pages/charts/flot.html"><i className="fa fa-circle-o"></i> Flot</a></li>
                            <li><a href="pages/charts/inline.html"><i className="fa fa-circle-o"></i> Inline charts</a></li>
                        </ul>
                        </li>
                        <li>
                        <ul className="treeview-menu">
                            <li><a href="pages/tables/simple.html"><i className="fa fa-circle-o"></i> Simple tables</a></li>
                            <li><a href="pages/tables/data.html"><i className="fa fa-circle-o"></i> Data tables</a></li>
                        </ul>
                        </li>
                        <li>
                        <a href="pages/calendar.html">
                            <i className="fa fa-calendar"></i> <span>Calendar</span>
                            <span className="pull-right-container">
                            <small className="label pull-right bg-red">3</small>
                            <small className="label pull-right bg-blue">17</small>
                            </span>
                        </a>
                        </li>
                        <li>
                        <a href="pages/mailbox/mailbox.html">
                            <i className="fa fa-envelope"></i> <span>Mailbox</span>
                            <span className="pull-right-container">
                            <small className="label pull-right bg-yellow">12</small>
                            <small className="label pull-right bg-green">16</small>
                            <small className="label pull-right bg-red">5</small>
                            </span>
                        </a>
                        </li>
                        */}
                    </ul>
                </section>
            </aside> 
        )
    }
}
