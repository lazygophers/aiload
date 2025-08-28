import React from 'react';
import { Typography, Card, Row, Col, Button } from 'antd';
import { useTranslation } from 'react-i18next';

const { Title, Paragraph } = Typography;

/**
 * 首页欢迎组件
 * 展示系统介绍和快速导航
 */
const Home: React.FC = () => {
  const { t } = useTranslation();

  return (
    <div
      style={{
        backgroundColor: '#fff',
        padding: '24px',
        borderRadius: '8px',
        minHeight: 'calc(100vh - 112px)',
        boxShadow: '0 1px 2px rgba(0, 0, 0, 0.03)',
      }}
    >
      <Typography>
        <Title level={2}>{t('welcome.title', '欢迎使用 AI Load 系统')}</Title>
        <Paragraph>
          {t('welcome.description', '这是一个强大的 AI 模型加载和管理系统，支持多种 AI 服务的集成和管理。')}
        </Paragraph>
      </Typography>

      <Row gutter={[16, 16]} style={{ marginTop: '32px' }}>
        <Col xs={24} sm={12} md={8}>
          <Card
            title="模型管理"
            bordered={false}
            hoverable
            style={{ height: '100%' }}
          >
            <Paragraph>
              管理和配置各种 AI 模型，包括 OpenAI、Gemini、Ollama 等主流服务。
            </Paragraph>
            <Button type="primary" block>
              开始配置
            </Button>
          </Card>
        </Col>
        
        <Col xs={24} sm={12} md={8}>
          <Card
            title="API 文档"
            bordered={false}
            hoverable
            style={{ height: '100%' }}
          >
            <Paragraph>
              查看详细的 API 文档，了解如何使用系统提供的各种接口。
            </Paragraph>
            <Button type="primary" block>
              查看文档
            </Button>
          </Card>
        </Col>
        
        <Col xs={24} sm={12} md={8}>
          <Card
            title="系统监控"
            bordered={false}
            hoverable
            style={{ height: '100%' }}
          >
            <Paragraph>
              实时监控系统状态，查看模型加载情况和性能指标。
            </Paragraph>
            <Button type="primary" block>
              查看状态
            </Button>
          </Card>
        </Col>
      </Row>

      <Card title="快速开始" style={{ marginTop: '32px' }}>
        <Row gutter={[16, 16]}>
          <Col span={24}>
            <Paragraph>
              1. 在侧边栏中选择"模型管理"开始配置您的第一个 AI 模型<br />
              2. 配置完成后，您可以通过 API 调用来使用这些模型<br />
              3. 查看"API 文档"了解详细的接口说明
            </Paragraph>
          </Col>
        </Row>
      </Card>
    </div>
  );
};

export default Home;