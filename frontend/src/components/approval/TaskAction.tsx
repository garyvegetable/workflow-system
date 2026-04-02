import { Button, Space, Modal, Form, Input, message } from 'antd';
import { CheckOutlined, CloseOutlined, SwapOutlined } from '@ant-design/icons';
import { useState } from 'react';
import { approvalApi } from '@/api/approval';

interface TaskActionProps {
  taskId: number
  onActionComplete: () => void
}

export const TaskAction = ({ taskId, onActionComplete }: TaskActionProps) => {
  const [approveModalOpen, setApproveModalOpen] = useState(false);
  const [rejectModalOpen, setRejectModalOpen] = useState(false);
  const [transferModalOpen, setTransferModalOpen] = useState(false);
  const [loading, setLoading] = useState(false);
  const [form] = Form.useForm();

  const handleApprove = async () => {
    try {
      setLoading(true);
      await approvalApi.approve(taskId, { comment: form.getFieldValue('comment') || '' });
      message.success('审批已通过');
      setApproveModalOpen(false);
      onActionComplete();
    } catch (error) {
      message.error('操作失败');
    } finally {
      setLoading(false);
    }
  };

  const handleReject = async () => {
    try {
      setLoading(true);
      await approvalApi.reject(taskId, { comment: form.getFieldValue('comment') || '' });
      message.success('已驳回');
      setRejectModalOpen(false);
      onActionComplete();
    } catch (error) {
      message.error('操作失败');
    } finally {
      setLoading(false);
    }
  };

  const handleTransfer = async () => {
    const values = await form.validateFields();
    try {
      setLoading(true);
      await approvalApi.transfer(taskId, { new_assignee_id: values.targetUserId, comment: values.comment });
      message.success('已转交');
      setTransferModalOpen(false);
      onActionComplete();
    } catch (error) {
      message.error('操作失败');
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <Space>
        <Button type="primary" icon={<CheckOutlined />} onClick={() => setApproveModalOpen(true)}>
          同意
        </Button>
        <Button danger icon={<CloseOutlined />} onClick={() => setRejectModalOpen(true)}>
          驳回
        </Button>
        <Button icon={<SwapOutlined />} onClick={() => setTransferModalOpen(true)}>
          转交
        </Button>
      </Space>

      <Modal title="同意" open={approveModalOpen} onOk={handleApprove} confirmLoading={loading}>
        <Form form={form} layout="vertical">
          <Form.Item label="审批意见" name="comment">
            <Input.TextArea rows={3} placeholder="请输入审批意见（可选）" />
          </Form.Item>
        </Form>
      </Modal>

      <Modal title="驳回" open={rejectModalOpen} onOk={handleReject} confirmLoading={loading}>
        <Form form={form} layout="vertical">
          <Form.Item label="驳回原因" name="comment" rules={[{ required: true, message: '请输入驳回原因' }]}>
            <Input.TextArea rows={3} placeholder="请输入驳回原因" />
          </Form.Item>
        </Form>
      </Modal>

      <Modal title="转交任务" open={transferModalOpen} onOk={handleTransfer} confirmLoading={loading}>
        <Form form={form} layout="vertical">
          <Form.Item label="转交给" name="targetUserId" rules={[{ required: true, message: '请输入目标用户ID' }]}>
            <Input placeholder="请输入用户ID" />
          </Form.Item>
          <Form.Item label="转交原因" name="comment">
            <Input.TextArea rows={2} placeholder="请输入转交原因（可选）" />
          </Form.Item>
        </Form>
      </Modal>
    </>
  );
};
