import { Handle, Position } from '@xyflow/react';

export function ApprovalNode({ data, selected }: any) {
  return (
    <div style={{
      padding: '15px 20px',
      background: '#1890ff',
      borderRadius: '8px',
      color: 'white',
      minWidth: '120px',
      boxShadow: selected ? '0 0 10px rgba(24,144,255,0.5)' : 'none',
    }}>
      <Handle type="target" position={Position.Left} />
      <div style={{ fontWeight: 'bold', marginBottom: '5px' }}>审批节点</div>
      <div style={{ fontSize: '12px', opacity: 0.9 }}>
        {data?.approverName || '未配置'}
      </div>
      <Handle type="source" position={Position.Right} />
    </div>
  );
}
