import React, { Component } from 'react';

class FileView extends Component {
  constructor(props) {
    super(props);

    this.state = {
      file: null
    };
  }

  componentDidMount() {
    var fileid = this.props.match.params.hash;
    console.log(fileid);
    fetch('http://localhost:3001/file/' + fileid)
    .then(r => {
      return r.json();
    })
    .then(data => {
      this.setState({ file: data });
    });
  }

  render() {
    var file = this.state.file;

    if (file === null) {
      return (<div>no file</div>);
    }

    return(
      <div className="FileView">
        <div>
          <strong>{file.name}</strong>
        </div>
        <div>Size: {file.size}</div>
        <div>Date: {file.date}</div>
        <a href={"http://localhost:3001/file/" + file.token + '/download'}>download</a>
      </div>
    );
  }
};

export default FileView;
