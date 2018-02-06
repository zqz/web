import React, { Component } from 'react';
import { Link } from 'react-router-dom'
import './Header.css';

class Header extends Component {
  render() {
    return (
      <div className="Header">
        <div className="Container">
          <div className="Logo">
            zqz.ca
          </div>
          <div className="Navigation">
            <Link to="/">Files</Link>
          </div>
        </div>
      </div>
    )
  }
};

export default Header;
