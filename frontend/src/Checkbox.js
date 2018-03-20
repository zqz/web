import React, { Component } from 'react';
import './Checkbox.css';

class Checkbox extends Component {
  render() {
    var classes = ['Checkbox'];

    if (this.props.checked === true) {
      classes.push('Active');
    }

    return (
      <div onClick={this.props.onClick} className="Checkbox-Container">
        <div className="Label">
          {this.props.label}
        </div>
        <div className={classes.join(" ")}>
          <div className="Inside"></div>
        </div>
        <div className="Help">
          {this.props.desc}
        </div>
      </div>
    );
  }
};

export default Checkbox;
