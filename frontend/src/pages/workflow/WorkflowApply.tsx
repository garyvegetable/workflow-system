import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Card, Form, Button, Space, message, Upload } from 'antd';
import { InboxOutlined, SendOutlined } from '@ant-design/icons';
import { workflowApi } from '@/api/workflow';
import { DynamicForm } from '@/components/form/DynamicForm';
import type { FormField } from '@/components/form/FieldEditor';

const { Dragger } = Upload;

interface WorkflowDetail {
  id: number
  code: string
  name: string
  form_fields: FormField[]
}

export const WorkflowApply = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [workflow, setWorkflow] = useState<WorkflowDetail | null>(null);
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    if (id) {
      fetchWorkflow(parseInt(id));
    }
  }, [id]);

  const fetchWorkflow = async (workflowId: number) => {
    setLoading(true);
    try {
      const response = await workflowApi.get(workflowId);
      const wf = response.data;
      if (typeof wf.form_fields === 'string') {
        wf.form_fields = JSON.parse(wf.form_fields);
      }
      setWorkflow(wf);
    } catch {
      message.error('获取流程信息失败');
    } finally {
      setLoading(false);
    }
  };

  const uploadProps = {
    name: 'file',
    action: '/api/v1/attachments/upload',
    headers: {
      Authorization: `Bearer ${localStorage.getItem('token')}`,
    },
    onChange: (info: any) => {
      if (info.file.status === 'done') {
        message.success(`${info.file.name} 上传成功`);
      } else if (info.file.status === 'error') {
        message.error('上传失败');
      }
    },
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      setSubmitting(true);

      await workflowApi.createInstance({
        definition_id: workflow?.id,
        title: `${workflow?.name} - 申请`,
        form_data: values,
      });

      message.success('提交成功');
      navigate('/my-applications');
    } catch {
      message.error('提交失败');
    } finally {
      setSubmitting(false);
    }
  };

  if (!workflow) {
    return <Card loading={loading} />;
  }

  return (
    <Card
      title={`申请: ${workflow.name}`}
      extra={
        <Space>
          <Button onClick={() => navigate('/workflows')}>取消</Button>
          <Button type="primary" icon={<SendOutlined />} onClick={handleSubmit} loading={submitting}>
            提交申请
          </Button>
        </Space>
      }
    >
      <Form form={form} layout="vertical">
        <DynamicForm fields={workflow.form_fields || []} form={form} />

        <div style={{ marginTop: 24 }}>
          <h4>附件上传</h4>
          <Dragger {...uploadProps}>
            <p className="ant-upload-drag-icon"><InboxOutlined /></p>
            <p className="ant-upload-text">点击或拖拽文件到此区域上传</p>
            <p className="ant-upload-hint">支持 PDF、Excel、图片，单文件不超过 10MB</p>
          </Dragger>
        </div>
      </Form>
    </Card>
  );
};
