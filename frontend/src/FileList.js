import React, { Component } from 'react';

import Config from './Config';
import FileListEntry from './FileListEntry';
import './FileList2.css';
import './FileList.css';

class Search extends Component {
  render() {
    return <input className="ExtraSearch" onChange={this.props.onChange}/>;
  };
}

class FileList extends Component {
  constructor(props) {
    super(props);

    this.state = { rows: true, filter: null, files: [] };
  };

  componentWillMount() {
    fetch(Config.root() + '/files')
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

  onToggleView = () => {
    this.setState({rows: !this.state.rows});
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
        <FileListEntry key={"fe" + file.slug} file={file}/>
      );
    }

    var className = 'FileList';
    if (this.state.rows) {
      className += ' Rows';
    } else {
      className += ' Grid';
    }

    return(
      <div className={className}>
        <div className="Section">
        Files {files.length}
        <div className="Extra">
          <div onClick={this.onToggleView} className="foo">SUP</div>
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
