import { useState } from 'react';
import { Input, Select, Switch, Button, Card, Space } from 'antd';

const fieldTypes = [
  { label: '文本', value: 'text' },
  { label: '多行文本', value: 'textarea' },
  { label: '数字', value: 'number' },
  { label: '日期', value: 'date' },
  { label: '下拉选择', value: 'select' },
  { label: '单选', value: 'radio' },
  { label: '多选', value: 'checkbox' },
  { label: '文件上传', value: 'file' },
  { label: '明细表', value: 'table' },
];

export interface FormField {
  name: string
  type: string
  label: string
  required?: boolean
  options?: { label: string; value: string }[]
  placeholder?: string
  children?: FormField[]
}

interface Props {
  value?: FormField[]
  onChange?: (fields: FormField[]) => void
}

export const FieldEditor = ({ value = [], onChange }: Props) => {
  const [fields, setFields] = useState<FormField[]>(value);

  const addField = () => {
    const newField: FormField = {
      name: `field_${Date.now()}`,
      type: 'text',
      label: '新字段',
      required: false,
    };
    const updated = [...fields, newField];
    setFields(updated);
    onChange?.(updated);
  };

  const updateField = (index: number, updates: Partial<FormField>) => {
    const updated = fields.map((f, i) => i === index ? { ...f, ...updates } : f);
    setFields(updated);
    onChange?.(updated);
  };

  const removeField = (index: number) => {
    const updated = fields.filter((_, i) => i !== index);
    setFields(updated);
    onChange?.(updated);
  };

  return (
    <div>
      {fields.map((field, index) => (
        <Card key={field.name} size="small" style={{ marginBottom: 8 }}>
          <Space direction="vertical" style={{ width: '100%' }}>
            <Space>
              <Select
                value={field.type}
                onChange={(type) => updateField(index, { type })}
                options={fieldTypes}
                style={{ width: 120 }}
              />
              <Input
                placeholder="字段名"
                value={field.name}
                onChange={(e) => updateField(index, { name: e.target.value })}
                style={{ width: 120 }}
              />
              <Input
                placeholder="标签"
                value={field.label}
                onChange={(e) => updateField(index, { label: e.target.value })}
                style={{ width: 150 }}
              />
              <Switch
                checked={field.required}
                onChange={(checked) => updateField(index, { required: checked })}
                checkedChildren="必填"
                unCheckedChildren="选填"
              />
              <Button type="link" danger onClick={() => removeField(index)}>
                删除
              </Button>
            </Space>
          </Space>
        </Card>
      ))}
      <Button type="dashed" onClick={addField} block>
        + 添加字段
      </Button>
    </div>
  );
};
