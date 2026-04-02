import { Handle, Position } from '@xyflow/react';

export function EndNode({ data, selected }: any) {
  return (
    <div style={{
      padding: '10px 20px',
      background: '#ff4d4f',
      borderRadius: '20px',
      color: 'white',
      fontWeight: 'bold',
      boxShadow: selected ? '0 0 10px rgba(255,77,79,0.5)' : 'none',
    }}>
      <Handle type="target" position={Position.Left} />
      {data?.label || '结束'}
    </div>
  );
}
