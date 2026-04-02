import { apiClient } from './client';

export const workflowApi = {
  list: (companyId?: number) => apiClient.get('/workflows', { params: { company_id: companyId } }),
  get: (id: number) => apiClient.get(`/workflows/${id}`),
  create: (data: any) => apiClient.post('/workflows', data),
  update: (id: number, data: any) => apiClient.put(`/workflows/${id}`, data),
  delete: (id: number) => apiClient.delete(`/workflows/${id}`),
  publish: (id: number) => apiClient.post(`/workflows/${id}/publish`),
  disable: (id: number) => apiClient.post(`/workflows/${id}/disable`),
  copy: (id: number, targetCompanyId: number) =>
    apiClient.post(`/workflows/${id}/copy`, { target_company_id: targetCompanyId }),
  // 申请提交
  createInstance: (data: any) => apiClient.post('/workflows/instances', data),
  getMyApplications: () => apiClient.get('/workflows/instances/my'),
  cancelInstance: (id: number) => apiClient.post(`/workflows/instances/${id}/cancel`),
};
