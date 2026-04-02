import { useState, useEffect } from 'react';
import { Table, Card, Button, Tag, Space, message } from 'antd';
import { BellOutlined, CheckCircleOutlined } from '@ant-design/icons';
import { notificationApi } from '@/api/notification';

interface Notification {
  id: number
  title: string
  content: string
  type: string
  is_read: boolean
  created_at: string
}

export const NotificationList = () => {
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    fetchNotifications();
  }, []);

  const fetchNotifications = async () => {
    setLoading(true);
    try {
      const response = await notificationApi.list();
      setNotifications(response.data);
    } catch (error) {
      message.error('获取通知列表失败');
    } finally {
      setLoading(false);
    }
  };

  const handleMarkRead = async (id: number) => {
    try {
      await notificationApi.markAsRead(id);
      message.success('已标记为已读');
      fetchNotifications();
    } catch (error) {
      message.error('操作失败');
    }
  };

  const handleMarkAllRead = async () => {
    try {
      await notificationApi.markAllRead();
      message.success('已全部标记为已读');
      fetchNotifications();
    } catch (error) {
      message.error('操作失败');
    }
  };

  const columns = [
    {
      title: '状态',
      dataIndex: 'is_read',
      render: (isRead: boolean) =>
        isRead ? <Tag color="default">已读</Tag> : <Tag color="processing">未读</Tag>,
    },
    { title: '标题', dataIndex: 'title' },
    { title: '内容', dataIndex: 'content' },
    { title: '类型', dataIndex: 'type' },
    { title: '时间', dataIndex: 'created_at' },
    {
      title: '操作',
      render: (_: any, record: Notification) =>
        !record.is_read && (
          <Button type="link" size="small" onClick={() => handleMarkRead(record.id)}>
            标记已读
          </Button>
        ),
    },
  ];

  return (
    <Card
      title={
        <Space>
          <BellOutlined />
          <span>通知中心</span>
        </Space>
      }
      extra={
        <Button icon={<CheckCircleOutlined />} onClick={handleMarkAllRead}>
          全部标记已读
        </Button>
      }
    >
      <Table
        columns={columns}
        dataSource={notifications}
        rowKey="id"
        loading={loading}
        pagination={{ pageSize: 10 }}
      />
    </Card>
  );
};
