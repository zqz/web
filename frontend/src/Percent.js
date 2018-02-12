import React, { Component } from 'react';

class Percent extends Component {
  render() {
    var v = this.props.value;
    var n = '0';
    if (isNaN(v)) {
      n = '0';
    } else if (v == 0) {
      n = '0'
    } else if (v >= 1) {
      n = '100'
    } else {
      n = (v * 100).toFixed(2)
    }

    var perc = n + '%';
    return <span>{perc}</span>;
  }
};

export default Percent;
