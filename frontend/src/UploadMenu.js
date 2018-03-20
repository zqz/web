import React, { Component } from 'react';
import UploadButton from './UploadButton.js';
import './UploadMenu.css';
import './Buttons.css';

class UploadMenu extends Component {
  render() {
    if (this.props.files === 0) {
      return(
        <div className="UploadMenu">
          <UploadButton full={true} label="Click here to upload a file or drag a file onto the page" onChange={this.props.onChange}/>
        </div>
      );
    } else {
      return(
        <div className="UploadMenu">
          <UploadButton label="Browse" onChange={this.props.onChange}/>
          <span className="Buttons">
            <span className="Button" onClick={this.props.startAll}>Start All</span>
            <span className="Button" onClick={this.props.stopAll}>Stop All</span>
            <span className="Button" onClick={this.props.removeAll}>Remove All</span>
          </span>
        </div>
      );
    }
  }
}

export default UploadMenu;
