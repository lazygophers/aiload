import React, { useState, useEffect } from 'react';
import { Layout, Menu, Dropdown, Button, Typography, theme } from 'antd';
import {
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  GlobalOutlined,
  CheckOutlined,
} from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useLayout } from './LayoutContext';
import './Header.css';

const { Header: AntHeader } = Layout;
const { Text } = Typography;

/**
 * 语言选项配置
 */
const languageOptions = [
  { key: 'zh', label: '简体中文', flag: '🇨🇳' },
  { key: 'en', label: 'English', flag: '🇺🇸' },
  { key: 'ja', label: '日本語', flag: '🇯🇵' },
  { key: 'ko', label: '한국어', flag: '🇰🇷' },
  { key: 'es', label: 'Español', flag: '🇪🇸' },
  { key: 'fr', label: 'Français', flag: '🇫🇷' },
  { key: 'de', label: 'Deutsch', flag: '🇩🇪' },
  { key: 'zh-TW', label: '繁體中文', flag: '🇹🇼' },
];

/**
 * 顶部导航栏组件
 * 提供页面标题、语言切换和移动端菜单折叠功能
 */
const Header: React.FC<{ pageTitle?: string }> = ({ pageTitle }) => {
  const { sidebarCollapsed, setSidebarCollapsed } = useLayout();
  const { i18n, t } = useTranslation();
  const [currentLang, setCurrentLang] = useState(i18n.language);
  const [isChanging, setIsChanging] = useState(false);

  // 使用 Ant Design 的主题 token
  const {
    token: { colorBgContainer },
  } = theme.useToken();

  // 监听语言变化
  useEffect(() => {
    const handleLanguageChanged = (lng: string) => {
      setCurrentLang(lng);
    };
    
    i18n.on('languageChanged', handleLanguageChanged);
    
    return () => {
      i18n.off('languageChanged', handleLanguageChanged);
    };
  }, [i18n]);

  /**
   * 切换侧边栏状态
   */
  const toggleSidebar = () => {
    setSidebarCollapsed(!sidebarCollapsed);
  };

  /**
   * 切换语言
   * @param lang - 语言代码
   */
  const changeLanguage = (lang: string) => {
    if (lang !== currentLang) {
      setIsChanging(true);
      i18n.changeLanguage(lang);
      
      // 添加过渡动画完成后重置状态
      setTimeout(() => {
        setIsChanging(false);
      }, 300);
    }
  };

  /**
   * 获取当前语言信息
   */
  const getCurrentLanguage = () => {
    return languageOptions.find(option => option.key === currentLang) || languageOptions[0];
  };

  /**
   * 渲染语言切换菜单
   */
  const languageMenu = (
    <Menu
      items={languageOptions.map(option => ({
        key: option.key,
        label: (
          <span className="language-menu-item">
            <span className="language-flag">{option.flag}</span>
            <span className="language-label">{option.label}</span>
            {option.key === currentLang && (
              <CheckOutlined className="language-check" />
            )}
          </span>
        ),
        onClick: () => changeLanguage(option.key),
      }))}
    />
  );

  // 获取当前语言信息
  const currentLanguage = getCurrentLanguage();

  return (
    <AntHeader
      className={`app-header ${isChanging ? 'language-changing' : ''}`}
      style={{ background: colorBgContainer }}
    >
      <div className="header-left">
        {/* 移动端菜单折叠按钮 */}
        <Button
          type="text"
          icon={sidebarCollapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
          onClick={toggleSidebar}
          className="sidebar-trigger"
          aria-label={sidebarCollapsed ? t('header.expandSidebar') : t('header.collapseSidebar')}
        />
        
        {/* 页面标题 */}
        <Text className="page-title" strong>
          {pageTitle || t('header.defaultTitle')}
        </Text>
      </div>

      <div className="header-right">
        {/* 语言切换下拉菜单 */}
        <Dropdown
          overlay={languageMenu}
          placement="bottomRight"
          arrow
          trigger={['click']}
        >
          <Button
            type="text"
            className={`language-switcher ${isChanging ? 'switching' : ''}`}
          >
            <GlobalOutlined />
            <span className="current-language">
              <span className="language-flag">{currentLanguage.flag}</span>
              <span className="language-name">{currentLanguage.label}</span>
            </span>
          </Button>
        </Dropdown>
      </div>
    </AntHeader>
  );
};

export default Header;