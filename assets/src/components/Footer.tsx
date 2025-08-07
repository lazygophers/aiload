import React from 'react';
import { Layout, Space, Dropdown, Menu } from 'antd';
import { GlobalOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';

const { Footer: AntFooter } = Layout;

const languageMap: { [key: string]: string } = {
  en: 'English',
  zh: '简体中文',
  'zh-TW': '繁體中文',
  ja: '日本語',
  ko: '한국어',
  fr: 'Français',
  de: 'Deutsch',
  es: 'Español',
};

const Footer: React.FC = () => {
  const { i18n } = useTranslation();

  const changeLanguage = (lng: string) => {
    i18n.changeLanguage(lng);
  };


  const menuItems = Object.keys(i18n.options.resources || {}).map((lng) => ({
    key: lng,
    label: languageMap[lng] || lng,
    onClick: () => changeLanguage(lng),
  }));

  return (
    <AntFooter style={{ textAlign: 'center' }}>
      Copyright © 2025 Your Company Name. All Rights Reserved.
      <br />
      <Space>
        <a href="#">关于我们</a>
        <a href="#">联系方式</a>
        <a href="#">隐私政策</a>
        <Dropdown menu={{ items: menuItems }}>
          <a onClick={(e) => e.preventDefault()}>
            <GlobalOutlined /> {languageMap[i18n.language] || '语言'}
          </a>
        </Dropdown>
      </Space>
    </AntFooter>
  );
};

export default Footer;