import { Handle, Position } from '@xyflow/react';

export function StartNode({ data, selected }: any) {
  console.log('StartNode rendered with:', { data, selected });
  return (
    <div style={{
      padding: '10px 20px',
      background: '#52c41a',
      borderRadius: '20px',
      color: 'white',
      fontWeight: 'bold',
      boxShadow: selected ? '0 0 10px rgba(82,196,26,0.5)' : 'none',
    }}>
      <Handle type="source" position={Position.Right} />
      {data?.label || '开始'}
    </div>
  );
}
