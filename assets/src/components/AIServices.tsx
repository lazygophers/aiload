import React from 'react';
import './AIServices.css';

import openaiLogo from '../assets/openai.svg';
import geminiLogo from '../assets/gemini.svg';
import anthropicLogo from '../assets/anthropic.svg';
import mistralaiLogo from '../assets/mistralai.svg';
import metaLogo from '../assets/meta.svg';
import cohereLogo from '../assets/cohere.svg';

const logos = [
  { src: openaiLogo, alt: 'OpenAI', url: 'https://openai.com' },
  { src: geminiLogo, alt: 'Google Gemini', url: 'https://deepmind.google/technologies/gemini/' },
  { src: anthropicLogo, alt: 'Anthropic', url: 'https://www.anthropic.com' },
  { src: mistralaiLogo, alt: 'Mistral AI', url: 'https://mistral.ai' },
  { src: metaLogo, alt: 'Meta Llama', url: 'https://llama.meta.com/' },
  { src: cohereLogo, alt: 'Cohere', url: 'https://cohere.com' },
];

const AIServices: React.FC = () => {
  // Duplicate logos for a seamless scroll effect
  const extendedLogos = [...logos, ...logos];

  return (
    <section id="services" className="ai-services-section">
      <div className="container">
        <h2 className="section-title">Supported AI Services</h2>
        <div className="scroller">
          <div className="scroller-inner">
            {extendedLogos.map((logo, index) => (
              <a href={logo.url} target="_blank" rel="noopener noreferrer" key={index}>
                <img src={logo.src} alt={logo.alt} className="logo-item" />
              </a>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
};

export default AIServices;