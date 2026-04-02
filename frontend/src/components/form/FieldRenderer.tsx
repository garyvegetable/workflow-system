import { Input, Select, DatePicker, Radio, Checkbox, Upload, InputNumber } from 'antd';
import { InboxOutlined } from '@ant-design/icons';
import type { FormField } from './FieldEditor';

interface Props {
  field: FormField
  value?: any
  onChange?: (value: any) => void
}

const { TextArea } = Input;

export const FieldRenderer = ({ field, value, onChange }: Props) => {
  switch (field.type) {
    case 'text':
      return <Input
        placeholder={field.placeholder}
        value={value}
        onChange={(e) => onChange?.(e.target.value)}
        required={field.required}
      />;
    case 'textarea':
      return <TextArea
        placeholder={field.placeholder}
        value={value}
        onChange={(e) => onChange?.(e.target.value)}
        rows={3}
      />;
    case 'number':
      return <InputNumber
        style={{ width: '100%' }}
        value={value}
        onChange={onChange}
      />;
    case 'date':
      return <DatePicker style={{ width: '100%' }} value={value} onChange={onChange} />;
    case 'select':
      return <Select
        placeholder={field.placeholder}
        value={value}
        onChange={onChange}
        options={field.options}
        style={{ width: '100%' }}
      />;
    case 'radio':
      return <Radio.Group value={value} onChange={(e) => onChange?.(e.target.value)}>
        {field.options?.map(opt => (
          <Radio key={opt.value} value={opt.value}>{opt.label}</Radio>
        ))}
      </Radio.Group>;
    case 'checkbox':
      return <Checkbox.Group
        value={value}
        onChange={onChange}
        options={field.options}
      />;
    case 'file':
      return <Upload.Dragger>
        <p className="ant-upload-drag-icon"><InboxOutlined /></p>
        <p>点击或拖拽文件到此区域上传</p>
      </Upload.Dragger>;
    default:
      return <Input value={value} onChange={(e) => onChange?.(e.target.value)} />;
  }
};
