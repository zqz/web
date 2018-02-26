import React, { Component } from 'react';
import FileSize from './FileSize';

class Size extends Component {
  bytesToSize(bytes) {
    return FileSize.ToSize(bytes);
  }

  render() {
    return <span>{this.bytesToSize(this.props.bytes)}</span>;
  }
};

export default Size;
