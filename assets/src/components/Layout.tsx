import React from 'react';
import { Layout as AntLayout } from 'antd';
import { LayoutProvider, useLayout } from './LayoutContext';
import Sidebar from './Sidebar';
import Header from './Header';

const { Content } = AntLayout;

/**
 * 主布局组件
 * 包含侧边栏和主内容区的控制台布局
 */
const LayoutContent: React.FC<{ children?: React.ReactNode }> = ({ children }) => {
	const { sidebarCollapsed } = useLayout();

	// 根据侧边栏状态计算主内容区的左边距
	const contentMarginLeft = sidebarCollapsed ? 80 : 250;

	return (
		<AntLayout style={{ minHeight: '100vh' }}>
			{/* 顶部导航栏组件 */}
			<Header pageTitle="控制台" />
			
			{/* 侧边栏组件 */}
			<Sidebar />

			{/* 主内容区 - 占据剩余空间 */}
			<Content
				style={{
					marginLeft: contentMarginLeft,
					marginTop: '64px', // 为固定定位的Header留出空间
					padding: '24px',
					minHeight: 'calc(100vh - 64px)', // 减去Header高度
					backgroundColor: '#f5f5f5',
					transition: 'margin-left 0.3s ease, margin-top 0.3s ease',
				}}
			>
				{children}
			</Content>
		</AntLayout>
	);
};

/**
 * 布局组件包装器
 * 提供布局状态管理
 */
const Layout: React.FC<{ children?: React.ReactNode }> = ({ children }) => {
	return (
		<LayoutProvider>
			<LayoutContent>{children}</LayoutContent>
		</LayoutProvider>
	);
};

export default Layout;
