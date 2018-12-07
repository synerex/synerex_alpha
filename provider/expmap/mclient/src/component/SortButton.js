import React, { Component } from 'react';

export default class SortButton extends Component {
    sortByAsc(e) {
        e.preventDefault();
        this.props.onSortByAsc(e.target.value)
    }
    sortByDesc(e) {
        e.preventDefault();
        this.props.onSortByDesc(e.target.value)
    }
    render() {
        return (
            <div class="info-box">
                <span class="info-box-icon bg-yellow"><i class="fa fa-sort-amount-desc"></i></span>
                <div class="info-box-content">
                    <span class="info-box-text">日付でソートする</span>
                    <div class="btn-group">
                        <button className="btn btn-default" onClick={this.sortByAsc.bind(this)} value="time">昇順</button>
                        <button className="btn btn-default" onClick={this.sortByDesc.bind(this)} value="time">降順</button>
                    </div>
                </div>
            </div>
        );
    }
}