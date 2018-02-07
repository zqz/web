import React, { Component } from 'react';
import { Link } from 'react-router-dom'
import Size from './Size';
import './FileList.css';

class FileList extends Component {
  constructor(props) {
    super(props);

    this.state = { files: [] }
  }

  componentWillMount() {
    fetch('http://localhost:3001/files')
    .then(r => {
      return r.json();
    })
    .then(data => {
      this.setState({ files: data });
    });
  }

  render() {
    var files = []

    for(var i = 0; i < this.state.files.length; i++) {
      var file = this.state.files[i];
      var key = file.hash;

      files.push(
        <div key={key} className="File">
          <Link to={"/file/" + file.hash}>{file.name}</Link>
          <span>{file.hash} - <Size bytes={file.size}/></span>
        </div>
      );
    }

    return(
      <div className="FileList">
        <strong>Files</strong>
        {files}
      </div>
    );
  }
}

export default FileList;
