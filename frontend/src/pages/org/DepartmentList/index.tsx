import { useState, useEffect, useCallback } from 'react';
import { Button, Modal, Form, Input, Tree, message, Card, Alert, Select } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, SettingOutlined, DragOutlined } from '@ant-design/icons';
import type { DataNode } from 'antd/es/tree';
import { departmentApi, Department, ApprovalChainStep } from '@/api/department';
import { employeeApi, Employee } from '@/api/employee';
import { useSelector } from 'react-redux';
import type { RootState } from '@/store';

export const DepartmentList = () => {
  const companyId = useSelector((state: RootState) => state.auth.user?.companyId);
  const [departments, setDepartments] = useState<Department[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingDept, setEditingDept] = useState<Department | null>(null);
  const [chainModalVisible, setChainModalVisible] = useState(false);
  const [selectedDeptId, setSelectedDeptId] = useState<number | null>(null);
  const [approvalChain, setApprovalChain] = useState<ApprovalChainStep[]>([]);
  const [chainLoading, setChainLoading] = useState(false);
  const [deleteModalVisible, setDeleteModalVisible] = useState(false);
  const [deleteTargetDept, setDeleteTargetDept] = useState<Department | null>(null);
  const [transferDeptId, setTransferDeptId] = useState<number | null>(null);
  const [childCount, setChildCount] = useState(0);
  const [form] = Form.useForm();
  const [chainEmployeeOptions, setChainEmployeeOptions] = useState<Employee[]>([]);
  const [leaderOptions, setLeaderOptions] = useState<Employee[]>([]);
  const [searchingEmployee, setSearchingEmployee] = useState(false);

  const fetchDepartments = useCallback(async () => {
    setLoading(true);
    try {
      const response = await departmentApi.list();
      setDepartments(response.data || []);
    } catch (error) {
      message.error('获取部门列表失败');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchDepartments();
  }, [fetchDepartments]);

  // 将数据转换为树形结构
  const buildTreeData = (data: Department[]): DataNode[] => {
    const map: Record<number, DataNode> = {};
    const roots: DataNode[] = [];

    data.forEach(item => {
      map[item.id] = {
        key: String(item.id),
        title: item.name,
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

  // 统计子部门数量
  const countChildren = (deptId: number, allDepts: Department[]): number => {
    let count = 0;
    allDepts.forEach(d => {
      if (d.parent_id === deptId) {
        count++;
        count += countChildren(d.id, allDepts);
      }
    });
    return count;
  };

  const handleAdd = () => {
    setEditingDept(null);
    form.resetFields();
    // 加载员工列表用于选择负责人
    employeeApi.list(companyId).then(res => {
      setLeaderOptions(res.data || []);
    }).catch(() => {});
    setModalVisible(true);
  };

  const handleEdit = (dept: Department) => {
    setEditingDept(dept);
    form.setFieldsValue({
      name: dept.name,
      parent_id: dept.parent_id,
      leader_id: dept.leader_id,
    });
    // 加载员工列表用于选择负责人
    employeeApi.list(companyId).then(res => {
      setLeaderOptions(res.data || []);
    }).catch(() => {});
    setModalVisible(true);
  };

  const handleDeleteClick = (dept: Department) => {
    const children = countChildren(dept.id, departments);
    setDeleteTargetDept(dept);
    setChildCount(children);
    setTransferDeptId(null);

    if (children > 0) {
      // 有子部门，需要选择接收部门
      setDeleteModalVisible(true);
    } else {
      // 没有子部门，直接删除
      Modal.confirm({
        title: '确认删除',
        content: `确定删除部门"${dept.name}"吗？`,
        onOk: () => handleDelete(dept.id, 0),
      });
    }
  };

  const handleDelete = async (deptId: number, transferId: number) => {
    try {
      await departmentApi.delete(deptId, transferId || undefined);
      message.success('删除成功');
      setDeleteModalVisible(false);
      fetchDepartments();
    } catch (error: any) {
      const errMsg = error?.response?.data?.error || '删除失败';
      message.error(errMsg);
    }
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      // 处理 leader_id：如果为空或undefined，转换为null
      const submitData: any = { ...values };
      if (!submitData.leader_id) {
        submitData.leader_id = null;
      }
      if (editingDept) {
        await departmentApi.update(editingDept.id, submitData);
        message.success('更新成功');
      } else {
        await departmentApi.create({ ...submitData, company_id: companyId });
        message.success('创建成功');
      }
      setModalVisible(false);
      fetchDepartments();
    } catch (error: any) {
      const errMsg = error?.response?.data?.error || '操作失败';
      message.error(errMsg);
    }
  };

  const handleConfigChain = async (dept: Department) => {
    setSelectedDeptId(dept.id);
    setChainLoading(true);
    try {
      const response = await departmentApi.getApprovalChain(dept.id);
      setApprovalChain(response.data.map((step, index) => ({
        ...step,
        step_order: index + 1,
      })));
    } catch (error) {
      setApprovalChain([]);
    } finally {
      setChainLoading(false);
    }
    setChainModalVisible(true);
  };

  const handleAddChainStep = () => {
    setApprovalChain([
      ...approvalChain,
      { employee_id: 0, step_order: approvalChain.length + 1 },
    ]);
  };

  const handleRemoveChainStep = (index: number) => {
    const newChain = approvalChain.filter((_, i) => i !== index)
      .map((step, i) => ({ ...step, step_order: i + 1 }));
    setApprovalChain(newChain);
  };

  const handleChainStepChange = (index: number, field: keyof ApprovalChainStep, value: number) => {
    const newChain = [...approvalChain];
    newChain[index] = { ...newChain[index], [field]: value };
    setApprovalChain(newChain);
  };

  const handleMoveChainStep = (index: number, direction: 'up' | 'down') => {
    if (
      (direction === 'up' && index === 0) ||
      (direction === 'down' && index === approvalChain.length - 1)
    ) {
      return;
    }
    const newChain = [...approvalChain];
    const targetIndex = direction === 'up' ? index - 1 : index + 1
    ;[newChain[index], newChain[targetIndex]] = [newChain[targetIndex], newChain[index]];
    newChain.forEach((step, i) => { step.step_order = i + 1; });
    setApprovalChain(newChain);
  };

  const handleSaveChain = async () => {
    if (!selectedDeptId) {return;}
    try {
      await departmentApi.setApprovalChain(
        selectedDeptId,
        approvalChain.map(step => ({
          employee_id: step.employee_id,
          step_order: step.step_order,
        })),
      );
      message.success('审批链保存成功');
      setChainModalVisible(false);
    } catch (error) {
      message.error('保存失败');
    }
  };

  // 获取可用的转移部门列表（排除要删除的部门及其子部门）
  const getAvailableTransferDepts = (): Department[] => {
    if (!deleteTargetDept) {return [];}
    const excludeIds = new Set<number>();
    const collectIds = (deptId: number) => {
      excludeIds.add(deptId);
      departments.forEach(d => {
        if (d.parent_id === deptId) {
          collectIds(d.id);
        }
      });
    };
    collectIds(deleteTargetDept.id);
    return departments.filter(d => !excludeIds.has(d.id));
  };

  const treeData = buildTreeData(departments);

  // 使用 Tree 组件渲染树形结构
  const TreeTitle = ({ title, record }: { title: string; record: Department }) => (
    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', width: '100%' }}>
      <span>{title}</span>
      <div onClick={(e) => e.stopPropagation()}>
        <Button type="link" size="small" icon={<SettingOutlined />} onClick={() => handleConfigChain(record)}>
          审批链
        </Button>
        <Button type="link" size="small" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
          编辑
        </Button>
        <Button type="link" size="small" danger icon={<DeleteOutlined />} onClick={() => handleDeleteClick(record)}>
          删除
        </Button>
      </div>
    </div>
  );

  return (
    <div>
      <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd} style={{ marginBottom: 16 }}>
        新增部门
      </Button>

      {/* 树形展示 */}
      <Card>
        {loading ? null : departments.length === 0 ? (
          <Alert message="暂无部门数据" type="info" showIcon />
        ) : (
          <Tree
            showLine
            defaultExpandAll
            treeData={treeData}
            titleRender={(nodeData) => {
              const dept = departments.find(d => d.id === Number(nodeData.key));
              if (!dept) {return nodeData.title as string;}
              return (
                <TreeTitle
                  title={dept.name}
                  record={dept}
                />
              );
            }}
          />
        )}
      </Card>

      {/* 新增/编辑 Modal */}
      <Modal
        title={editingDept ? '编辑部门' : '新增部门'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
        width={500}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="name"
            label="部门名称"
            rules={[{ required: true, message: '请输入部门名称' }]}
          >
            <Input placeholder="请输入部门名称" />
          </Form.Item>
          <Form.Item name="parent_id" label="上级部门">
            <Select
              allowClear
              placeholder="选择上级部门（留空为顶级部门）"
              options={departments.map(d => ({ label: d.name, value: d.id }))}
            />
          </Form.Item>
          <Form.Item name="leader_id" label="部门负责人">
            <Select
              showSearch
              allowClear
              placeholder="选择部门负责人"
              optionFilterProp="label"
              filterOption={(input, option) =>
                (option?.label ?? '').toLowerCase().includes(input.toLowerCase())
              }
              options={leaderOptions.map(e => ({
                label: `${e.name || e.username} (${e.email || '无邮箱'})`,
                value: e.id,
              }))}
              onSearch={async (value) => {
                if (value.length < 1) {
                  setLeaderOptions([]);
                  return;
                }
                setSearchingEmployee(true);
                try {
                  const response = await employeeApi.search(value);
                  setLeaderOptions(response.data);
                } catch {
                  message.error('搜索员工失败');
                } finally {
                  setSearchingEmployee(false);
                }
              }}
              loading={searchingEmployee}
            />
          </Form.Item>
        </Form>
      </Modal>

      {/* 删除确认 Modal（用于选择转移部门） */}
      <Modal
        title="确认删除"
        open={deleteModalVisible}
        onOk={() => {
          if (childCount > 0 && !transferDeptId) {
            message.error('请选择接收部门');
            return;
          }
          handleDelete(deleteTargetDept!.id, transferDeptId || 0);
        }}
        onCancel={() => setDeleteModalVisible(false)}
      >
        {childCount > 0 ? (
          <>
            <p>该部门下有 <strong>{childCount}</strong> 个子部门，确认删除？</p>
            <p>涉及的员工将转移到：</p>
            <Select
              style={{ width: '100%', marginTop: 8 }}
              placeholder="选择接收部门"
              value={transferDeptId}
              onChange={setTransferDeptId}
              options={getAvailableTransferDepts().map(d => ({ label: d.name, value: d.id }))}
            />
          </>
        ) : (
          <p>确定删除部门&quot;{deleteTargetDept?.name}&quot;吗？</p>
        )}
      </Modal>

      {/* 审批链配置 Modal */}
      <Modal
        title="配置部门审批链"
        open={chainModalVisible}
        onOk={handleSaveChain}
        onCancel={() => setChainModalVisible(false)}
        width={600}
        destroyOnClose
      >
        {chainLoading ? null : (
          <>
            <Alert
              message="拖拽排序调整审批顺序"
              type="info"
              showIcon
              style={{ marginBottom: 16 }}
            />
            {approvalChain.length === 0 ? (
              <Alert message="暂未配置审批链" type="warning" showIcon />
            ) : (
              <div>
                {approvalChain.map((step, index) => (
                  <Card
                    key={index}
                    size="small"
                    style={{ marginBottom: 8 }}
                    extra={
                      <div style={{ display: 'flex', gap: 8 }}>
                        <Button
                          size="small"
                          icon={<DragOutlined />}
                          disabled={index === 0}
                          onClick={() => handleMoveChainStep(index, 'up')}
                        />
                        <Button
                          size="small"
                          icon={<DragOutlined />}
                          disabled={index === approvalChain.length - 1}
                          onClick={() => handleMoveChainStep(index, 'down')}
                        />
                        <Button
                          size="small"
                          danger
                          onClick={() => handleRemoveChainStep(index)}
                        >
                          移除
                        </Button>
                      </div>
                    }
                  >
                    <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                      <span style={{ fontWeight: 'bold', color: '#1890ff' }}>
                        步骤 {step.step_order}
                      </span>
                      <Select
                        showSearch
                        placeholder="选择审批人"
                        allowClear
                        filterOption={false}
                        style={{ width: 200 }}
                        value={step.employee_id || undefined}
                        onSearch={async (value) => {
                          if (value.length < 1) {
                            setChainEmployeeOptions([]);
                            return;
                          }
                          setSearchingEmployee(true);
                          try {
                            const response = await employeeApi.search(value);
                            setChainEmployeeOptions(response.data);
                          } catch {
                            message.error('搜索员工失败');
                          } finally {
                            setSearchingEmployee(false);
                          }
                        }}
                        onChange={(value) => handleChainStepChange(index, 'employee_id', value || 0)}
                        loading={searchingEmployee}
                        options={chainEmployeeOptions.map(e => ({
                          label: `${e.name || e.username} (${e.email})`,
                          value: e.id,
                        }))}
                      />
                    </div>
                  </Card>
                ))}
              </div>
            )}
            <Button
              type="dashed"
              onClick={handleAddChainStep}
              style={{ marginTop: 8, width: '100%' }}
              icon={<PlusOutlined />}
            >
              添加审批步骤
            </Button>
          </>
        )}
      </Modal>
    </div>
  );
};
