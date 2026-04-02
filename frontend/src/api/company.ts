import { apiClient } from './client';

export const companyApi = {
  list: () => apiClient.get('/companies'),
  create: (data: any) => apiClient.post('/companies', data),
  update: (id: number, data: any) => apiClient.put(`/companies/${id}`, data),
  delete: (id: number) => apiClient.delete(`/companies/${id}`),
};
