import { Outlet, useNavigate } from 'react-router-dom';
import { Layout as AntLayout, Menu, Button, Space, Badge } from 'antd';
import {
  HomeOutlined,
  TeamOutlined,
  ApartmentOutlined,
  UserOutlined,
  BellOutlined,
  CheckSquareOutlined,
  FileDoneOutlined,
  ProjectOutlined,
  TrophyOutlined,
  SettingOutlined,
} from '@ant-design/icons';
import { useDispatch } from 'react-redux';
import { useState, useEffect } from 'react';
import { logout } from '@/store/authSlice';
import { approvalApi } from '@/api/approval';
import { notificationApi } from '@/api/notification';

const { Header, Sider, Content } = AntLayout;

export const Layout = () => {
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const [pendingCount, setPendingCount] = useState(0);
  const [unreadCount, setUnreadCount] = useState(0);

  useEffect(() => {
    const fetchPendingCount = async () => {
      try {
        const response = await approvalApi.listPending();
        setPendingCount(response.data?.length || 0);
      } catch {
        // ignore error
      }
    };
    const fetchUnreadCount = async () => {
      try {
        const response = await notificationApi.list();
        const unread = response.data?.filter((n: { is_read: boolean }) => !n.is_read).length || 0;
        setUnreadCount(unread);
      } catch {
        // ignore error
      }
    };
    fetchPendingCount();
    fetchUnreadCount();
  }, []);

  const menuItems = [
    { key: '/', icon: <HomeOutlined />, label: '工作台' },
    { key: '/companies', icon: <HomeOutlined />, label: '公司管理' },
    { key: '/departments', icon: <ApartmentOutlined />, label: '部门管理' },
    { key: '/employees', icon: <TeamOutlined />, label: '员工管理' },
    { key: '/positions', icon: <TrophyOutlined />, label: '职位管理' },
    { key: '/workflows', icon: <ProjectOutlined />, label: '流程模板' },
    { key: '/tasks/pending', icon: <Badge count={pendingCount} size="small"><CheckSquareOutlined /></Badge>, label: '待审批' },
    { key: '/tasks/handled', icon: <FileDoneOutlined />, label: '已审批' },
    { key: '/notifications', icon: <Badge count={unreadCount} size="small"><BellOutlined /></Badge>, label: '通知' },
    { key: '/system-settings', icon: <SettingOutlined />, label: '系统设置' },
  ];

  const handleLogout = () => {
    dispatch(logout());
    navigate('/login');
  };

  return (
    <AntLayout style={{ minHeight: '100vh' }}>
      <Sider theme="light" width={200}>
        <div style={{ padding: '16px', fontSize: '18px', fontWeight: 'bold', textAlign: 'center' }}>
          Workflow
        </div>
        <Menu
          mode="inline"
          defaultSelectedKeys={['/workflows']}
          items={menuItems}
          onClick={({ key }) => navigate(key)}
        />
      </Sider>
      <AntLayout>
        <Header style={{ background: '#fff', padding: '0 24px', display: 'flex', justifyContent: 'flex-end', alignItems: 'center' }}>
          <Space>
            <UserOutlined /> <span>admin</span>
            <Button onClick={handleLogout}>退出</Button>
          </Space>
        </Header>
        <Content style={{ padding: '24px' }}>
          <Outlet />
        </Content>
      </AntLayout>
    </AntLayout>
  );
};
