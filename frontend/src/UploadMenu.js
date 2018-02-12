import React, { Component } from 'react';
import UploadButton from './UploadButton.js';
import './UploadMenu.css';
import './Buttons.css';

class UploadMenu extends Component {
  render() {
    return(
      <div className="UploadMenu">
        <UploadButton onChange={this.props.onChange}/>
        <span className="Buttons">
          <span className="Button" onClick={this.props.startAll}>Start All</span>
          <span className="Button" onClick={this.props.stopAll}>Stop All</span>
          <span className="Button" onClick={this.props.removeAll}>Remove All</span>
        </span>
      </div>
    )
  }
}

export default UploadMenu;
