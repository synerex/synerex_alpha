import React, { Component } from 'react';

export default class FilterForm extends Component {
    filterVal() {
        const val = this.refs.myInput.value;
        this.props.onFilterVal(val);
    }
    render() {
        return (
            <div className="box-tools">
                <div className="input-group input-group-sm" style={{ width: '200px' }}>
                    <input
                        type="text"
                        ref="myInput"
                        defaultValue=""
                        onKeyUp={this.filterVal.bind(this)}
                        name="table_search"
                        className="form-control pull-right"
                        placeholder="キーワードで絞り込む"
                    />
                    <div class="input-group-btn">
                        <button className="btn btn-default"><i class="fa fa-search"></i></button>
                    </div>
                </div>
            </div>
        );
    }
}