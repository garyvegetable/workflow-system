import { useState, useEffect } from 'react';
import { Table, Card, Space, Drawer, Descriptions, Button, message, Modal, Input } from 'antd';
import { approvalApi } from '@/api/approval';
import { TaskAction } from '@/components/approval/TaskAction';

interface Task {
  id: number
  instance_id: number
  node_name: string
  status: number
  create_time: string
  instance_title: string
  initiator_name: string
}

export const MyTasks = () => {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(false);
  const [detailOpen, setDetailOpen] = useState(false);
  const [selectedTask, setSelectedTask] = useState<Task | null>(null);
  const [selectedRowKeys, setSelectedRowKeys] = useState<number[]>([]);
  const [batchModalOpen, setBatchModalOpen] = useState(false);
  const [batchAction, setBatchAction] = useState<'approve' | 'reject'>('approve');
  const [batchComment, setBatchComment] = useState('');

  useEffect(() => {
    fetchTasks();
  }, []);

  const fetchTasks = async () => {
    setLoading(true);
    try {
      const response = await approvalApi.listPending();
      setTasks(response.data);
    } catch (error) {
      message.error('获取任务列表失败');
    } finally {
      setLoading(false);
    }
  };

  const handleViewDetail = (task: Task) => {
    setSelectedTask(task);
    setDetailOpen(true);
  };

  const handleBatchApprove = () => {
    if (selectedRowKeys.length === 0) {
      message.warning('请先选择要审批的任务');
      return;
    }
    setBatchAction('approve');
    setBatchComment('');
    setBatchModalOpen(true);
  };

  const handleBatchReject = () => {
    if (selectedRowKeys.length === 0) {
      message.warning('请先选择要驳回的任务');
      return;
    }
    setBatchAction('reject');
    setBatchComment('');
    setBatchModalOpen(true);
  };

  const handleBatchSubmit = async () => {
    try {
      if (batchAction === 'approve') {
        await approvalApi.batchApprove(selectedRowKeys, batchComment);
        message.success('批量审批成功');
      } else {
        await approvalApi.batchReject(selectedRowKeys, batchComment);
        message.success('批量驳回成功');
      }
      setBatchModalOpen(false);
      setSelectedRowKeys([]);
      fetchTasks();
    } catch (error) {
      message.error('批量操作失败');
    }
  };

  const rowSelection = {
    selectedRowKeys,
    onChange: (keys: React.Key[]) => setSelectedRowKeys(keys as number[]),
  };

  const columns = [
    { title: '任务名称', dataIndex: 'node_name' },
    { title: '所属申请', dataIndex: 'instance_title' },
    { title: '申请人', dataIndex: 'initiator_name' },
    { title: '创建时间', dataIndex: 'create_time' },
    {
      title: '操作',
      render: (_: any, record: Task) => (
        <Space>
          <Button type="link" onClick={() => handleViewDetail(record)}>
            查看详情
          </Button>
          <TaskAction taskId={record.id} onActionComplete={fetchTasks} />
        </Space>
      ),
    },
  ];

  return (
    <Card
      title="待处理任务"
      extra={
        <Space>
          <span>已选择 {selectedRowKeys.length} 项</span>
          <Button type="primary" onClick={handleBatchApprove} disabled={selectedRowKeys.length === 0}>
            批量通过
          </Button>
          <Button danger onClick={handleBatchReject} disabled={selectedRowKeys.length === 0}>
            批量驳回
          </Button>
        </Space>
      }
    >
      <Table
        columns={columns}
        dataSource={tasks}
        rowKey="id"
        loading={loading}
        rowSelection={rowSelection}
      />

      <Drawer
        title="任务详情"
        open={detailOpen}
        onClose={() => setDetailOpen(false)}
        width={600}
        extra={
          selectedTask && (
            <TaskAction taskId={selectedTask.id} onActionComplete={() => {
              setDetailOpen(false);
              fetchTasks();
            }} />
          )
        }
      >
        {selectedTask && (
          <Descriptions column={1} bordered>
            <Descriptions.Item label="任务名称">{selectedTask.node_name}</Descriptions.Item>
            <Descriptions.Item label="所属申请">{selectedTask.instance_title}</Descriptions.Item>
            <Descriptions.Item label="申请人">{selectedTask.initiator_name}</Descriptions.Item>
            <Descriptions.Item label="创建时间">{selectedTask.create_time}</Descriptions.Item>
          </Descriptions>
        )}
      </Drawer>

      <Modal
        title={batchAction === 'approve' ? '批量通过' : '批量驳回'}
        open={batchModalOpen}
        onOk={handleBatchSubmit}
        onCancel={() => setBatchModalOpen(false)}
        okText="确定"
        cancelText="取消"
      >
        <p>确认对 {selectedRowKeys.length} 个任务执行{batchAction === 'approve' ? '通过' : '驳回'}操作？</p>
        <Input.TextArea
          placeholder="审批意见（可选）"
          value={batchComment}
          onChange={(e) => setBatchComment(e.target.value)}
          rows={3}
        />
      </Modal>
    </Card>
  );
};
