import React, { Component } from 'react';
import Config from './Config';
import './Preview.css';

class Preview extends Component {

  inside() {
    var f = this.props.file;
    if (f == null) {
      return <div>no file</div>;
    }

    if (f.type.startsWith('image')) {
      var imgpath = Config.cdnroot() + f.slug;
      var alt = "preview of " + f.name;
      return (<img alt={alt} src={imgpath}/>);
    }
  }

  render() {
    return (
      <div>
        <div>Preview</div>
        {this.inside()}
      </div>
    );
  }
}

export default Preview;
