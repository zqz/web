import React, { Component } from 'react';

class FileView extends Component {
  render() {
    return <span>FILEID = {this.props.match.params.id}</span>;
  }
};

export default FileView;
