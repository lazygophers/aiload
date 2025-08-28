import { useMemo } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useLayout } from '../components/LayoutContext';
import { documentsMeta, getDocumentsByCategory } from '../data/docs';
import type { MenuProps } from 'antd';
import {
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  FileTextOutlined,
  ApiOutlined,
  QuestionCircleOutlined,
  GlobalOutlined,
  SettingOutlined,
  BookOutlined,
} from '@ant-design/icons';

export type MenuItem = Required<MenuProps>['items'][number];

/**
 * 导航 Hook
 * 提供动态菜单项生成和导航功能
 */
export const useNavigation = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const { t } = useTranslation();
  const { sidebarCollapsed, setSidebarCollapsed } = useLayout();

  // 获取当前选中的菜单项
  const selectedKeys = useMemo(() => {
    const pathname = location.pathname;
    
    // 如果是文档详情页，高亮 docs
    if (pathname.startsWith('/docs/') && pathname !== '/docs') {
      return ['docs'];
    }
    
    // 其他情况直接匹配路径
    return [pathname === '/' ? '/home' : pathname];
  }, [location.pathname]);

  // 生成文档子菜单 - 显示所有文档而不仅仅是 getting-started
  const documentMenuItems = useMemo((): MenuItem[] => {
    return documentsMeta.map(doc => ({
      key: `/docs/${doc.id}`,
      icon: <BookOutlined />,
      label: t(doc.titleKey),
    }));
  }, [t]);

  // 生成完整的菜单项
  const menuItems = useMemo((): MenuItem[] => {
    const baseItems: MenuItem[] = [
      {
        key: 'docs',
        icon: <FileTextOutlined />,
        label: t('sidebar.docs'),
        children: [
          {
            key: 'docs-group',
            type: 'group' as const,
            label: t('sidebar.documents'),
            children: documentMenuItems,
          },
        ],
      },
      {
        key: 'api',
        icon: <ApiOutlined />,
        label: t('sidebar.api'),
      },
      {
        key: 'faq',
        icon: <QuestionCircleOutlined />,
        label: t('sidebar.faq'),
      },
      {
        key: 'i18n',
        icon: <GlobalOutlined />,
        label: t('sidebar.i18n'),
      },
      {
        key: 'settings',
        icon: <SettingOutlined />,
        label: t('sidebar.settings'),
      },
    ];

    return baseItems;
  }, [t, documentMenuItems]);

  // 处理菜单点击
  const handleMenuClick = (key: string) => {
    // 根据key进行导航
    switch (key) {
      case 'docs':
        navigate('/docs');
        break;
      case 'api':
        // API页面的路由
        break;
      case 'faq':
        // FAQ页面的路由
        break;
      case 'i18n':
        // 国际化页面的路由
        break;
      case 'settings':
        // 设置页面的路由
        break;
      default:
        // 如果是文档详情页的路径
        if (key.startsWith('/docs/')) {
          navigate(key);
        }
    }

    // 在移动端自动收起侧边栏
    if (window.innerWidth <= 768) {
      setSidebarCollapsed(true);
    }
  };

  // 切换侧边栏折叠状态
  const toggleSidebar = () => {
    setSidebarCollapsed(!sidebarCollapsed);
  };

  // 获取折叠按钮图标
  const collapseIcon = sidebarCollapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />;

  return {
    menuItems,
    selectedKeys,
    collapseIcon,
    sidebarCollapsed,
    toggleSidebar,
    handleMenuClick,
  };
};