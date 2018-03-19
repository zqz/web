import React, { Component } from 'react';
import FileItem from './FileItem.js';
import './UploadQueue.css';

class UploadQueue extends Component {
  constructor(props) {
    super(props);

    // for notifying he uploader component
    this.start = this.start.bind(this);
    this.stop = this.stop.bind(this);
  }

  start(file) {
    return () => {
      console.log(file);
      this.props.onStart(file);
    }
  }

  stop(file) {
    return () => {
      console.log(file);
      this.props.onStop(file);
    }
  }

  remove(file) {
    return () => {
      this.props.onRemove(file);
    }
  }

  // done(file) {
  //   return () => {
  //     console.log('finished');
  //     this.props.onStop();
  //   }
  // }

  render() {
    var files = this.props.files;

    if (files.length === 0) {
      return <div/>;
    }

    var fileItems = [];
    for ( var i = 0; i < files.length; i++ ) {
      var f = files[i];
      var k = 'fileitem_' + f.key();
      console.log('key:' + k);

      fileItems.push(
        <FileItem
          key={k}
          filedata={f}
          start={this.start(f)}
          stop={this.stop(f)}
          remove={this.remove(f)}
          // done={this.done(f)}
        />
      )
    }

    return (
      <div className="UploadQueue">
      {fileItems}
      </div>
    );
  }
}

export default UploadQueue;
