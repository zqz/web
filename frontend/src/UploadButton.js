import React, { Component } from 'react';
import './UploadButton.css'

class UploadButton extends Component {
  constructor(props) {
    super(props);

    this.onClick = this.onClick.bind(this);
  }

  onClick(e) {
    this.refs.uploader.click();
  }

  render() {
    var classes="Button UploadButton";
    if (this.props.full) {
      classes += " Full"
    }
    return(
      <div className={classes}>
        <div className="Fake" onClick={this.onClick}>
          {this.props.label}
        </div>
        <input
          ref="uploader"
          className="Real"
          type="file"
          onChange={this.props.onChange}
          multiple
        />
      </div>
    )
  }
}

export default UploadButton;
