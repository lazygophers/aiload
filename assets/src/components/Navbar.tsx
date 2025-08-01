import React from 'react';
import './Navbar.css';

const Navbar: React.FC = () => {
  return (
    <nav className="navbar">
      <div className="navbar-logo">
        <a href="/">AI Load</a>
      </div>
      <div className="navbar-links">
        <a href="#features">Features</a>
        <a href="#services">Services</a>
        <a href="#about">About</a>
      </div>
      <div className="navbar-actions">
        <a href="https://github.com/aiload" target="_blank" rel="noopener noreferrer" className="github-button">
          GitHub
        </a>
      </div>
    </nav>
  );
};

export default Navbar;