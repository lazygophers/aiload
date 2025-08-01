import React from 'react';
import './Features.css';

interface FeatureCardProps {
  icon: string;
  title: string;
  description: string;
}

const FeatureCard: React.FC<FeatureCardProps> = ({ icon, title, description }) => (
  <div className="feature-card">
    <div className="feature-icon">{icon}</div>
    <h3 className="feature-title">{title}</h3>
    <p className="feature-description">{description}</p>
  </div>
);

const Features: React.FC = () => {
  const featuresData = [
    {
      icon: '🚀',
      title: 'High-Concurrency Simulation',
      description: 'Simulate thousands of concurrent users to test your AI application\'s limits.',
    },
    {
      icon: '📊',
      title: 'Real-time Analytics',
      description: 'Get detailed, real-time reports on performance, latency, and error rates.',
    },
    {
      icon: '☁️',
      title: 'Cloud-Native',
      description: 'Deploy and scale your load tests effortlessly on any cloud provider.',
    },
  ];

  return (
    <section id="features" className="features-section">
      <div className="container">
        <h2 className="section-title">Why AI Load?</h2>
        <div className="features-grid">
          {featuresData.map((feature, index) => (
            <FeatureCard key={index} {...feature} />
          ))}
        </div>
      </div>
    </section>
  );
};

export default Features;