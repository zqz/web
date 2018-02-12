import React, { Component } from 'react';
import './ProgressBar.css';

class ProgressBar extends Component {
  render() {
    var v = this.props.value;

    if (isNaN(v)) {
      v = 0;
    }

    var style = {
      width: (v * 100) + '%'
    }

    return(
      <div className="ProgressBar">
        <div className="Bar" style={style}></div>
      </div>
    )
  }
}

export default ProgressBar;
