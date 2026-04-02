import { Handle, Position } from '@xyflow/react';

export function CountersignNode({ data, selected }: any) {
  const approvers = data?.approvers || [];
  return (
    <div style={{
      padding: '15px',
      background: '#722ed1',
      borderRadius: '8px',
      color: 'white',
      minWidth: '120px',
      boxShadow: selected ? '0 0 10px rgba(114,46,209,0.5)' : 'none',
    }}>
      <Handle type="target" position={Position.Left} />
      <div style={{ fontWeight: 'bold', marginBottom: '5px' }}>会签节点</div>
      <div style={{ fontSize: '11px', marginBottom: '5px' }}>
        审批人: {approvers.length > 0 ? approvers.join(', ') : '未配置'}
      </div>
      <div style={{ fontSize: '10px', opacity: 0.8 }}>
        需要 {data?.requiredCount || approvers.length} 人全部同意
      </div>
      <Handle type="source" position={Position.Right} />
    </div>
  );
}
