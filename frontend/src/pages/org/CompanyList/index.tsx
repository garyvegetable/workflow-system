import { useState, useEffect } from 'react';
import { Table, Button, Modal, Form, Input, message, Popconfirm, Space } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { apiClient } from '@/api/client';

interface Company {
  id: number
  code: string
  name: string
  short_name: string
  status: number
}

export const CompanyList = () => {
  const [companies, setCompanies] = useState<Company[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingCompany, setEditingCompany] = useState<Company | null>(null);
  const [form] = Form.useForm();

  useEffect(() => {
    fetchCompanies();
  }, []);

  const fetchCompanies = async () => {
    setLoading(true);
    try {
      const response = await apiClient.get('/companies');
      setCompanies(response.data);
    } catch (error) {
      message.error('获取公司列表失败');
    } finally {
      setLoading(false);
    }
  };

  const handleAdd = () => {
    setEditingCompany(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = (record: Company) => {
    setEditingCompany(record);
    form.setFieldsValue(record);
    setModalVisible(true);
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      if (editingCompany) {
        await apiClient.put(`/companies/${editingCompany.id}`, values);
        message.success('更新成功');
      } else {
        await apiClient.post('/companies', values);
        message.success('创建成功');
      }
      setModalVisible(false);
      fetchCompanies();
    } catch (error) {
      message.error('操作失败');
    }
  };

  const handleDelete = async (id: number) => {
    try {
      await apiClient.delete(`/companies/${id}`);
      message.success('删除成功');
      fetchCompanies();
    } catch (error) {
      message.error('删除失败');
    }
  };

  const handleToggleStatus = async (company: Company) => {
    try {
      const newStatus = company.status === 1 ? 0 : 1;
      await apiClient.put(`/companies/${company.id}`, { status: newStatus });
      message.success(newStatus === 1 ? '已激活' : '已禁用');
      fetchCompanies();
    } catch (error) {
      message.error('操作失败');
    }
  };

  const columns = [
    { title: '公司代码', dataIndex: 'code' },
    { title: '公司名称', dataIndex: 'name' },
    { title: '简称', dataIndex: 'short_name' },
    {
      title: '状态',
      dataIndex: 'status',
      render: (status: number, record: Company) => (
        <Button
          type="link"
          onClick={() => handleToggleStatus(record)}
          style={{ color: status === 1 ? '#52c41a' : '#ff4d4f' }}
        >
          {status === 1 ? '正常' : '禁用'}
        </Button>
      ),
    },
    {
      title: '操作',
      render: (_: any, record: Company) => (
        <Space>
          <Button type="link" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
            编辑
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
      <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd} style={{ marginBottom: 16 }}>
        新增公司
      </Button>
      <Table columns={columns} dataSource={companies} rowKey="id" loading={loading} />
      <Modal
        title={editingCompany ? '编辑公司' : '新增公司'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
        width={500}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="code" label="公司代码" rules={[{ required: true, message: '请输入公司代码' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="name" label="公司名称" rules={[{ required: true, message: '请输入公司名称' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="short_name" label="简称">
            <Input />
          </Form.Item>
          {editingCompany && (
            <Form.Item name="status" label="状态">
              <Button
                type="primary"
                onClick={() => {
                  const currentStatus = form.getFieldValue('status');
                  const newStatus = currentStatus === 1 ? 0 : 1;
                  form.setFieldsValue({ status: newStatus });
                }}
                style={{ marginRight: 8 }}
              >
                {form.getFieldValue('status') === 1 ? '点击禁用' : '点击激活'}
              </Button>
              <span style={{ color: form.getFieldValue('status') === 1 ? '#52c41a' : '#ff4d4f' }}>
                当前状态: {form.getFieldValue('status') === 1 ? '正常' : '禁用'}
              </span>
            </Form.Item>
          )}
        </Form>
      </Modal>
    </div>
  );
};
