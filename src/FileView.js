import React, { Component } from 'react';

class FileView extends Component {
  constructor(props) {
    super(props);

    this.state = {
      file: null
    };
  }

  componentDidMount() {
    var fileid = this.props.match.params.id;
    console.log(fileid);
    fetch('http://localhost:3001/upload/' + fileid)
    .then(r => {
      return r.json();
    })
    .then(data => {
      this.setState({ file: data });
    });
  }

  render() {
    if (this.state.file === null) {
      return (<div>no file</div>);
    }

    return(
      <div className="FileView">
        <a href={"http://localhost:3001/file/" + this.state.file.hash}>download</a>
      </div>
    );
  }
};

export default FileView;
