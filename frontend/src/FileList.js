import React, { Component } from 'react';
import { Link } from 'react-router-dom'

import Size from './Size';
import './FileList2.css';

class Search extends Component {
  render() {
    return <input className="ExtraSearch" onChange={this.props.onChange}/>;
  };
}

class FileList extends Component {
  constructor(props) {
    super(props);

    this.state = { filter: null, files: [] };
  };

  componentWillMount() {
    fetch('http://localhost:3001/files')
    .then(r => {
      return r.json();
    })
    .then(data => {
      this.setState({ files: data });
    })
    .catch(e => { console.log('error', e); })
  }

  onSearch = (e) => {
    var f = e.target.value;
    if (f === undefined || f === null || f.length === 0) {
      f = null;
    }

    this.setState({filter: f});
  }

  render() {
    var files = []

    for(var i = 0; i < this.state.files.length; i++) {
      var file = this.state.files[i];

      var filter = this.state.filter;
      if (filter !== null) {
        if (file.name.toLowerCase().indexOf(filter) < 0) {
          continue;
        }
      }

      files.push(
        <Link className="File" to={"/file/" + file.hash} alt={file.name}>
          <div className="Name">
          {file.name}
          </div>
          <span><Size bytes={file.size}/></span>
        </Link>
      );
    }

    return(
      <div className="FileList">
        <div className="Section">
        Files {files.length}
        <div className="Extra">
          <Search onChange={this.onSearch}/>
        </div>
        </div>
        <div className="List">
          {files}
        </div>
      </div>
    );
  }
}

export default FileList;
