import React, { Component } from 'react';
import { Route, BrowserRouter } from 'react-router-dom'
import Uploader from './Uploader.js';
import FileList from './FileList.js';
import FileView from './FileView.js';
import './App.css';

class App extends Component {
  render() {
    return (
      <div className="App">
      <div>zqz.ca</div>
        <Uploader/>
        <BrowserRouter>
          <div>
            <Route exact path="/" component={FileList}/>
            <Route path="/file/:id" component={FileView}/>
          </div>
        </BrowserRouter>
      </div>
    );
  }
}

export default App;
