import React from 'react';
import { Layout, Menu, Button, theme } from 'antd';
import { useTranslation } from 'react-i18next';
import { useNavigation } from '../hooks/useNavigation';
import './Sidebar.css';

const { Sider } = Layout;

/**
 * 侧边栏组件
 * 提供可折叠的导航菜单，支持国际化、动态高亮和文档子菜单
 */
const Sidebar: React.FC = () => {
  const { t } = useTranslation();
  const {
    menuItems,
    selectedKeys,
    collapseIcon,
    sidebarCollapsed,
    toggleSidebar,
    handleMenuClick,
  } = useNavigation();
  
  // 使用 Ant Design 的主题 token
  const {
    token: { colorBgContainer },
  } = theme.useToken();

  return (
    <Sider
      trigger={null}
      collapsible
      collapsed={sidebarCollapsed}
      theme="light"
      width={250}
      collapsedWidth={80}
      style={{
        overflow: 'auto',
        height: '100vh',
        position: 'fixed',
        left: 0,
        top: 0,
        bottom: 0,
        borderRight: '1px solid #f0f0f0',
        background: colorBgContainer,
      }}
    >
      {/* 折叠/展开按钮 */}
      <div className="sidebar-trigger">
        <Button
          type="text"
          icon={collapseIcon}
          onClick={toggleSidebar}
          style={{
            fontSize: '16px',
            width: 64,
            height: 64,
          }}
        />
      </div>
      
      {/* 导航菜单 */}
      <Menu
        mode="inline"
        selectedKeys={selectedKeys}
        items={menuItems}
        className="sidebar-menu"
        onClick={({ key }) => handleMenuClick(key)}
        // 支持子菜单展开
        openKeys={sidebarCollapsed ? [] : undefined}
      />
    </Sider>
  );
};

export default Sidebar;