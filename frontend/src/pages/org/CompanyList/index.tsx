import { useState, useEffect } from 'react';
import { Table, Button, Modal, Form, Input, message } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
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

  const columns = [
    { title: '公司代码', dataIndex: 'code' },
    { title: '公司名称', dataIndex: 'name' },
    { title: '简称', dataIndex: 'short_name' },
    { title: '状态', dataIndex: 'status', render: (status: number) => status === 1 ? '正常' : '禁用' },
    {
      title: '操作',
      render: (_: any, record: Company) => (
        <Button type="link" onClick={() => handleEdit(record)}>编辑</Button>
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
      >
        <Form form={form} layout="vertical">
          <Form.Item name="code" label="公司代码" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="name" label="公司名称" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="short_name" label="简称">
            <Input />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};
