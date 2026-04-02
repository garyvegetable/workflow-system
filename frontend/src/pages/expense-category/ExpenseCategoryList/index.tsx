import { useState, useEffect, useCallback } from 'react';
import { Button, Modal, Form, Input, Tree, message, Popconfirm, Select, Card, Alert } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import type { DataNode } from 'antd/es/tree';
import { expenseCategoryApi, ExpenseCategory } from '@/api/expense-category';

export const ExpenseCategoryList = () => {
  const [categories, setCategories] = useState<ExpenseCategory[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingCategory, setEditingCategory] = useState<ExpenseCategory | null>(null);
  const [form] = Form.useForm();

  const fetchCategories = useCallback(async () => {
    setLoading(true);
    try {
      const response = await expenseCategoryApi.list();
      setCategories(response.data);
    } catch (error) {
      message.error('获取费用科目列表失败');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchCategories();
  }, [fetchCategories]);

  // 将数据转换为树形结构
  const buildTreeData = (data: ExpenseCategory[]): DataNode[] => {
    const map: Record<number, DataNode> = {};
    const roots: DataNode[] = [];

    data.forEach(item => {
      map[item.id] = {
        key: String(item.id),
        title: `${item.code} - ${item.name}`,
        children: [],
      };
    });

    data.forEach(item => {
      if (item.parent_id && map[item.parent_id]) {
        map[item.parent_id].children?.push(map[item.id]);
      } else {
        roots.push(map[item.id]);
      }
    });

    return roots;
  };

  const handleAdd = () => {
    setEditingCategory(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = (category: ExpenseCategory) => {
    setEditingCategory(category);
    form.setFieldsValue({
      code: category.code,
      name: category.name,
      parent_id: category.parent_id,
    });
    setModalVisible(true);
  };

  const handleDelete = async (id: number) => {
    try {
      await expenseCategoryApi.delete(id);
      message.success('删除成功');
      fetchCategories();
    } catch (error) {
      message.error('删除失败');
    }
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      if (editingCategory) {
        await expenseCategoryApi.update(editingCategory.id, values);
        message.success('更新成功');
      } else {
        await expenseCategoryApi.create(values);
        message.success('创建成功');
      }
      setModalVisible(false);
      fetchCategories();
    } catch (error) {
      message.error('操作失败');
    }
  };

  const treeData = buildTreeData(categories);

  // 使用 Tree 组件渲染树形结构
  const TreeTitle = ({ title, record }: { title: string; record: ExpenseCategory }) => (
    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', width: '100%' }}>
      <span>{title}</span>
      <div onClick={(e) => e.stopPropagation()}>
        <Button type="link" size="small" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
          编辑
        </Button>
        <Popconfirm title="确定删除此费用科目?" onConfirm={() => handleDelete(record.id)}>
          <Button type="link" size="small" danger icon={<DeleteOutlined />}>
            删除
          </Button>
        </Popconfirm>
      </div>
    </div>
  );

  return (
    <div>
      <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd} style={{ marginBottom: 16 }}>
        新增费用科目
      </Button>

      {/* 树形展示 */}
      <Card>
        {loading ? null : categories.length === 0 ? (
          <Alert message="暂无费用科目数据" type="info" showIcon />
        ) : (
          <Tree
            showLine
            defaultExpandAll
            treeData={treeData}
            titleRender={(nodeData) => {
              const category = categories.find(c => c.id === Number(nodeData.key));
              if (!category) {return nodeData.title as string;}
              return (
                <TreeTitle
                  title={`${category.code} - ${category.name}`}
                  record={category}
                />
              );
            }}
          />
        )}
      </Card>

      {/* 新增/编辑 Modal */}
      <Modal
        title={editingCategory ? '编辑费用科目' : '新增费用科目'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
        width={500}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="code"
            label="科目编码"
            rules={[{ required: true, message: '请输入科目编码' }]}
          >
            <Input placeholder="请输入科目编码" />
          </Form.Item>
          <Form.Item
            name="name"
            label="科目名称"
            rules={[{ required: true, message: '请输入科目名称' }]}
          >
            <Input placeholder="请输入科目名称" />
          </Form.Item>
          <Form.Item name="parent_id" label="上级科目">
            <Select
              allowClear
              placeholder="选择上级科目（留空为顶级科目）"
              options={categories.map(c => ({ label: `${c.code} - ${c.name}`, value: c.id }))}
            />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};