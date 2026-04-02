import { useState, useEffect } from 'react';
import { Table, Tag, Card, Button, message, Modal, Timeline, Space } from 'antd';
import { workflowApi } from '@/api/workflow';
import { approvalApi } from '@/api/approval';

interface Application {
  id: number
  title: string
  definition_id: number
  status: number
  started_at: string
}

interface HistoryItem {
  id: number
  node_name: string
  assignee_id: number
  status: number
  action: string
  comment: string
  completed_at: string
}

const statusMap = {
  0: { text: '草稿', color: 'default' },
  1: { text: '审批中', color: 'processing' },
  2: { text: '已通过', color: 'success' },
  3: { text: '已驳回', color: 'error' },
  4: { text: '已撤回', color: 'warning' },
};

export const MyApplications = () => {
  const [applications, setApplications] = useState<Application[]>([]);
  const [loading, setLoading] = useState(false);
  const [historyVisible, setHistoryVisible] = useState(false);
  const [history, setHistory] = useState<HistoryItem[]>([]);
  const [historyLoading, setHistoryLoading] = useState(false);

  useEffect(() => {
    fetchApplications();
  }, []);

  const fetchApplications = async () => {
    setLoading(true);
    try {
      const response = await workflowApi.getMyApplications();
      setApplications(response.data);
    } catch {
      message.error('获取申请列表失败');
    } finally {
      setLoading(false);
    }
  };

  const handleCancel = async (id: number) => {
    try {
      await workflowApi.cancelInstance(id);
      message.success('已撤回');
      fetchApplications();
    } catch {
      message.error('撤回失败');
    }
  };

  const handleViewHistory = async (id: number) => {
    setHistoryVisible(true);
    setHistoryLoading(true);
    try {
      const response = await approvalApi.getHistory(id);
      setHistory(response.data);
    } catch {
      message.error('获取审批历史失败');
    } finally {
      setHistoryLoading(false);
    }
  };

  const columns = [
    { title: '标题', dataIndex: 'title' },
    {
      title: '状态',
      dataIndex: 'status',
      render: (status: number) => {
        const s = statusMap[status as keyof typeof statusMap] || statusMap[0];
        return <Tag color={s.color}>{s.text}</Tag>;
      },
    },
    { title: '提交时间', dataIndex: 'started_at' },
    {
      title: '操作',
      render: (_: any, record: Application) => (
        <Space>
          <Button type="link" onClick={() => handleViewHistory(record.id)}>
            查看详情
          </Button>
          {record.status === 1 && (
            <Button type="link" danger onClick={() => handleCancel(record.id)}>
              撤回
            </Button>
          )}
        </Space>
      ),
    },
  ];

  return (
    <Card title="我的申请">
      <Table columns={columns} dataSource={applications} rowKey="id" loading={loading} />

      <Modal
        title="审批历史"
        open={historyVisible}
        onCancel={() => setHistoryVisible(false)}
        footer={null}
        width={600}
      >
        {historyLoading ? null : history.length === 0 ? (
          <p>暂无审批记录</p>
        ) : (
          <Timeline
            items={history.map(h => ({
              color: h.action === 'approve' ? 'green' : h.action === 'reject' ? 'red' : 'blue',
              children: (
                <div>
                  <strong>{h.node_name}</strong>
                  <br />
                  <span>审批人ID: {h.assignee_id}</span>
                  <br />
                  <span>结果: {h.action === 'approve' ? '通过' : h.action === 'reject' ? '驳回' : '待处理'}</span>
                  {h.comment && (
                    <>
                      <br />
                      <span>意见: {h.comment}</span>
                    </>
                  )}
                  <br />
                  <span style={{ fontSize: '12px', color: '#999' }}>{h.completed_at}</span>
                </div>
              ),
            }))}
          />
        )}
      </Modal>
    </Card>
  );
};
