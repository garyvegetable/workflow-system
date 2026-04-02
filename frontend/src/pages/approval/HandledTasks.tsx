import { useState, useEffect } from 'react';
import { Table, Tag, message } from 'antd';
import { apiClient } from '@/api/client';

interface Task {
  id: number
  instance_id: number
  node_name: string
  status: number
  action?: string
  comment?: string
}

export const HandledTasks = () => {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    fetchTasks();
  }, []);

  const fetchTasks = async () => {
    setLoading(true);
    try {
      const response = await apiClient.get('/tasks/handled');
      setTasks(response.data);
    } catch (error) {
      message.error('获取已处理任务失败');
    } finally {
      setLoading(false);
    }
  };

  const getActionTag = (action?: string) => {
    switch (action) {
      case 'approve': return <Tag color="success">同意</Tag>;
      case 'reject': return <Tag color="error">驳回</Tag>;
      default: return <Tag>未知</Tag>;
    }
  };

  const columns = [
    { title: '任务ID', dataIndex: 'id' },
    { title: '流程实例ID', dataIndex: 'instance_id' },
    { title: '节点名称', dataIndex: 'node_name' },
    { title: '审批动作', dataIndex: 'action', render: getActionTag },
    { title: '审批意见', dataIndex: 'comment' },
  ];

  return (
    <div>
      <Table columns={columns} dataSource={tasks} rowKey="id" loading={loading} />
    </div>
  );
};
