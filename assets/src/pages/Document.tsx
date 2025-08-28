import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Card, Row, Col, Button, Typography, Empty, Anchor, BackTop, Tag, Skeleton } from 'antd';
import { ArrowLeftOutlined, ReadOutlined, ClockCircleOutlined, CalendarOutlined, FolderOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import MarkdownRenderer from '../components/MarkdownRenderer';
import { documentsMeta as documents, getDocumentMetaById } from '../data/docs';
import { documents as documentContents } from '../pages/docs';

const { Title, Paragraph } = Typography;
const { Link } = Anchor;

interface DocumentProps {}

/**
 * 文档详情页面
 * 根据路由参数显示对应文档的内容，支持目录导航和返回列表
 */
const Document: React.FC<DocumentProps> = () => {
  const { docId } = useParams<{ docId: string }>();
  const navigate = useNavigate();
  const { t } = useTranslation();
  
  // 状态管理
  const [documentContent, setDocumentContent] = useState<string>('');
  const [documentMeta, setDocumentMeta] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [activeAnchor, setActiveAnchor] = useState<string>('');
  const [toc, setToc] = useState<any[]>([]);
  
  // 获取文档内容
  useEffect(() => {
    if (docId) {
      setLoading(true);
      
      // 查找文档元数据
      const meta = getDocumentMetaById(docId);
      setDocumentMeta(meta);
      
      // 查找文档内容
      const document = documentContents.find((doc: any) => doc.name === meta?.name);
      if (document) {
        setDocumentContent(document.content);
        
        // 生成目录
        const tocItems = MarkdownRenderer.generateToc(document.content);
        setToc(tocItems);
      }
      
      setLoading(false);
    }
  }, [docId]);
  
  // 处理返回文档列表
  const handleBackToList = () => {
    navigate('/docs');
  };
  
  // 处理锚点点击
  const handleAnchorClick = (anchor: string) => {
    setActiveAnchor(anchor);
    const element = document.getElementById(anchor);
    if (element) {
      element.scrollIntoView({ behavior: 'smooth' });
      // 更新URL哈希
      window.history.pushState(null, '', `#${anchor}`);
    }
  };
  
  // 监听滚动事件，更新活动锚点
  useEffect(() => {
    const handleScroll = () => {
      const headingElements = document.querySelectorAll('h1[id], h2[id], h3[id], h4[id], h5[id], h6[id]');
      
      for (let i = headingElements.length - 1; i >= 0; i--) {
        const element = headingElements[i] as HTMLElement;
        const rect = element.getBoundingClientRect();
        
        if (rect.top <= 100) {
          setActiveAnchor(element.id);
          break;
        }
      }
    };
    
    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);
  
  // 初始化时检查URL哈希
  useEffect(() => {
    const hash = window.location.hash.slice(1);
    if (hash) {
      setActiveAnchor(hash);
      setTimeout(() => {
        const element = document.getElementById(hash);
        if (element) {
          element.scrollIntoView({ behavior: 'smooth' });
        }
      }, 100);
    }
  }, [documentContent]);
  
  if (loading) {
    return (
      <div style={{ padding: '24px', textAlign: 'center' }}>
        <Paragraph>加载中...</Paragraph>
      </div>
    );
  }
  
  if (!documentMeta || !documentContent) {
    return (
      <div style={{ padding: '24px' }}>
        <Empty description="文档不存在" />
        <Button type="primary" icon={<ArrowLeftOutlined />} onClick={handleBackToList}>
          返回文档列表
        </Button>
      </div>
    );
  }
  
  return (
    <div className="document-page-container">
      {/* 返回按钮和文档标题 */}
      <Row gutter={[16, 16]} className="document-header">
        <Col span={24}>
          <Button
            icon={<ArrowLeftOutlined />}
            onClick={handleBackToList}
            className="back-button"
          >
            返回文档列表
          </Button>
          
          <Card className="document-info-card">
            <div className="document-header-content">
              <div className="document-icon-wrapper">
                <ReadOutlined className="document-icon" />
              </div>
              <div className="document-title-wrapper">
                <Title level={2} className="document-title">
                  {t(documentMeta.titleKey)}
                </Title>
                <Paragraph type="secondary" className="document-description">
                  {t(documentMeta.descriptionKey)}
                </Paragraph>
              </div>
            </div>
            
            <div className="document-meta">
              <span className="meta-item">
                <ClockCircleOutlined /> {documentMeta.readingTime} 分钟阅读
              </span>
              <span className="meta-item">
                <CalendarOutlined /> {new Date(documentMeta.lastUpdated).toLocaleDateString()}
              </span>
              <span className="meta-item">
                <FolderOutlined /> {t(`docs.categories.${documentMeta.category}`)}
              </span>
            </div>
            
            {/* 标签列表 */}
            <div className="document-tags">
              {documentMeta.tags.map((tag: string) => (
                <Tag key={tag} className="document-tag">
                  {tag}
                </Tag>
              ))}
            </div>
          </Card>
        </Col>
      </Row>
      
      {/* 文档内容和目录 */}
      <Row gutter={[16, 16]} className="document-content-wrapper">
        {/* 主内容区 */}
        <Col xs={24} sm={24} md={18}>
          <Card className="document-content-card">
            <MarkdownRenderer
              content={documentContent}
              onHeadingClick={handleAnchorClick}
            />
          </Card>
        </Col>
        
        {/* 目录侧边栏 */}
        <Col xs={0} sm={0} md={6}>
          <Card
            title="目录"
            size="small"
            className="toc-card"
          >
            <Anchor
              affix={false}
              getCurrentAnchor={() => activeAnchor}
              items={toc.map(item => ({
                key: item.anchor,
                href: `#${item.anchor}`,
                title: item.text,
                level: item.level,
              }))}
              className="custom-anchor"
            />
          </Card>
        </Col>
      </Row>
      
      {/* 返回顶部按钮 */}
      <BackTop />
      
      {/* 全局样式 */}
      <style>{`
        .document-page-container {
          max-width: 1200px;
          margin: 0 auto;
          padding: 24px;
          background-color: #f5f5f5;
          min-height: calc(100vh - 64px);
        }
        
        .document-header {
          margin-bottom: 24px;
        }
        
        .back-button {
          margin-bottom: 16px;
          border-radius: 6px;
          height: 40px;
          display: flex;
          align-items: center;
          font-weight: 500;
        }
        
        .document-info-card {
          border-radius: 8px;
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
          border: 1px solid #f0f0f0;
          overflow: hidden;
        }
        
        .document-header-content {
          display: flex;
          align-items: flex-start;
          margin-bottom: 16px;
        }
        
        .document-icon-wrapper {
          width: 64px;
          height: 64px;
          background: linear-gradient(135deg, #1890ff 0%, #36cfc9 100%);
          border-radius: 12px;
          display: flex;
          align-items: center;
          justify-content: center;
          margin-right: 16px;
          flex-shrink: 0;
        }
        
        .document-icon {
          font-size: 32px;
          color: rgba(255, 255, 255, 0.9);
        }
        
        .document-title-wrapper {
          flex: 1;
        }
        
        .document-title {
          margin: 0 !important;
          color: rgba(0, 0, 0, 0.85) !important;
          font-weight: 600 !important;
        }
        
        .document-description {
          margin: 0 !important;
          font-size: 16px;
          line-height: 1.6;
        }
        
        .document-meta {
          display: flex;
          flex-wrap: wrap;
          gap: 24px;
          margin-bottom: 16px;
          color: rgba(0, 0, 0, 0.45);
          font-size: 14px;
        }
        
        .meta-item {
          display: flex;
          align-items: center;
          gap: 6px;
        }
        
        .document-tags {
          display: flex;
          flex-wrap: wrap;
          gap: 8px;
        }
        
        .document-tag {
          font-size: 12px;
          padding: 4px 12px;
          border-radius: 12px;
          background-color: #f0f5ff;
          border: 1px solid #adc6ff;
          color: #1890ff;
          margin: 0;
        }
        
        .document-content-wrapper {
          margin-bottom: 24px;
        }
        
        .document-content-card {
          border-radius: 8px;
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
          border: 1px solid #f0f0f0;
        }
        
        .document-content-card .ant-card-body {
          padding: 32px;
        }
        
        .toc-card {
          border-radius: 8px;
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
          border: 1px solid #f0f0f0;
          position: sticky;
          top: 88px;
        }
        
        .toc-card .ant-card-head {
          border-bottom: 1px solid #f0f0f0;
          background-color: #fafafa;
        }
        
        .toc-card .ant-card-head-title {
          font-weight: 600;
          color: rgba(0, 0, 0, 0.85);
        }
        
        .custom-anchor {
          padding-left: 0 !important;
        }
        
        .custom-anchor .ant-anchor-link {
          padding: 4px 0 !important;
          line-height: 1.8;
        }
        
        .custom-anchor .ant-anchor-link-title {
          font-size: 14px !important;
          color: rgba(0, 0, 0, 0.65) !important;
          transition: color 0.3s;
        }
        
        .custom-anchor .ant-anchor-link-active > .ant-anchor-link-title {
          color: #1890ff !important;
          font-weight: 500;
        }
        
        .custom-anchor .ant-anchor-link:hover .ant-anchor-link-title {
          color: #40a9ff !important;
        }
        
        /* 目录层级缩进 */
        .custom-anchor .ant-anchor-link:nth-child(1) { padding-left: 0 !important; }
        .custom-anchor .ant-anchor-link:nth-child(2) { padding-left: 16px !important; }
        .custom-anchor .ant-anchor-link:nth-child(3) { padding-left: 32px !important; }
        .custom-anchor .ant-anchor-link:nth-child(4) { padding-left: 48px !important; }
        .custom-anchor .ant-anchor-link:nth-child(5) { padding-left: 64px !important; }
        .custom-anchor .ant-anchor-link:nth-child(6) { padding-left: 80px !important; }
        
        @media (max-width: 768px) {
          .document-page-container {
            padding: 16px;
          }
          
          .document-header-content {
            flex-direction: column;
            align-items: center;
            text-align: center;
          }
          
          .document-icon-wrapper {
            margin-right: 0;
            margin-bottom: 16px;
          }
          
          .document-meta {
            justify-content: center;
          }
          
          .document-content-card .ant-card-body {
            padding: 16px;
          }
          
          .toc-card {
            position: static;
            margin-top: 16px;
          }
        }
      `}</style>
    </div>
  );
};

export default Document;