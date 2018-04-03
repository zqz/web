import React, { Component } from 'react';
import Size from './Size';
import Config from './Config';
import { Link } from 'react-router-dom'

class FileListEntry extends Component {
  constructor(props) {
    super(props);

    this.state = {
      showThumb: false
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

    var t = (<div></div>);
    if (f.thumbnail !== "" && this.state.showThumb) {
      var alt = "thumbnail of " + f.name;
      t = <img alt={alt} src={Config.root() + "/meta/" + f.hash + "/thumbnail"}/>;
    }

    return (
      <Link onMouseEnter={this.mouseIn} onMouseLeave={this.mouseOut} className="File" to={"/file/" + f.hash} alt={f.name}>
        <div className="Name">
      {t}{f.name}
        </div>
        <span><Size bytes={f.size}/></span>
      </Link>
    )
  }
}

export default FileListEntry;
