import { useState, useCallback, useEffect } from 'react';
import { ReactFlow, Controls, Background, addEdge, applyNodeChanges, applyEdgeChanges, Connection, NodeTypes, Panel } from '@xyflow/react';
import '@xyflow/react/dist/style.css';
import { Button, message, Space, Form, Input, Modal, Select } from 'antd';
import { apiClient } from '@/api/client';
import { employeeApi, Employee } from '@/api/employee';
import { departmentApi, Department } from '@/api/department';
import { StartNode, EndNode, ApprovalNode, ConditionNode, CountersignNode } from '@/components/workflow/nodes';

const nodeTypes: NodeTypes = {
  start: StartNode,
  end: EndNode,
  approval: ApprovalNode,
  condition: ConditionNode,
  countersign: CountersignNode,
};

export const WorkflowDesigner = () => {
  const [nodes, setNodes] = useState<any[]>([]);
  const [edges, setEdges] = useState<any[]>([]);
  const [workflowId, setWorkflowId] = useState<string | null>(null);
  const [selectedNode, setSelectedNode] = useState<any>(null);
  const [propertyModalOpen, setPropertyModalOpen] = useState(false);
  const [form] = Form.useForm();
  const [employeeOptions, setEmployeeOptions] = useState<Employee[]>([]);
  const [departmentOptions, setDepartmentOptions] = useState<Department[]>([]);
  const [searchingEmployee, setSearchingEmployee] = useState(false);

  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    const id = params.get('id');
    if (id) {
      setWorkflowId(id);
      fetchWorkflow(parseInt(id));
    } else {
      // 加载草稿
      const draftKey = `workflow_draft_${window.location.pathname}`;
      const savedDraft = localStorage.getItem(draftKey);
      if (savedDraft) {
        try {
          const draft = JSON.parse(savedDraft);
          if (draft.nodes) {setNodes(draft.nodes);}
          if (draft.edges) {setEdges(draft.edges);}
          message.info('已恢复上次编辑的草稿');
        } catch {
          // Ignore parse errors - draft may be corrupted or from old version
        }
      }
    }
  }, []);

  // 草稿自动保存
  useEffect(() => {
    if (nodes.length === 0 && edges.length === 0) {return;}
    const draftKey = `workflow_draft_${window.location.pathname}`;
    const draft = { nodes, edges, savedAt: new Date().toISOString() };
    localStorage.setItem(draftKey, JSON.stringify(draft));
  }, [nodes, edges]);

  // 清除草稿
  const clearDraft = () => {
    const draftKey = `workflow_draft_${window.location.pathname}`;
    localStorage.removeItem(draftKey);
    message.success('草稿已清除');
  };

  const fetchWorkflow = async (id: number) => {
    try {
      const response = await apiClient.get(`/workflows/${id}`);
      if (response.data.graph_data) {
        const graphData = typeof response.data.graph_data === 'string'
          ? JSON.parse(response.data.graph_data)
          : response.data.graph_data;
        if (graphData.nodes) {setNodes(graphData.nodes);}
        if (graphData.edges) {setEdges(graphData.edges);}
      }
    } catch (error) {
      message.error('获取流程数据失败');
    }
  };

  const onConnect = useCallback(
    (params: Connection) => setEdges((eds: any[]) => addEdge(params, eds)),
    [setEdges],
  );

  const handleSave = async () => {
    try {
      const graphData = { nodes, edges };
      if (workflowId) {
        await apiClient.put(`/workflows/${workflowId}`, { graph_data: graphData });
        message.success('保存成功');
        // 保存成功后清除草稿
        const draftKey = `workflow_draft_${window.location.pathname}`;
        localStorage.removeItem(draftKey);
      } else {
        message.info('请先创建流程模板');
      }
    } catch (error) {
      message.error('保存失败');
    }
  };

  const handlePublish = async () => {
    if (!workflowId) {
      message.info('请先保存流程');
      return;
    }
    try {
      await apiClient.post(`/workflows/${workflowId}/publish`);
      message.success('发布成功');
    } catch (error) {
      message.error('发布失败');
    }
  };

  const onNodeClick = useCallback((_: any, node: any) => {
    console.log('Node clicked:', node);
    setSelectedNode(node);
    form.setFieldsValue({
      label: node.data?.label || '',
      approverId: node.data?.approverId || undefined,
      approverName: node.data?.approverName || '',
      departmentId: node.data?.departmentId || undefined,
      condition: node.data?.condition || '',
      trueLabel: node.data?.trueLabel || '是',
      falseLabel: node.data?.falseLabel || '否',
    });
    // 加载部门列表
    departmentApi.list().then(res => {
      setDepartmentOptions(res.data || []);
    }).catch(() => {});
    setPropertyModalOpen(true);
  }, [form]);

  const handlePropertyUpdate = () => {
    if (!selectedNode) {return;}
    const values = form.getFieldsValue();
    // 互斥处理：部门和直接审批人只能选一个
    const nodeData: any = { ...selectedNode.data };
    if (values.approverId) {
      nodeData.approverId = values.approverId;
      nodeData.approverName = values.approverName || '';
      nodeData.departmentId = undefined;
    } else if (values.departmentId) {
      nodeData.departmentId = values.departmentId;
      nodeData.approverId = undefined;
      nodeData.approverName = '';
    }
    nodeData.label = values.label;
    setNodes((nds: any[]) =>
      nds.map((node: any) =>
        node.id === selectedNode.id ? { ...node, data: nodeData } : node,
      ),
    );
    setPropertyModalOpen(false);
    message.success('属性已更新');
  };

  const addNode = (type: string) => {
    console.log('addNode called:', type);
    const id = `${type}_${Date.now()}`;
    const position = { x: 250, y: 150 };

    const nodeConfig: Record<string, any> = {
      start: {
        type: 'start',
        data: { label: '开始' },
      },
      end: {
        type: 'end',
        data: { label: '结束' },
      },
      approval: {
        type: 'approval',
        data: { label: '审批节点', approverName: '' },
      },
      condition: {
        type: 'condition',
        data: { label: '条件节点', condition: '', trueLabel: '是', falseLabel: '否' },
      },
      countersign: {
        type: 'countersign',
        data: { label: '会签节点', approvers: [], requiredCount: 1 },
      },
    };

    const config = nodeConfig[type];
    if (!config) {return;}

    const newNode = {
      id,
      position,
      ...config,
    };

    console.log('Creating node:', newNode);
    setNodes((nds) => {
      console.log('Previous nodes:', nds);
      const updated = [...nds, newNode];
      console.log('Updated nodes:', updated);
      return updated;
    });
  };

  const onDragOver = useCallback((event: React.DragEvent) => {
    event.preventDefault();
    event.dataTransfer.dropEffect = 'move';
  }, []);

  const onDrop = useCallback(
    (event: React.DragEvent) => {
      event.preventDefault();
      const type = event.dataTransfer.getData('application/reactflow');
      if (!type) {return;}

      const reactFlowBounds = event.currentTarget.getBoundingClientRect();
      const position = {
        x: event.clientX - reactFlowBounds.left,
        y: event.clientY - reactFlowBounds.top,
      };

      const id = `${type}_${Date.now()}`;
      const nodeConfig: Record<string, any> = {
        start: { type: 'start', data: { label: '开始' } },
        end: { type: 'end', data: { label: '结束' } },
        approval: { type: 'approval', data: { label: '审批节点', approverName: '' } },
        condition: { type: 'condition', data: { label: '条件节点', condition: '', trueLabel: '是', falseLabel: '否' } },
      };

      const config = nodeConfig[type];
      if (!config) {return;}

      const newNode = { id, position, ...config };
      setNodes((nds) => [...nds, newNode]);
    },
    [],
  );

  const onDragStart = (event: React.DragEvent, nodeType: string) => {
    event.dataTransfer.setData('application/reactflow', nodeType);
    event.dataTransfer.effectAllowed = 'move';
  };

  const onNodesChange = useCallback((changes: any[]) => {
    setNodes((nds) => applyNodeChanges(changes, nds));
  }, []);

  const onEdgesChange = useCallback((changes: any[]) => {
    setEdges((eds) => applyEdgeChanges(changes, eds));
  }, []);

  const onInit = (_: any) => {
    console.log('ReactFlow initialized');
    console.log('Available node types:', Object.keys(nodeTypes));
  };

  const renderNodeProperties = () => {
    if (!selectedNode) {return null;}

    switch (selectedNode.type) {
      case 'start':
      case 'end':
        return (
          <>
            <Form.Item label="标签" name="label">
              <Input />
            </Form.Item>
          </>
        );
      case 'approval':
        return (
          <>
            <Form.Item label="标签" name="label">
              <Input />
            </Form.Item>
            <Form.Item label="指定审批人" name="approverId">
              <Select
                showSearch
                placeholder="搜索选择审批人（直接指定）"
                allowClear
                filterOption={false}
                onSearch={async (value) => {
                  if (value.length < 1) {
                    setEmployeeOptions([]);
                    return;
                  }
                  setSearchingEmployee(true);
                  try {
                    const response = await employeeApi.search(value);
                    setEmployeeOptions(response.data);
                  } catch {
                    message.error('搜索员工失败');
                  } finally {
                    setSearchingEmployee(false);
                  }
                }}
                onChange={(value) => {
                  const emp = employeeOptions.find(e => e.id === value);
                  form.setFieldsValue({ approverName: emp?.name || emp?.username || '' });
                }}
                loading={searchingEmployee}
                options={employeeOptions.map(e => ({
                  label: `${e.name || e.username} (${e.email})`,
                  value: e.id,
                }))}
              />
            </Form.Item>
            <Form.Item name="approverName" hidden>
              <Input />
            </Form.Item>
            <Form.Item label="或按部门负责人" name="departmentId">
              <Select
                showSearch
                placeholder="选择部门（自动使用部门负责人）"
                allowClear
                optionFilterProp="label"
                options={departmentOptions.map(d => ({
                  label: d.name,
                  value: d.id,
                }))}
                onChange={(value) => {
                  if (value) {
                    form.setFieldsValue({ approverId: undefined });
                  }
                }}
              />
            </Form.Item>
          </>
        );
      case 'condition':
        return (
          <>
            <Form.Item label="条件描述" name="condition">
              <Input.TextArea placeholder="请输入条件描述" />
            </Form.Item>
            <Form.Item label="是分支标签" name="trueLabel">
              <Input placeholder="如：通过" />
            </Form.Item>
            <Form.Item label="否分支标签" name="falseLabel">
              <Input placeholder="如：拒绝" />
            </Form.Item>
          </>
        );
      default:
        return null;
    }
  };

  return (
    <div style={{ height: 'calc(100vh - 120px)', display: 'flex', flexDirection: 'column' }}>
      <div style={{ padding: '12px 16px', borderBottom: '1px solid #eee', display: 'flex', justifyContent: 'space-between', alignItems: 'center', background: '#fff' }}>
        <span style={{ fontSize: '16px', fontWeight: 'bold' }}>流程设计器</span>
        <Space>
          <Button onClick={() => addNode('start')}>添加开始</Button>
          <Button onClick={() => addNode('approval')}>添加审批</Button>
          <Button onClick={() => addNode('condition')}>添加条件</Button>
          <Button onClick={() => addNode('countersign')}>添加会签</Button>
          <Button onClick={() => addNode('end')}>添加结束</Button>
          <Button type="primary" onClick={handleSave}>保存</Button>
          <Button onClick={handlePublish}>发布</Button>
          <Button onClick={clearDraft}>清除草稿</Button>
        </Space>
      </div>
      <div style={{ flex: 1, position: 'relative', minHeight: '400px' }}>
        <ReactFlow
          nodes={nodes}
          edges={edges}
          onNodesChange={onNodesChange}
          onEdgesChange={onEdgesChange}
          onConnect={onConnect}
          onInit={onInit}
          onNodeClick={onNodeClick}
          onDragOver={onDragOver}
          onDrop={onDrop}
          nodeTypes={nodeTypes}
          defaultViewport={{ x: 0, y: 0, zoom: 1 }}
        >
          <Controls />
          <Background />
          <Panel position="top-left">
            <div style={{ background: '#fff', padding: '10px', borderRadius: '4px', boxShadow: '0 2px 8px rgba(0,0,0,0.15)' }}>
              <div style={{ fontSize: '12px', marginBottom: '8px', fontWeight: 'bold' }}>节点面板</div>
              <Space direction="vertical" size="small">
                <div
                  draggable
                  onDragStart={(e) => onDragStart(e, 'start')}
                  style={{
                    padding: '8px 16px',
                    background: '#52c41a',
                    color: 'white',
                    borderRadius: '20px',
                    cursor: 'move',
                    textAlign: 'center',
                  }}
                >
                  开始节点
                </div>
                <div
                  draggable
                  onDragStart={(e) => onDragStart(e, 'approval')}
                  style={{
                    padding: '8px 16px',
                    background: '#1890ff',
                    color: 'white',
                    borderRadius: '8px',
                    cursor: 'move',
                    textAlign: 'center',
                  }}
                >
                  审批节点
                </div>
                <div
                  draggable
                  onDragStart={(e) => onDragStart(e, 'condition')}
                  style={{
                    padding: '8px 16px',
                    background: '#faad14',
                    color: 'white',
                    borderRadius: '8px',
                    cursor: 'move',
                    textAlign: 'center',
                  }}
                >
                  条件节点
                </div>
                <div
                  draggable
                  onDragStart={(e) => onDragStart(e, 'countersign')}
                  style={{
                    padding: '8px 16px',
                    background: '#722ed1',
                    color: 'white',
                    borderRadius: '8px',
                    cursor: 'move',
                    textAlign: 'center',
                  }}
                >
                  会签节点
                </div>
                <div
                  draggable
                  onDragStart={(e) => onDragStart(e, 'end')}
                  style={{
                    padding: '8px 16px',
                    background: '#ff4d4f',
                    color: 'white',
                    borderRadius: '20px',
                    cursor: 'move',
                    textAlign: 'center',
                  }}
                >
                  结束节点
                </div>
              </Space>
            </div>
          </Panel>
        </ReactFlow>
      </div>

      <Modal
        title="节点属性"
        open={propertyModalOpen}
        onOk={handlePropertyUpdate}
        onCancel={() => setPropertyModalOpen(false)}
        okText="确定"
        cancelText="取消"
      >
        <Form form={form} layout="vertical">
          {renderNodeProperties()}
        </Form>
      </Modal>
    </div>
  );
};
