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

class Toggle extends Component {
  render() {
    var text = 'grid';
    if (this.props.rows) {
      text = 'rows';
    }

    return (
      <div className={"pointer ico " + text} onClick={this.props.onClick}></div>
    );
  }
}

const FILELIST_CFG = 'filelist-rows';

class FileList extends Component {
  constructor(props) {
    super(props);

    this.state = {
      rows: Config.get(FILELIST_CFG) === 'true',
      filter: null,
      files: []
    };
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
    var v = !this.state.rows;
    console.log('v is now', v);
    Config.set(FILELIST_CFG, v.toString());
    this.setState({rows: v});
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
        <FileListEntry rows={this.state.rows} key={"fe" + file.slug} file={file}/>
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
          <Search onChange={this.onSearch}/>
          <Toggle rows={this.state.rows} onClick={this.onToggleView}/>
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
