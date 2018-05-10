import React, { Component } from "react";
import FileMissing from "./FileMissing";
import Size from "./Size";
import Config from "./Config";
import Preview from "./Preview";
import "./FileView.css";

class FileView extends Component {
  constructor(props) {
    super(props);

    this.state = {
      file: null,
      gotFile: false,
    };
  }

  getDerivedStateFromProps() {

  }

  componentWillReceiveProps(nextProps) {
    var hash = nextProps.match.params.hash;
    if (hash !== this.state.hash) {
      this.setState({
        hash: hash,
        file: null
      });
    }
  }

  componentDidMount() {
    var hash = this.props.match.params.hash;

    this.setState({
      hash: hash,
      file: null
    });
  }

  work() {
    var hash = this.props.match.params.hash;

    fetch(Config.root() + "/meta/" + hash)
    .then(r => {
      if (r.status === 404) {
        return null;
      } else {
        return r.json();
      }
    })
    .then(data => {
      this.setState({ file: data });
    });
  }

  render() {
    if (this.state.file === null) {
      console.log("lols");
      this.work();
      return(<div>loading</div>);
    }
    var file = this.state.file;
      console.log("lols");

    if (file === null) {
      return <FileMissing/>;
    }

    var dlurl = Config.cdnroot() + file.slug;

    return(
      <div className="FileView">
        <div className="Section">
          {file.name}
        </div>
        <div className="Preview">
          <Preview file={file}/>
        </div>
        <div className="Content">
          <div className="Left">
            <div>Size: <Size bytes={file.size}/></div>
            <div>Date: {file.date}</div>
            <div>Slug: {file.slug}</div>
            <div>Hash: {file.hash}</div>
            <a className="Download Button" href={dlurl}>download</a>
          </div>
          <div className="Right">
            <span className="Link">{dlurl}</span>
          </div>
        </div>
      </div>
    );
  }
};

export default FileView;
