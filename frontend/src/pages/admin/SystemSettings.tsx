import { useState, useEffect } from 'react';
import { Card, Form, Input, Button, message, Space } from 'antd';
import { systemSettingApi } from '@/api/systemSetting';

export const SystemSettingsPage = () => {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    fetchSettings();
  }, []);

  const fetchSettings = async () => {
    setLoading(true);
    try {
      const response = await systemSettingApi.list();
      form.setFieldsValue(response.data);
    } catch (error) {
      message.error('获取设置失败');
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async () => {
    try {
      const values = await form.validateFields();
      setSaving(true);
      // 保存每个设置
      for (const [key, value] of Object.entries(values)) {
        if (value !== undefined && value !== null) {
          await systemSettingApi.set(key, String(value));
        }
      }
      message.success('保存成功');
    } catch (error) {
      message.error('保存失败');
    } finally {
      setSaving(false);
    }
  };

  return (
    <Card title="系统设置" loading={loading}>
      <Form
        form={form}
        layout="vertical"
        initialValues={{
          smtp_host: '',
          smtp_port: '587',
          smtp_user: '',
          smtp_password: '',
          smtp_from: '',
        }}
      >
        <Form.Item label="SMTP 服务器地址" name="smtp_host">
          <Input placeholder="如：smtp.qq.com" />
        </Form.Item>

        <Form.Item label="SMTP 端口" name="smtp_port">
          <Input placeholder="如：587" />
        </Form.Item>

        <Form.Item label="SMTP 用户名" name="smtp_user">
          <Input placeholder="邮箱地址" />
        </Form.Item>

        <Form.Item label="SMTP 密码" name="smtp_password">
          <Input.Password placeholder="邮箱密码或授权码" />
        </Form.Item>

        <Form.Item label="发件人邮箱" name="smtp_from">
          <Input placeholder="如：noreply@example.com" />
        </Form.Item>

        <Space>
          <Button type="primary" onClick={handleSave} loading={saving}>
            保存设置
          </Button>
          <Button onClick={fetchSettings}>
            重置
          </Button>
        </Space>
      </Form>

      <Card title="说明" style={{ marginTop: 24 }}>
        <ul style={{ paddingLeft: 20 }}>
          <li>SMTP 用于发送邮件通知，如审批提醒、审批结果通知等</li>
          <li>配置后需要重启服务生效</li>
          <li>推荐使用 QQ 邮箱或企业邮箱，发送前请确保已开启 SMTP 服务</li>
        </ul>
      </Card>
    </Card>
  );
};
