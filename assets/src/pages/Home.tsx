import React, { useState, useEffect } from 'react';
import { Tabs, Spin } from 'antd';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';

const markdownFiles = [
  'README.md',
  'api_spec.md',
  'architecture.md',
  'cache.md',
  'configuration.md',
  'contributing.md',
  'database.md',
  'deployment.md',
  'getting-started.md',
  'introduction.md',
  'llms_api.md',
  'core-modules.md',
];

interface Doc {
  fileName: string;
  content: string;
}

const Home: React.FC = () => {
  const [docs, setDocs] = useState<Doc[]>([]);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    const fetchDocs = async () => {
      try {
        const fetchedDocs = await Promise.all(
          markdownFiles.map(async (file) => {
            const response = await fetch(`/docs/${file}`);
            if (!response.ok) {
              throw new Error(`Failed to fetch ${file}`);
            }
            const content = await response.text();
            return { fileName: file, content };
          })
        );
        setDocs(fetchedDocs);
      } catch (error) {
        console.error("Error fetching markdown files:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchDocs();
  }, []);

  if (loading) {
    return <Spin size="large" style={{ display: 'block', marginTop: '50px' }} />;
  }

  const tabItems = docs.map((doc, index) => ({
    key: String(index),
    label: doc.fileName,
    children: <ReactMarkdown remarkPlugins={[remarkGfm]}>{doc.content}</ReactMarkdown>,
  }));

  return (
    <div style={{ padding: '24px' }}>
      <Tabs defaultActiveKey="0" type="card" items={tabItems} />
    </div>
  );
};

export default Home;