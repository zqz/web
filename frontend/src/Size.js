import React, { Component } from 'react';

class Size extends Component {
  bytesToSize(bytes) {
    var sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
    if (bytes === 0) return '0';
    var i = parseInt(Math.floor(Math.log(bytes) / Math.log(1024)), 10);
    var v = bytes / Math.pow(1024, i);
    var n = '';

    if (i === 0) {
      n = Math.round(v, 2);
    } else {
      n = v.toFixed(2);
    }

    return n + ' ' + sizes[i];
  }

  render() {
    return <span>{this.bytesToSize(this.props.bytes)}</span>;
  }
};

export default Size;
