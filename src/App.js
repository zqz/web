import React, { Component } from 'react';
import Uploader from './Uploader.js';
import FileList from './FileList.js';
import './App.css';

class App extends Component {
  render() {
    return (
      <div className="App">
      <div>zqz.ca</div>
        <Uploader/>
        <FileList/>
      </div>
    );
  }
}

export default App;
