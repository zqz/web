import React, { Component } from "react";

import Size from "./Size.js";
import Percent from "./Percent.js";
import ProgressBar from "./ProgressBar.js";
import FileItemButtons from "./FileItemButtons.js";
import "./FileItem.css";

class FileItem extends Component {
  constructor(props) {
    super(props);

    this.state = {
      uploadState: "queue",
      loaded: 0,
      total: 0,
      slug: null,
      response: null,
      hash: null
    };

    this.props.filedata.onFound((data) => {
      this.setState({
        loaded: data.bytes_received,
        total: data.size,
        slug: data.slug,
        response: data,
      });
    });

    this.props.filedata.onNotFound((data) => {});

    this.props.filedata.onExists((data) => {
      this.setState({loaded: data.bytes_received, total: data.size});
    });

    this.props.filedata.onHash((h) => {
      this.setState({hash: h});
      this.props.filedata.check(h);
    });

    this.props.filedata.onProgress((e) => {
      this.setState({loaded: e.loaded, total: e.total});
    });

    this.props.filedata.onStart(() => {
      this.setState({uploadState: "started"});
    });

    this.props.filedata.onAbort(() => {
      this.setState({uploadState: "aborted"});
    });

    this.props.filedata.onError(() => {
      this.setState({uploadState: "errored"});
    });

    this.props.filedata.onLoad(() => {
      this.setState({uploadState: "loaded"});
    });

    this.props.filedata.onResponse((data) => {
      this.setState({
        uploadState: "response",
        response: data,
        slug: data.slug
      });
    });
  }

  componentDidMount() {
    this.props.filedata.hash();
  }

  start = () => {
    // this.props.filedata.start();
    this.props.filedata.start();
    this.props.start();
  }

  stop = () => {
    this.props.filedata.stop();
    this.props.stop();
  }

  remove = () => {
    this.props.remove();
  }

  buttons = () => {
    return <FileItemButtons
      slug={this.state.slug}
      done={this.fileDone}
      onStart={this.start}
      onStop={this.stop}
      onResume={this.start}
      onRemove={this.remove}
      uploadState={this.state.uploadState}
    />;
  }

  onClickName = (e) => {
    if (this.fileDone()) {
      return;
    }

    var t = e.target;
    t.contentEditable = true;
    t.oninput = function(e) {
      t.textContent = t.textContent.replace(/(\r\n|\n|\r)/gm,"");
    }
  }

  fileDone = () => {
    var meta = this.props.filedata.meta();
    if (this.state.loaded === 0) {
      return false;
    }
    return this.state.loaded === meta.size;
  }

  render() {
    var meta = this.props.filedata.meta();

    var name = meta.name + " ("  + this.state.uploadState + ") ";
    var perc = this.state.loaded / meta.size;
    var basic = null;

    if (this.fileDone()) {
      basic = <span>finished</span>;
    } else {
      basic =(
        <span className="Percent">
          <Size bytes={this.state.loaded}/> / <Size bytes={meta.size}/> - <Percent value={perc}/>
        </span>
      );
    };

    var progress = (
      <span>
      {basic}
      </span>
    );

    var buttons = this.buttons();
    var hash = this.state.hash;
    var speed = "";

    if (!this.fileDone()) {
      var progressBar = <ProgressBar value={perc}/>;
      speed = this.props.filedata.speed();
    }

    return (
      <div className={"FileItem " + this.state.uploadState}>
        <div className="Side">
        </div>
        <div className="Main">
          <div className="Top">
            <span className="Left">
              <span className="Name" onClick={this.onClickName}>{name}</span>
            </span>
            <span className="Buttons">
              {buttons}
            </span>
          </div>
          {progressBar}
          <div className="Bottom">
            <span>{hash}</span>
            <span className="Progress">{speed} {progress}</span>
          </div>
        </div>
      </div>
    );
  }
}

export default FileItem;
