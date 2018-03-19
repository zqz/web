import React, { Component } from 'react';
import FileMissing from './FileMissing';
import Size from './Size';
import './FileView.css';

class FileView extends Component {
  constructor(props) {
    super(props);

    this.state = {
      file: null
    };
  }

  componentDidMount() {
    var fileid = this.props.match.params.hash;

    fetch('http://localhost:3001/meta/' + fileid)
    .then(r => {
      if (r.status === 404) {
        return null;
      } else {
        return r.json();
      }
    })
    .then(data => {
      this.setState({ file: data });
    });
  }

  render() {
    var file = this.state.file;

    if (file === null || file === undefined) {
      return <FileMissing/>;
    }

    return(
      <div className="FileView">
        <div className="Section">
          {file.name}
        </div>
        <div className="Content">
          <div className="Left">
            <div>Size: <Size bytes={file.size}/></div>
            <div>Date: {file.date}</div>
            <div>Slug: {file.slug}</div>
            <div>Hash: {file.hash}</div>
            <a className="Download Button" href={"http://localhost:3001/d/" + file.slug}>download</a>
          </div>
          <div className="Right">
            <span className="Link">https://zqz.ca/d/{file.slug}</span>
          </div>
        </div>
      </div>
    );
  }
};

export default FileView;
