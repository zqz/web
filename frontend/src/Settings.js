import React, { Component } from 'react';
import Checkbox from './Checkbox';
import Config from './Config';
import './Settings.css';

class Settings extends Component {
  constructor(props) {
    super(props);

    this.state = {
      instant: this.get('instant') === 'true',
      dark: this.get('dark') === 'true',
      other: this.get('other') === 'true',
      'filelist-rows': this.get('filelist-rows') !== 'true'
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

      if (option === 'dark') {
        // Config.toggleDark();
      }
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
            label="Night Mode"
            checked={this.state.dark}
            onClick={this.change('dark')}/>
          <Checkbox
            label="Another Checkbox"
            checked={this.state.other}
            onClick={this.change('other')}/>
          <Checkbox
            label="Grid View"
            desc="Show files as a grid of thumbnails on the homepage"
            checked={this.state['filelist-rows']}
            onClick={this.change('filelist-rows')}/>
        </div>
      </div>
    );
  }
};

export default Settings;
