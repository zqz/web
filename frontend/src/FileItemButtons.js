import React, { Component } from "react";
import { Link } from "react-router-dom";
import "./Buttons.css";


class Button extends Component {
  render() {
    var p = this.props;
    return(
      <span key={"button_" + p.key} className="Button" onClick={p.onClick}>
        {p.text}
      </span>
    )
  }
}


class FileItemButtons extends Component {
  open = () => <Link to={"/file/" + this.props.hash} key="button_open" className="Button">Open</Link>
  start = () => <Button key="start" text="Start" onClick={this.props.onStart}/>
  stop = () => <Button key="stop" text="Stop" onClick={this.props.onStop}/>
  resume = () => <Button key="resume" text="Resume" onClick={this.props.onResume}/>
  remove = () => <Button key="remove" text="Remove" onClick={this.props.onRemove}/>

  buttons() {
    var buttons = [];
    var s = this.props.uploadState;

    if (s === "started") {
      buttons.push(this.stop());
    } else if (s === "aborted") {
      buttons.push(this.resume());
    } else if (this.props.done()) {
      buttons.push(this.open());
    } else {
      buttons.push(this.start());
    }
    buttons.push(this.remove());

    return buttons;
  }

  render() {
    return (
      <div className="Buttons">
        {this.buttons()}
      </div>
    )
  }
}

export default FileItemButtons;
