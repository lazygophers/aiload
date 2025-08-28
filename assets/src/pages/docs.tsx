import React, { useState, useMemo } from 'react';
import { Typography, Card, Row, Col, Tag, Button, Input, Select, Space, Empty, Tooltip } from 'antd';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import { SearchOutlined, BookOutlined, ClockCircleOutlined, TagOutlined, FolderOutlined } from '@ant-design/icons';
import { documentsMeta, getDocumentsByCategory, getAllCategories, getAllTags } from '../data/docs';

const { Title, Paragraph, Text } = Typography;
const { Search } = Input;
const { Option } = Select;

/**
 * 文档列表页面
 * 展示所有可用文档，支持分类筛选和搜索
 */
const Docs: React.FC = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  
  // 状态管理
  const [searchText, setSearchText] = useState('');
  const [selectedCategory, setSelectedCategory] = useState<string>('all');
  const [selectedTag, setSelectedTag] = useState<string>('all');
  
  // 获取所有分类和标签
  const categories = getAllCategories();
  const tags = getAllTags();
  
  // 过滤和搜索文档
  const filteredDocuments = useMemo(() => {
    let docs = documentsMeta;
    
    // 按分类筛选
    if (selectedCategory !== 'all') {
      docs = getDocumentsByCategory(selectedCategory as any);
    }
    
    // 按标签筛选
    if (selectedTag !== 'all') {
      docs = docs.filter(doc => doc.tags.includes(selectedTag));
    }
    
    // 按搜索文本筛选
    if (searchText) {
      const searchLower = searchText.toLowerCase();
      docs = docs.filter(doc => 
        t(doc.titleKey).toLowerCase().includes(searchLower) ||
        t(doc.descriptionKey).toLowerCase().includes(searchLower) ||
        doc.tags.some(tag => tag.toLowerCase().includes(searchLower))
      );
    }
    
    return docs;
  }, [searchText, selectedCategory, selectedTag, t]);
  
  // 分组文档按分类
  const groupedDocuments = useMemo(() => {
    const groups: Record<string, typeof documentsMeta> = {};
    
    filteredDocuments.forEach(doc => {
      if (!groups[doc.category]) {
        groups[doc.category] = [];
      }
      groups[doc.category].push(doc);
    });
    
    return groups;
  }, [filteredDocuments]);
  
  // 处理文档卡片点击
  const handleDocumentClick = (docId: string) => {
    // 导航到文档详情页面
    navigate(`/docs/${docId}`);
  };
  
  // 获取分类的显示名称
  const getCategoryDisplayName = (category: string) => {
    const categoryMap: Record<string, string> = {
      'getting-started': t('docs.categories.getting-started'),
      'core-concepts': t('docs.categories.core-concepts'),
      'api-reference': t('docs.categories.api-reference'),
      'guides': t('docs.categories.guides'),
      'development': t('docs.categories.development'),
    };
    return categoryMap[category] || category;
  };
  
  // 获取分类的颜色
  const getCategoryColor = (category: string) => {
    const colorMap: Record<string, string> = {
      'getting-started': 'green',
      'core-concepts': 'blue',
      'api-reference': 'purple',
      'guides': 'orange',
      'development': 'red',
    };
    return colorMap[category] || 'default';
  };
  
  return (
    <>
      {/* 页面内容 */}
      <div className="docs-page-container">
        {/* 页面标题 */}
        <Typography>
          <Title level={2}>{t('docs.page.title')}</Title>
          <Paragraph type="secondary">
            {t('docs.page.description')}
          </Paragraph>
        </Typography>
        
        {/* 搜索和筛选区域 */}
        <Card className="filter-card" style={{ marginBottom: '24px' }}>
          <Space wrap style={{ width: '100%' }}>
            <Search
              placeholder={t('docs.search.placeholder')}
              allowClear
              style={{ width: 300 }}
              value={searchText}
              onChange={(e) => setSearchText(e.target.value)}
              prefix={<SearchOutlined />}
              size="large"
            />
            
            <Select
              placeholder={t('docs.filter.category')}
              style={{ width: 200 }}
              value={selectedCategory}
              onChange={setSelectedCategory}
              allowClear
              size="large"
              suffixIcon={<FolderOutlined />}
            >
              <Option value="all">{t('docs.filter.all-categories')}</Option>
              {categories.map(category => (
                <Option key={category} value={category}>
                  {getCategoryDisplayName(category)}
                </Option>
              ))}
            </Select>
            
            <Select
              placeholder={t('docs.filter.tag')}
              style={{ width: 200 }}
              value={selectedTag}
              onChange={setSelectedTag}
              allowClear
              size="large"
              suffixIcon={<TagOutlined />}
            >
              <Option value="all">{t('docs.filter.all-tags')}</Option>
              {tags.map(tag => (
                <Option key={tag} value={tag}>
                  #{tag}
                </Option>
              ))}
            </Select>
            
            <Button
              type="default"
              onClick={() => {
                setSearchText('');
                setSelectedCategory('all');
                setSelectedTag('all');
              }}
              style={{ marginLeft: '8px' }}
            >
              {t('docs.filter.clear')}
            </Button>
          </Space>
        </Card>
        
        {/* 文档列表 */}
        {Object.keys(groupedDocuments).length === 0 ? (
          <Empty 
            description={t('docs.no-results')} 
            image={Empty.PRESENTED_IMAGE_SIMPLE}
            style={{ padding: '60px 0' }}
          />
        ) : (
          Object.entries(groupedDocuments).map(([category, docs]) => (
            <div key={category} className="document-category">
              <Title level={3} className="category-title">
                <Tag color={getCategoryColor(category)} className="category-tag">
                  {getCategoryDisplayName(category)}
                </Tag>
                <span className="doc-count">({docs.length})</span>
              </Title>
              
              <Row gutter={[20, 20]}>
                {docs.map(doc => (
                  <Col xs={24} sm={12} md={8} lg={6} key={doc.id}>
                    <Card
                      className="document-card"
                      hoverable
                      onClick={() => handleDocumentClick(doc.id)}
                      cover={
                        <div className="card-cover">
                          <BookOutlined className="card-icon" />
                        </div>
                      }
                    >
                      <Card.Meta
                        title={
                          <Tooltip title={t(doc.titleKey)}>
                            <Text className="doc-title">
                              {t(doc.titleKey)}
                            </Text>
                          </Tooltip>
                        }
                        description={
                          <div className="doc-meta">
                            <Paragraph
                              type="secondary"
                              ellipsis={{ rows: 2 }}
                              className="doc-description"
                            >
                              {t(doc.descriptionKey)}
                            </Paragraph>
                            <div className="doc-tags">
                              {doc.tags.slice(0, 3).map(tag => (
                                <Tag key={tag} className="doc-tag">
                                  {tag}
                                </Tag>
                              ))}
                              {doc.tags.length > 3 && (
                                <Tag className="doc-tag more-tag">
                                  +{doc.tags.length - 3}
                                </Tag>
                              )}
                            </div>
                            <div className="doc-info">
                              <span className="reading-time">
                                <ClockCircleOutlined /> {doc.readingTime} {t('docs.minutes')}
                              </span>
                              <span className="update-date">
                                {new Date(doc.lastUpdated).toLocaleDateString()}
                              </span>
                            </div>
                          </div>
                        }
                      />
                    </Card>
                  </Col>
                ))}
              </Row>
            </div>
          ))
        )}
      </div>
      
      {/* 全局样式 */}
      <style>{`
        .docs-page-container {
          background-color: #f5f5f5;
          padding: 24px;
          min-height: calc(100vh - 64px);
        }
        
        .filter-card {
          border-radius: 8px;
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
          border: 1px solid #f0f0f0;
        }
        
        .document-category {
          margin-bottom: 40px;
        }
        
        .category-title {
          margin-bottom: 20px !important;
          color: rgba(0, 0, 0, 0.85);
          display: flex;
          align-items: center;
          gap: 12px;
        }
        
        .category-tag {
          font-size: 14px;
          padding: 4px 12px;
          border-radius: 4px;
        }
        
        .doc-count {
          color: rgba(0, 0, 0, 0.45);
          font-size: 16px;
          font-weight: normal;
        }
        
        .document-card {
          height: 100%;
          border-radius: 8px;
          overflow: hidden;
          transition: all 0.3s ease;
          border: 1px solid #f0f0f0;
        }
        
        .document-card:hover {
          transform: translateY(-4px);
          box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
          border-color: #1890ff;
        }
        
        .card-cover {
          height: 120px;
          background: linear-gradient(135deg, #1890ff 0%, #36cfc9 100%);
          display: flex;
          align-items: center;
          justify-content: center;
          position: relative;
          overflow: hidden;
        }
        
        .card-cover::before {
          content: '';
          position: absolute;
          top: -50%;
          left: -50%;
          width: 200%;
          height: 200%;
          background: radial-gradient(circle, rgba(255,255,255,0.1) 0%, transparent 70%);
          animation: shimmer 3s infinite;
        }
        
        @keyframes shimmer {
          0% { transform: rotate(0deg); }
          100% { transform: rotate(360deg); }
        }
        
        .card-icon {
          font-size: 48px;
          color: rgba(255, 255, 255, 0.9);
          z-index: 1;
        }
        
        .doc-title {
          font-size: 16px !important;
          font-weight: 600 !important;
          color: rgba(0, 0, 0, 0.85) !important;
          display: block;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }
        
        .doc-meta {
          display: flex;
          flex-direction: column;
          gap: 12px;
        }
        
        .doc-description {
          margin-bottom: 0 !important;
          font-size: 13px !important;
          line-height: 1.6;
          color: rgba(0, 0, 0, 0.45) !important;
        }
        
        .doc-tags {
          display: flex;
          flex-wrap: wrap;
          gap: 6px;
        }
        
        .doc-tag {
          font-size: 11px !important;
          padding: 2px 8px !important;
          border-radius: 10px !important;
          background-color: #f0f5ff !important;
          border: 1px solid #adc6ff !important;
          color: #1890ff !important;
        }
        
        .more-tag {
          background-color: #f5f5f5 !important;
          border: 1px solid #d9d9d9 !important;
          color: rgba(0, 0, 0, 0.45) !important;
        }
        
        .doc-info {
          display: flex;
          justify-content: space-between;
          align-items: center;
          font-size: 12px;
          color: rgba(0, 0, 0, 0.45);
          margin-top: 8px;
        }
        
        .reading-time,
        .update-date {
          display: flex;
          align-items: center;
          gap: 4px;
        }
        
        @media (max-width: 768px) {
          .docs-page-container {
            padding: 16px;
          }
          
          .document-card:hover {
            transform: none;
          }
          
          .card-cover {
            height: 100px;
          }
          
          .card-icon {
            font-size: 36px;
          }
        }
      `}</style>
    </>
  );
};

export default Docs;