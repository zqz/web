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
    return(
      <div className="Button UploadButton">
      <div className="Fake" onClick={this.onClick}>
        Browse or drag file onto page
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
