import React, { Component } from 'react';
import { Link } from 'react-router-dom'

import Size from './Size.js';
import Percent from './Percent.js';
import ProgressBar from './ProgressBar.js';
import './FileItem.css';
import './Buttons.css';

class FileItem extends Component {
  constructor(props) {
    super(props);

    this.state = {
      uploadState: 'queue',
      loaded: 0,
      total: 0,
      response: null,
      hash: null
    };

    this.props.filedata.onFound((data) => {
      console.log('found');
      this.setState({
        loaded: data.bytes_received,
        total: data.size,
        response: data
      });
    });

    this.props.filedata.onNotFound((data) => {
      console.log('not found');
    });

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
      this.setState({uploadState: 'started'});
    });

    this.props.filedata.onAbort(() => {
      this.setState({uploadState: 'aborted'});
    });

    this.props.filedata.onError(() => {
      this.setState({uploadState: 'errored'});
    });

    this.props.filedata.onLoad(() => {
      this.setState({uploadState: 'loaded'});
    });

    this.props.filedata.onResponse((data) => {
      this.setState({
        uploadState: 'response',
        response: data
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

  buttons() {
    var buttons = [];
    var slug = '';
    var r = this.state.response;
    if (r !== undefined && r != null) {
      slug = r.slug;
    }

    var btnOpen = (
      <Link to={"/files/" + slug} key='button_open' className="Button">
      Open
      </Link>
    );

    var btnStart = (
      <span key='button_start' className="Button" onClick={this.start}>
      Start
      </span>
    );

    var btnStop = (
        <span key = 'button_stop' className="Button" onClick={this.stop}>
          Stop
        </span>
    );

    var btnResume = (
      <span key='button_resume' className="Button" onClick={this.start}>
      Resume
      </span>
    );

    var btnRemove = (
      <span key='button_remove' className="Button" onClick={this.remove}>
        Remove
      </span>
    )

    // if (this.fileDone()) {

    // } else if (this.props.filedata.started()) {
    if (this.state.uploadState === 'started') {
      buttons.push(btnStop);
    } else if (this.state.uploadState === 'aborted') {
      buttons.push(btnResume);
    } else if (this.fileDone()) {
      buttons.push(btnOpen);
    } else {
      buttons.push(btnStart);
    }

    buttons.push(btnRemove);

    return(
      <div className="Buttons">
        {buttons}
      </div>
    )
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
    return this.state.loaded === meta.size;
  }

  render() {
    var meta = this.props.filedata.meta();

    var name = meta.name + ' ('  + this.state.uploadState + ') ';
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
    var speed = '';

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
