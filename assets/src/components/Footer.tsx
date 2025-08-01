import React from 'react';
import './Footer.css';

const Footer: React.FC = () => {

  const scrollToTop = () => {
    window.scrollTo({
      top: 0,
      behavior: 'smooth'
    });
  };

  return (
    <footer className="footer">
      <div className="container footer-container">
        <div className="footer-section about">
          <h4 className="footer-title">AI Load</h4>
          <p>The ultimate load testing tool for AI applications. Built for developers, by developers.</p>
        </div>
        <div className="footer-section links">
          <h4 className="footer-title">Quick Links</h4>
          <ul>
            <li><a href="#features">Features</a></li>
            <li><a href="/docs">Documentation</a></li>
            <li><a href="https://github.com/aiload/issues">Report a Bug</a></li>
          </ul>
        </div>
        <div className="footer-section legal">
          <h4 className="footer-title">Legal</h4>
          <ul>
            <li><a href="/privacy">Privacy Policy</a></li>
            <li><a href="/terms">Terms of Service</a></li>
          </ul>
        </div>
      </div>
      <div className="footer-bottom">
        <p>&copy; {new Date().getFullYear()} AI Load. All rights reserved.</p>
        <button onClick={scrollToTop} className="back-to-top">
          Back to Top &uarr;
        </button>
      </div>
    </footer>
  );
};

export default Footer;