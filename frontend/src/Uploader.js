import React, { Component } from 'react';
import FileData from './File.js';
import UploadQueue from './UploadQueue.js';
import UploadMenu from './UploadMenu.js';
import './Uploader.css';

class Uploader extends Component {
  constructor(props) {
    super();

    this.state = {
      files: [],
      uploadState: 'blank'
    };

    this.onPaste = this.onPaste.bind(this);
    this.onDrop = this.onDrop.bind(this);
    this.onDragLeave = this.onDragLeave.bind(this);
    this.onDragOver = this.onDragOver.bind(this);

    this.onRemove = this.onRemove.bind(this);
    this.onStart = this.onStart.bind(this);
    this.onStop = this.onStop.bind(this);
    this.onFinish = this.onFinish.bind(this);
    this.onResponse = this.onResponse.bind(this);

    this.onChange = this.onChange.bind(this);
    this.startAll = this.startAll.bind(this);
    this.stopAll = this.stopAll.bind(this);
    this.removeAll = this.removeAll.bind(this);
  }

  componentDidMount() {
    document.addEventListener('paste', this.onPaste);
    document.body.addEventListener('dragover', this.onDragOver);
    document.body.addEventListener('dragleave', this.onDragLeave);
    document.body.addEventListener('drop', this.onDrop);
  }

  addFiles = (files) => {
    var filedatas = [];
    var i = null;

    for(i = 0; i < files.length; i++) {
      var raw = files[i];

      // empty files can be ignored.
      if (raw.size === 0) {
        continue;
      }

      var file = new FileData(raw);
      filedatas.push(file);
    }

    var newFiles = this.state.files;

    for(i = 0; i < filedatas.length; i++) {
      var fd = filedatas[i];
      var canAdd = true;
      for(var j = 0; j < newFiles.length; j++) {
        var f = newFiles[j];
        if (fd._data.name === f._data.name && fd._data.size === f._data.size) {
          canAdd = false;
          break;
        }
      }

      if (canAdd) {
        newFiles.push(fd);
      }
    }

    this.setState({files: newFiles});
  }

  removeFile = (file) => {
    var newFiles = this.state.files;
    var index = newFiles.indexOf(file);
    newFiles.splice(index, 1);
    this.setState({files: newFiles});
  }

  onRemove(file) {
    this.removeFile(file);
  }

  onResponse(e) {
    console.log(e);
  }

  startAll() {
    var files = this.state.files;
    for(var i = 0; i < files.length; i++) {
      var file = files[i];
      file.start();
    }
  }

  stopAll() {
    var files = this.state.files;
    for(var i = 0; i < files.length; i++) {
      var file = files[i];
      file.stop();
    }
  }

  removeAll() {
    this.setState({files: []});
  }

  onChange(e) {
    this.addFiles(e.target.files);
  };

  onStart() {
    this.setState({uploadState: 'started'});
  }

  onStop() {
    this.setState({uploadState: 'stopped'});
  }

  onFinish() {
    this.setState({uploadState: 'finished'});
  }

  // The paste event handler.
  onPaste(e) {
    var files = e.clipboardData.items || [];

    for (var i = 0; i < files.length; i++) {
      var blob = files[i].getAsFile();

      if (blob === null) {
        continue;
      }

      this.addFiles([blob]);
    }
  }

  onDragOver(e) {
    e.stopPropagation();
    e.preventDefault();
  }

  onDragLeave(e) {
    e.stopPropagation();
    e.preventDefault();
  }

  onDrop(e) {
    e.stopPropagation();
    e.preventDefault();

    var files = e.target.files || e.dataTransfer.files;
    for (var i = 0; i < files.length; i++) {
      var file = files[i];

      this.addFiles([file]);
    }
  }

  render() {
    return (
      <div className="Uploader">
        <UploadMenu
          files={this.state.files.length}
          onChange={this.onChange}
          startAll={this.startAll}
          stopAll={this.stopAll}
          removeAll={this.removeAll}
        />
        <UploadQueue
          files={this.state.files}
          onStart={this.onStart}
          onStop={this.onStop}
          onFinish={this.onFinish}
          onRemove={this.onRemove}
        />
      </div>
    );
  }
}

export default Uploader;
