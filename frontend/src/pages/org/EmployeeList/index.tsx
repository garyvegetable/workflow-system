import { useState, useEffect } from 'react';
import {
  Table,
  Button,
  Modal,
  Form,
  Input,
  Select,
  message,
  Popconfirm,
  Space,
  Tag,
  Drawer,
  Card,
  Descriptions,
  Divider,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  BankOutlined,
  CheckCircleFilled,
} from '@ant-design/icons';
import { employeeApi, Employee, BankAccount } from '@/api/employee';

const levelOptions = [
  { label: '员工', value: '员工' },
  { label: '主管', value: '主管' },
  { label: '经理', value: '经理' },
  { label: '总监', value: '总监' },
];

const statusOptions = [
  { label: '正常', value: 1 },
  { label: '禁用', value: 0 },
];

export const EmployeeList = () => {
  const [employees, setEmployees] = useState<Employee[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingEmployee, setEditingEmployee] = useState<Employee | null>(null);
  const [form] = Form.useForm();

  // Bank account drawer state
  const [bankDrawerVisible, setBankDrawerVisible] = useState(false);
  const [selectedEmployee, setSelectedEmployee] = useState<Employee | null>(null);
  const [bankAccounts, setBankAccounts] = useState<BankAccount[]>([]);
  const [bankLoading, setBankLoading] = useState(false);
  const [bankModalVisible, setBankModalVisible] = useState(false);
  const [editingBankAccount, setEditingBankAccount] = useState<BankAccount | null>(null);
  const [bankForm] = Form.useForm();

  const fetchEmployees = async () => {
    setLoading(true);
    try {
      const response = await employeeApi.list();
      setEmployees(response.data);
    } catch (error) {
      message.error('获取员工列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchEmployees();
  }, []);

  const handleAdd = () => {
    setEditingEmployee(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = (record: Employee) => {
    setEditingEmployee(record);
    form.setFieldsValue({
      username: record.username,
      email: record.email,
      level: record.level,
      status: record.status,
      company_id: record.company_id,
    });
    setModalVisible(true);
  };

  const handleDelete = async (id: number) => {
    try {
      await employeeApi.delete(id);
      message.success('删除成功');
      fetchEmployees();
    } catch (error) {
      message.error('删除失败');
    }
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      if (editingEmployee) {
        await employeeApi.update(editingEmployee.id, values);
        message.success('更新成功');
      } else {
        await employeeApi.create(values);
        message.success('创建成功');
      }
      setModalVisible(false);
      fetchEmployees();
    } catch (error) {
      message.error('操作失败');
    }
  };

  // Bank account management
  const openBankDrawer = async (employee: Employee) => {
    setSelectedEmployee(employee);
    setBankDrawerVisible(true);
    await fetchBankAccounts(employee.id);
  };

  const fetchBankAccounts = async (empId: number) => {
    setBankLoading(true);
    try {
      const response = await employeeApi.listBankAccounts(empId);
      setBankAccounts(response.data);
    } catch (error) {
      message.error('获取银行账户列表失败');
    } finally {
      setBankLoading(false);
    }
  };

  const handleAddBankAccount = () => {
    setEditingBankAccount(null);
    bankForm.resetFields();
    setBankModalVisible(true);
  };

  const handleEditBankAccount = (record: BankAccount) => {
    setEditingBankAccount(record);
    bankForm.setFieldsValue({
      bank_name: record.bank_name,
      bank_branch: record.bank_branch,
      bank_account: record.bank_account,
      account_holder: record.account_holder,
      is_default: record.is_default,
    });
    setBankModalVisible(true);
  };

  const handleDeleteBankAccount = async (aid: number) => {
    if (!selectedEmployee) {return;}
    try {
      await employeeApi.deleteBankAccount(selectedEmployee.id, aid);
      message.success('删除成功');
      fetchBankAccounts(selectedEmployee.id);
    } catch (error) {
      message.error('删除失败');
    }
  };

  const handleBankSubmit = async () => {
    if (!selectedEmployee) {return;}
    try {
      const values = await bankForm.validateFields();
      if (editingBankAccount) {
        await employeeApi.updateBankAccount(selectedEmployee.id, editingBankAccount.id, values);
        message.success('更新成功');
      } else {
        await employeeApi.createBankAccount(selectedEmployee.id, values);
        message.success('创建成功');
      }
      setBankModalVisible(false);
      fetchBankAccounts(selectedEmployee.id);
    } catch (error) {
      message.error('操作失败');
    }
  };

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 60 },
    { title: '用户名', dataIndex: 'username' },
    { title: '邮箱', dataIndex: 'email' },
    { title: '岗位级别', dataIndex: 'level' },
    {
      title: '状态',
      dataIndex: 'status',
      render: (status: number) => (
        <Tag color={status === 1 ? 'green' : 'red'}>{status === 1 ? '正常' : '禁用'}</Tag>
      ),
    },
    {
      title: '操作',
      render: (_: any, record: Employee) => (
        <Space>
          <Button type="link" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
            编辑
          </Button>
          <Button type="link" icon={<BankOutlined />} onClick={() => openBankDrawer(record)}>
            银行账户
          </Button>
          <Popconfirm title="确认删除?" onConfirm={() => handleDelete(record.id)}>
            <Button type="link" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  const bankColumns = [
    { title: '银行名称', dataIndex: 'bank_name' },
    { title: '支行', dataIndex: 'bank_branch' },
    { title: '账号', dataIndex: 'bank_account' },
    { title: '持卡人', dataIndex: 'account_holder' },
    {
      title: '默认',
      dataIndex: 'is_default',
      render: (isDefault: boolean) =>
        isDefault ? <CheckCircleFilled style={{ color: '#52c41a' }} /> : null,
    },
    {
      title: '操作',
      render: (_: any, record: BankAccount) => (
        <Space>
          <Button type="link" onClick={() => handleEditBankAccount(record)}>
            编辑
          </Button>
          <Popconfirm title="确认删除?" onConfirm={() => handleDeleteBankAccount(record.id)}>
            <Button type="link" danger>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <Button
        type="primary"
        icon={<PlusOutlined />}
        onClick={handleAdd}
        style={{ marginBottom: 16 }}
      >
        新增员工
      </Button>
      <Table columns={columns} dataSource={employees} rowKey="id" loading={loading} />

      {/* Employee Modal */}
      <Modal
        title={editingEmployee ? '编辑员工' : '新增员工'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
        width={500}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="username"
            label="用户名"
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input />
          </Form.Item>
          {!editingEmployee && (
            <Form.Item name="password" label="密码" rules={[{ required: true, message: '请输入密码' }]}>
              <Input.Password />
            </Form.Item>
          )}
          {editingEmployee && (
            <Form.Item name="password" label="新密码（留空则不修改）">
              <Input.Password />
            </Form.Item>
          )}
          <Form.Item name="email" label="邮箱">
            <Input />
          </Form.Item>
          <Form.Item name="level" label="岗位级别">
            <Select options={levelOptions} />
          </Form.Item>
          <Form.Item
            name="company_id"
            label="所属公司"
            rules={[{ required: true, message: '请输入公司ID' }]}
          >
            <Input type="number" />
          </Form.Item>
          <Form.Item name="status" label="状态" initialValue={1}>
            <Select options={statusOptions} />
          </Form.Item>
        </Form>
      </Modal>

      {/* Bank Account Drawer */}
      <Drawer
        title={
          <span>
            银行账户管理 - {selectedEmployee?.username}
          </span>
        }
        open={bankDrawerVisible}
        onClose={() => setBankDrawerVisible(false)}
        width={700}
        extra={
          <Button type="primary" icon={<PlusOutlined />} onClick={handleAddBankAccount}>
            新增账户
          </Button>
        }
      >
        {selectedEmployee && (
          <Card size="small" style={{ marginBottom: 16 }}>
            <Descriptions column={2}>
              <Descriptions.Item label="用户名">{selectedEmployee.username}</Descriptions.Item>
              <Descriptions.Item label="邮箱">{selectedEmployee.email}</Descriptions.Item>
              <Descriptions.Item label="岗位级别">{selectedEmployee.level}</Descriptions.Item>
              <Descriptions.Item label="状态">
                {selectedEmployee.status === 1 ? '正常' : '禁用'}
              </Descriptions.Item>
            </Descriptions>
          </Card>
        )}
        <Divider>银行账户列表</Divider>
        <Table
          columns={bankColumns}
          dataSource={bankAccounts}
          rowKey="id"
          loading={bankLoading}
          pagination={false}
        />

        {/* Bank Account Modal */}
        <Modal
          title={editingBankAccount ? '编辑银行账户' : '新增银行账户'}
          open={bankModalVisible}
          onOk={handleBankSubmit}
          onCancel={() => setBankModalVisible(false)}
        >
          <Form form={bankForm} layout="vertical">
            <Form.Item
              name="bank_name"
              label="银行名称"
              rules={[{ required: true, message: '请输入银行名称' }]}
            >
              <Input />
            </Form.Item>
            <Form.Item name="bank_branch" label="支行">
              <Input />
            </Form.Item>
            <Form.Item
              name="bank_account"
              label="账号"
              rules={[{ required: true, message: '请输入账号' }]}
            >
              <Input />
            </Form.Item>
            <Form.Item
              name="account_holder"
              label="持卡人"
              rules={[{ required: true, message: '请输入持卡人姓名' }]}
            >
              <Input />
            </Form.Item>
            <Form.Item name="is_default" label="设为默认" valuePropName="checked">
              <Input type="checkbox" />
            </Form.Item>
          </Form>
        </Modal>
      </Drawer>
    </div>
  );
};
