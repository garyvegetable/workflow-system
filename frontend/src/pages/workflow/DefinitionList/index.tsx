import { useState, useEffect } from 'react';
import { Table, Button, Space, Tag, message, Modal, Select, Popconfirm } from 'antd';
import { PlusOutlined, EditOutlined, UploadOutlined, StopOutlined, CopyOutlined, DeleteOutlined, CheckCircleOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { workflowApi } from '@/api/workflow';
import { companyApi } from '@/api/company';

interface Workflow {
  id: number
  code: string
  name: string
  version: number
  status: number
  company_id: number
}

interface Company {
  id: number
  name: string
}

export const WorkflowList = () => {
  const [workflows, setWorkflows] = useState<Workflow[]>([]);
  const [loading, setLoading] = useState(false);
  const [copyModalVisible, setCopyModalVisible] = useState(false);
  const [selectedWorkflow, setSelectedWorkflow] = useState<Workflow | null>(null);
  const [targetCompanyId, setTargetCompanyId] = useState<number | null>(null);
  const [companies, setCompanies] = useState<Company[]>([]);
  const [copyLoading, setCopyLoading] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    fetchWorkflows();
    fetchCompanies();
  }, []);

  const fetchWorkflows = async () => {
    setLoading(true);
    try {
      const response = await workflowApi.list();
      setWorkflows(response.data);
    } catch (error) {
      message.error('获取流程列表失败');
    } finally {
      setLoading(false);
    }
  };

  const fetchCompanies = async () => {
    try {
      const response = await companyApi.list();
      setCompanies(response.data);
    } catch (error) {
      console.error('获取公司列表失败', error);
    }
  };

  const getStatusTag = (status: number) => {
    switch (status) {
      case 1: return <Tag color="default">草稿</Tag>;
      case 2: return <Tag color="success">已发布</Tag>;
      case 3: return <Tag color="error">已禁用</Tag>;
      default: return <Tag>未知</Tag>;
    }
  };

  const handlePublish = async (id: number) => {
    try {
      await workflowApi.publish(id);
      message.success('发布成功');
      fetchWorkflows();
    } catch (error) {
      message.error('发布失败');
    }
  };

  const handleDisable = async (id: number) => {
    try {
      await workflowApi.disable(id);
      message.success('已禁用');
      fetchWorkflows();
    } catch (error) {
      message.error('禁用失败');
    }
  };

  const handleEnable = async (id: number) => {
    try {
      await workflowApi.enable(id);
      message.success('已启用');
      fetchWorkflows();
    } catch (error) {
      message.error('启用失败');
    }
  };

  const handleDelete = async (id: number) => {
    try {
      await workflowApi.delete(id);
      message.success('删除成功');
      fetchWorkflows();
    } catch (error) {
      message.error('删除失败');
    }
  };

  const handleCopy = async () => {
    if (!selectedWorkflow || !targetCompanyId) {
      message.error('请选择目标公司');
      return;
    }
    setCopyLoading(true);
    try {
      await workflowApi.copy(selectedWorkflow.id, targetCompanyId);
      message.success('复制成功');
      setCopyModalVisible(false);
      setSelectedWorkflow(null);
      setTargetCompanyId(null);
    } catch (error) {
      message.error('复制失败');
    } finally {
      setCopyLoading(false);
    }
  };

  const openCopyModal = (workflow: Workflow) => {
    setSelectedWorkflow(workflow);
    setCopyModalVisible(true);
  };

  const columns = [
    { title: '流程代码', dataIndex: 'code' },
    { title: '流程名称', dataIndex: 'name' },
    { title: '版本', dataIndex: 'version' },
    { title: '状态', dataIndex: 'status', render: (status: number) => getStatusTag(status) },
    {
      title: '操作',
      render: (_: any, record: Workflow) => (
        <Space>
          {record.status === 2 && (
            <Button type="link" onClick={() => navigate(`/workflows/apply/${record.id}`)}>
              申请
            </Button>
          )}
          <Button type="link" icon={<EditOutlined />} onClick={() => navigate(`/workflows/designer/${record.id}`)}>
            编辑
          </Button>
          {record.status === 1 && (
            <Button type="link" icon={<UploadOutlined />} onClick={() => handlePublish(record.id)}>
              发布
            </Button>
          )}
          {record.status === 2 && (
            <Button type="link" danger icon={<StopOutlined />} onClick={() => handleDisable(record.id)}>
              禁用
            </Button>
          )}
          {record.status === 3 && (
            <Button type="link" icon={<CheckCircleOutlined />} onClick={() => handleEnable(record.id)}>
              启用
            </Button>
          )}
          <Button type="link" icon={<CopyOutlined />} onClick={() => openCopyModal(record)}>
            复制
          </Button>
          <Popconfirm
            title="确认删除"
            description="删除后不可恢复，确定要删除吗？"
            onConfirm={() => handleDelete(record.id)}
            okText="删除"
            cancelText="取消"
            okButtonProps={{ danger: true }}
          >
            <Button type="link" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <div style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => navigate('/workflows/designer')}>
          新建流程
        </Button>
      </div>
      <Table columns={columns} dataSource={workflows} rowKey="id" loading={loading} />

      <Modal
        title="复制流程到其他公司"
        open={copyModalVisible}
        onOk={handleCopy}
        onCancel={() => {
          setCopyModalVisible(false);
          setSelectedWorkflow(null);
          setTargetCompanyId(null);
        }}
        confirmLoading={copyLoading}
      >
        <div style={{ marginBottom: 16 }}>
          <p>源流程: {selectedWorkflow?.name}</p>
        </div>
        <div>
          <label>目标公司: </label>
          <Select
            style={{ width: 200 }}
            placeholder="请选择公司"
            value={targetCompanyId}
            onChange={setTargetCompanyId}
            options={companies.map(c => ({ value: c.id, label: c.name }))}
          />
        </div>
      </Modal>
    </div>
  );
};
