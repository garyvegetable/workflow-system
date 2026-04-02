import { Form } from 'antd';
import { FieldRenderer } from './FieldRenderer';
import type { FormField } from './FieldEditor';

interface Props {
  fields: FormField[]
  form: any
  onValuesChange?: (changedValues: any, values: any) => void
}

export const DynamicForm = ({ fields, form, onValuesChange }: Props) => {
  return (
    <Form
      form={form}
      layout="vertical"
      onValuesChange={onValuesChange}
    >
      {fields.map((field) => (
        <Form.Item
          key={field.name}
          name={field.name}
          label={field.label}
          rules={[{ required: field.required, message: `请输入${field.label}` }]}
        >
          <FieldRenderer field={field} />
        </Form.Item>
      ))}
    </Form>
  );
};
