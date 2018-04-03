import React, { Component } from 'react';
import Checkbox from './Checkbox';
import Config from './Config';
import './Settings.css';

class Settings extends Component {
  constructor(props) {
    super(props);

    this.state = {
      instant: this.get('instant') === 'true',
      other: this.get('other') === 'true'
    };
  }

  set(k, v) {
    Config.set(k, v);
  }

  get(k) {
    return Config.get(k);
  }

  change(option) {
    return () => {
      var s = this.state;
      var v = !s[option];
      s[option] = v;
      this.set(option, v);
      this.setState(s);
    }
  }

  componentDidMount() {
  }

  render() {
    return(
      <div className="Settings">
        <div className="Section">
          Settings
        </div>
        <div className="Options">
          <Checkbox
            label="Instant Upload"
            desc="When adding a file, the upload will start immediately"
            checked={this.state.instant}
            onClick={this.change('instant')}/>
          <Checkbox
            label="Another Checkbox"
            checked={this.state.other}
            onClick={this.change('other')}/>
        </div>
      </div>
    );
  }
};

export default Settings;
