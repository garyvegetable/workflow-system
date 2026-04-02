import { useState, useEffect } from 'react';
import { Table, Button, Modal, Form, Input, message, Popconfirm } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { positionApi, Position } from '@/api/position';

export const PositionList = () => {
  const [positions, setPositions] = useState<Position[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingPosition, setEditingPosition] = useState<Position | null>(null);
  const [form] = Form.useForm();

  useEffect(() => {
    fetchPositions();
  }, []);

  const fetchPositions = async () => {
    setLoading(true);
    try {
      const response = await positionApi.list();
      setPositions(response.data || []);
    } catch (error) {
      message.error('获取职位列表失败');
    } finally {
      setLoading(false);
    }
  };

  const handleAdd = () => {
    setEditingPosition(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = (record: Position) => {
    setEditingPosition(record);
    form.setFieldsValue(record);
    setModalVisible(true);
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      if (editingPosition) {
        await positionApi.update(editingPosition.id, values);
        message.success('更新成功');
      } else {
        await positionApi.create(values);
        message.success('创建成功');
      }
      setModalVisible(false);
      fetchPositions();
    } catch (error) {
      message.error('操作失败');
    }
  };

  const handleDelete = async (id: number) => {
    try {
      await positionApi.delete(id);
      message.success('删除成功');
      fetchPositions();
    } catch (error) {
      message.error('删除失败');
    }
  };

  const columns = [
    { title: '职位名称', dataIndex: 'name' },
    { title: '职位代码', dataIndex: 'code' },
    { title: '状态', dataIndex: 'status', render: (status: number) => status === 1 ? '正常' : '禁用' },
    {
      title: '操作',
      render: (_: any, record: Position) => (
        <>
          <Button type="link" onClick={() => handleEdit(record)}>编辑</Button>
          <Popconfirm title="确定删除此职位?" onConfirm={() => handleDelete(record.id)}>
            <Button type="link" danger>删除</Button>
          </Popconfirm>
        </>
      ),
    },
  ];

  return (
    <div>
      <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd} style={{ marginBottom: 16 }}>
        新增职位
      </Button>
      <Table columns={columns} dataSource={positions} rowKey="id" loading={loading} />

      <Modal
        title={editingPosition ? '编辑职位' : '新增职位'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="职位名称" rules={[{ required: true, message: '请输入职位名称' }]}>
            <Input placeholder="如：经理、专员、总监" />
          </Form.Item>
          <Form.Item name="code" label="职位代码" rules={[{ required: true, message: '请输入职位代码' }]}>
            <Input placeholder="如：MGR、STAFF、DIR" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};
