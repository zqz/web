import React, { Component } from 'react';
import { Route, BrowserRouter } from 'react-router-dom'
import Uploader from './Uploader.js';
import FileList from './FileList.js';
import FileView from './FileView.js';
import Header from './Header.js';
import Footer from './Footer.js';
import './App.css';

class App extends Component {
  render() {
    return (
      <BrowserRouter>
      <div className="App">
        <Header />
        <Uploader/>
        <Route exact path="/" component={FileList}/>
        <Route path="/file/:id" component={FileView}/>
        <div className="Spacing"></div>
        <Footer />
      </div>
      </BrowserRouter>
    );
  }
}

export default App;
