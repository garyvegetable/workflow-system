import { Handle, Position } from '@xyflow/react';

export function ConditionNode({ data, selected }: any) {
  return (
    <div style={{
      padding: '15px',
      background: '#faad14',
      borderRadius: '8px',
      color: 'white',
      minWidth: '100px',
      boxShadow: selected ? '0 0 10px rgba(250,173,20,0.5)' : 'none',
    }}>
      <Handle type="target" position={Position.Left} />
      <div style={{ fontWeight: 'bold', marginBottom: '5px' }}>条件节点</div>
      <div style={{ fontSize: '11px' }}>
        {data?.condition || '请配置条件'}
      </div>
      <div style={{ marginTop: '10px', fontSize: '11px' }}>
        <span style={{ color: '#52c41a' }}>是 → {data?.trueLabel || 'Y'}</span><br/>
        <span style={{ color: '#ff4d4f' }}>否 → {data?.falseLabel || 'N'}</span>
      </div>
      {/* Yes branch handle - top right */}
      <Handle
        type="source"
        position={Position.Right}
        id="yes"
        style={{ top: '25%', background: '#52c41a' }}
      />
      {/* No branch handle - bottom right */}
      <Handle
        type="source"
        position={Position.Right}
        id="no"
        style={{ top: '75%', background: '#ff4d4f' }}
      />
    </div>
  );
}
