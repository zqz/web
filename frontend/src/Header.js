import React, { Component } from 'react';
import { Link } from 'react-router-dom'
import './Header.css';

class Header extends Component {
  render() {
    return (
      <div className="Header">
        <div className="Container">
          <Link to='/' className="Logo">
            zqz.ca
          </Link>
          <div className="Navigation">
            <Link to="/settings">Settings</Link>
            <Link to="/">Files</Link>
          </div>
        </div>
      </div>
    )
  }
};

export default Header;
