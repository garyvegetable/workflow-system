import { useState, useEffect } from 'react';
import { Card, Row, Col, Statistic, Button, List, Space, Tag, message } from 'antd';
import { useNavigate } from 'react-router-dom';
import { approvalApi } from '@/api/approval';
import { workflowApi } from '@/api/workflow';
import { notificationApi } from '@/api/notification';
import {
  CheckSquareOutlined,
  FileDoneOutlined,
  BellOutlined,
  PlusOutlined,
  ProjectOutlined,
  ClockCircleOutlined,
} from '@ant-design/icons';

interface DashboardStats {
  pendingCount: number
  handledCount: number
  notificationCount: number
  workflowCount: number
}

interface RecentTask {
  id: number
  node_name: string
  instance_title: string
  status: number
  create_time: string
}

interface Notification {
  id: number
  title: string
  is_read: boolean
  created_at: string
}

export const Dashboard = () => {
  const navigate = useNavigate();
  const [stats, setStats] = useState<DashboardStats>({
    pendingCount: 0,
    handledCount: 0,
    notificationCount: 0,
    workflowCount: 0,
  });
  const [recentTasks, setRecentTasks] = useState<RecentTask[]>([]);
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchDashboardData();
  }, []);

  const fetchDashboardData = async () => {
    setLoading(true);
    try {
      // Fetch all data in parallel
      const [pendingRes, handledRes, notifRes, workflowRes] = await Promise.allSettled([
        approvalApi.listPending(),
        approvalApi.listHandled(),
        notificationApi.list(),
        workflowApi.list(),
      ]);

      // Process pending tasks
      if (pendingRes.status === 'fulfilled') {
        const pending = pendingRes.value.data || [];
        setStats(prev => ({ ...prev, pendingCount: pending.length }));
        setRecentTasks(pending.slice(0, 5));
      }

      // Process handled tasks count
      if (handledRes.status === 'fulfilled') {
        setStats(prev => ({ ...prev, handledCount: handledRes.value.data?.length || 0 }));
      }

      // Process notifications
      if (notifRes.status === 'fulfilled') {
        const notifs = notifRes.value.data || [];
        setStats(prev => ({
          ...prev,
          notificationCount: notifs.filter((n: Notification) => !n.is_read).length,
        }));
        setNotifications(notifs.slice(0, 5));
      }

      // Process workflow count
      if (workflowRes.status === 'fulfilled') {
        const workflows = workflowRes.value.data || [];
        setStats(prev => ({
          ...prev,
          workflowCount: workflows.filter((w: any) => w.status === 2).length,
        }));
      }
    } catch (error) {
      message.error('获取数据失败');
    } finally {
      setLoading(false);
    }
  };

  const getStatusTag = (status: number) => {
    switch (status) {
      case 1: return <Tag color="blue">待处理</Tag>;
      case 2: return <Tag color="green">已处理</Tag>;
      case 3: return <Tag color="red">已驳回</Tag>;
      default: return <Tag>未知</Tag>;
    }
  };

  return (
    <div>
      <h1 style={{ marginBottom: 24 }}>工作台</h1>

      {/* 统计卡片 */}
      <Row gutter={16} style={{ marginBottom: 24 }}>
        <Col span={6}>
          <Card loading={loading}>
            <Statistic
              title="待审批任务"
              value={stats.pendingCount}
              prefix={<CheckSquareOutlined style={{ color: '#1890ff' }} />}
              valueStyle={{ color: '#1890ff' }}
            />
            <Button
              type="link"
              onClick={() => navigate('/tasks/pending')}
              style={{ padding: 0, marginTop: 8 }}
            >
              立即处理
            </Button>
          </Card>
        </Col>
        <Col span={6}>
          <Card loading={loading}>
            <Statistic
              title="已审批任务"
              value={stats.handledCount}
              prefix={<FileDoneOutlined style={{ color: '#52c41a' }} />}
              valueStyle={{ color: '#52c41a' }}
            />
            <Button
              type="link"
              onClick={() => navigate('/tasks/handled')}
              style={{ padding: 0, marginTop: 8 }}
            >
              查看历史
            </Button>
          </Card>
        </Col>
        <Col span={6}>
          <Card loading={loading}>
            <Statistic
              title="未读通知"
              value={stats.notificationCount}
              prefix={<BellOutlined style={{ color: '#faad14' }} />}
              valueStyle={{ color: '#faad14' }}
            />
            <Button
              type="link"
              onClick={() => navigate('/notifications')}
              style={{ padding: 0, marginTop: 8 }}
            >
              查看全部
            </Button>
          </Card>
        </Col>
        <Col span={6}>
          <Card loading={loading}>
            <Statistic
              title="可用流程"
              value={stats.workflowCount}
              prefix={<ProjectOutlined style={{ color: '#722ed1' }} />}
              valueStyle={{ color: '#722ed1' }}
            />
            <Button
              type="link"
              onClick={() => navigate('/workflows')}
              style={{ padding: 0, marginTop: 8 }}
            >
              查看流程
            </Button>
          </Card>
        </Col>
      </Row>

      {/* 快捷入口 */}
      <Card title="快捷入口" style={{ marginBottom: 24 }}>
        <Space size="large">
          <Button
            type="primary"
            icon={<PlusOutlined />}
            onClick={() => navigate('/workflows')}
          >
            发起申请
          </Button>
          <Button
            icon={<CheckSquareOutlined />}
            onClick={() => navigate('/tasks/pending')}
          >
            待办任务
          </Button>
          <Button
            icon={<ClockCircleOutlined />}
            onClick={() => navigate('/my-applications')}
          >
            我的申请
          </Button>
        </Space>
      </Card>

      {/* 待审批任务和通知 */}
      <Row gutter={16}>
        <Col span={12}>
          <Card
            title="待审批任务"
            extra={<Button type="link" onClick={() => navigate('/tasks/pending')}>查看更多</Button>}
          >
            {recentTasks.length === 0 ? (
              <div style={{ textAlign: 'center', color: '#999', padding: '20px' }}>
                暂无待审批任务
              </div>
            ) : (
              <List
                dataSource={recentTasks}
                renderItem={(item) => (
                  <List.Item>
                    <List.Item.Meta
                      title={item.node_name}
                      description={item.instance_title}
                    />
                    {getStatusTag(item.status)}
                  </List.Item>
                )}
              />
            )}
          </Card>
        </Col>
        <Col span={12}>
          <Card
            title="最新通知"
            extra={<Button type="link" onClick={() => navigate('/notifications')}>查看更多</Button>}
          >
            {notifications.length === 0 ? (
              <div style={{ textAlign: 'center', color: '#999', padding: '20px' }}>
                暂无通知
              </div>
            ) : (
              <List
                dataSource={notifications}
                renderItem={(item) => (
                  <List.Item>
                    <List.Item.Meta
                      title={item.title}
                      description={item.created_at}
                    />
                    {!item.is_read && <Tag color="red">未读</Tag>}
                  </List.Item>
                )}
              />
            )}
          </Card>
        </Col>
      </Row>
    </div>
  );
};
