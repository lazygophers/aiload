import React, { useState, useEffect } from 'react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { Typography, Anchor } from 'antd';
import type { Components } from 'react-markdown';
import './MarkdownRenderer.css';

const { Text } = Typography;
const { Link } = Anchor;

interface MarkdownRendererProps {
  content: string;
  darkMode?: boolean;
  onHeadingClick?: (anchor: string) => void;
}

interface TocItem {
  anchor: string;
  level: number;
  text: string;
}

// 动态导入 SyntaxHighlighter
let SyntaxHighlighter: any = null;
let oneLight: any = null;
let oneDark: any = null;

const loadSyntaxHighlighter = async () => {
  if (!SyntaxHighlighter) {
    try {
      const module = await import('react-syntax-highlighter');
      SyntaxHighlighter = module.Prism;
      const styles = await import('react-syntax-highlighter/dist/esm/styles/prism');
      oneLight = styles.oneLight;
      oneDark = styles.oneDark;
    } catch (error) {
      console.warn('Failed to load react-syntax-highlighter, falling back to pre tag');
    }
  }
};

/**
 * Markdown渲染器组件
 * 支持代码高亮、TOC生成和标题锚点
 */
const MarkdownRenderer: React.FC<MarkdownRendererProps> & {
  generateToc: (content: string) => TocItem[];
} = ({ content, darkMode = false, onHeadingClick }) => {
  const [highlighterLoaded, setHighlighterLoaded] = useState(false);

  useEffect(() => {
    loadSyntaxHighlighter().then(() => setHighlighterLoaded(true));
  }, []);

  // 生成标题ID
  const generateHeadingId = (text: string) => {
    return text
      .toLowerCase()
      .replace(/[^\w\s-]/g, '')
      .replace(/\s+/g, '-')
      .replace(/-+/g, '-')
      .trim();
  };

  // 自定义组件渲染
  const components: Components = {
    // 代码块渲染
    code(props: any) {
      const { inline, className, children, ...rest } = props;
      const match = /language-(\w+)/.exec(className || '');
      const language = match ? match[1] : '';
      
      if (!inline && language && SyntaxHighlighter && highlighterLoaded) {
        return (
          <SyntaxHighlighter
            style={darkMode ? oneDark : oneLight}
            language={language}
            PreTag="div"
            className="code-block"
            {...rest}
          >
            {String(children).replace(/\n$/, '')}
          </SyntaxHighlighter>
        );
      }
      
      return (
        <code className={className} {...rest}>
          {children}
        </code>
      );
    },
    
    // 标题渲染，添加锚点
    h1({ children }) {
      const text = children?.toString() || '';
      const id = generateHeadingId(text);
      return (
        <h1 id={id} className="heading-anchor">
          <Link 
            href={`#${id}`} 
            title={text}
          >
            {text}
          </Link>
        </h1>
      );
    },
    h2({ children }) {
      const text = children?.toString() || '';
      const id = generateHeadingId(text);
      return (
        <h2 id={id} className="heading-anchor">
          <Link 
            href={`#${id}`} 
            title={text}
          >
            {text}
          </Link>
        </h2>
      );
    },
    h3({ children }) {
      const text = children?.toString() || '';
      const id = generateHeadingId(text);
      return (
        <h3 id={id} className="heading-anchor">
          <Link 
            href={`#${id}`} 
            title={text}
          >
            {text}
          </Link>
        </h3>
      );
    },
    h4({ children }) {
      const text = children?.toString() || '';
      const id = generateHeadingId(text);
      return (
        <h4 id={id} className="heading-anchor">
          <Link 
            href={`#${id}`} 
            title={text}
          >
            {text}
          </Link>
        </h4>
      );
    },
    h5({ children }) {
      const text = children?.toString() || '';
      const id = generateHeadingId(text);
      return (
        <h5 id={id} className="heading-anchor">
          <Link 
            href={`#${id}`} 
            title={text}
          >
            {text}
          </Link>
        </h5>
      );
    },
    h6({ children }) {
      const text = children?.toString() || '';
      const id = generateHeadingId(text);
      return (
        <h6 id={id} className="heading-anchor">
          <Link 
            href={`#${id}`} 
            title={text}
          >
            {text}
          </Link>
        </h6>
      );
    },
    
    // 表格样式
    table({ children }) {
      return (
        <div style={{ overflowX: 'auto' }}>
          <table className="markdown-table">{children}</table>
        </div>
      );
    },
    
    // 链接在新标签页打开
    a({ children, href }) {
      return (
        <a href={href} target="_blank" rel="noopener noreferrer">
          {children}
        </a>
      );
    },
  };

  return (
    <div className="markdown-renderer">
      <ReactMarkdown
        remarkPlugins={[remarkGfm]}
        components={components}
      >
        {content}
      </ReactMarkdown>
      
    </div>
  );
};

// 静态方法：生成TOC
MarkdownRenderer.generateToc = (content: string): TocItem[] => {
  const toc: TocItem[] = [];
  const lines = content.split('\n');
  
  lines.forEach(line => {
    const match = line.match(/^(#{1,6})\s+(.+)$/);
    if (match) {
      const level = match[1].length;
      const text = match[2].trim();
      const anchor = text
        .toLowerCase()
        .replace(/[^\w\s-]/g, '')
        .replace(/\s+/g, '-')
        .replace(/-+/g, '-')
        .trim();
      
      toc.push({ anchor, level, text });
    }
  });
  
  return toc;
};

export default MarkdownRenderer;