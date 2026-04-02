import { apiClient } from './client';

export interface SystemSettings {
  smtp_host: string
  smtp_port: string
  smtp_user: string
  smtp_password: string
  smtp_from: string
}

export const systemSettingApi = {
  list: () => apiClient.get<SystemSettings>('/system-settings'),
  get: (key: string) => apiClient.get<{ key: string; value: string }>(`/system-settings/${key}`),
  set: (key: string, value: string) =>
    apiClient.put(`/system-settings/${key}`, { value }),
  delete: (key: string) => apiClient.delete(`/system-settings/${key}`),
};
