import React, { Component } from 'react';
import ExifOrientationImg from 'react-exif-orientation-img'

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

      if (f.type.endsWith('jpeg') || f.type.endsWith('jpg')) {
        return (<ExifOrientationImg src={imgpath} alt={alt}/>);
      } else {
        return (<img alt={alt} src={imgpath}/>);
      }
    }
  }

  render() {
    return (
      <div>
        {this.inside()}
      </div>
    );
  }
}

export default Preview;
