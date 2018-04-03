import React, { Component } from 'react';
import Size from './Size';
import Config from './Config';
import { Link } from 'react-router-dom'

class FileListEntry extends Component {
  constructor(props) {
    super(props);

    this.state = {
      showThumb: false,
    };
  }
  mouseIn = () => {
    this.setState({showThumb: true});
  }

  mouseOut = () =>  {
    this.setState({showThumb: false});
  }

  render() {
    var f = this.props.file;
    var style = {};

    if (f.thumbnail !== "" && !this.props.rows) {
      var alt = "thumbnail of " + f.name;
      var thumbnailUrl = Config.root() + "/meta/" + f.hash + "/thumbnail";

      style.backgroundImage = "url('" + thumbnailUrl + "')";

    }


    return (
      <Link style={style} onMouseEnter={this.mouseIn} onMouseLeave={this.mouseOut} className="File" to={"/file/" + f.hash} alt={f.name}>
        <div className="Name">
          {f.name}
        </div>
        <span className="Size"><Size bytes={f.size}/></span>
      </Link>
    )
  }
}

export default FileListEntry;
