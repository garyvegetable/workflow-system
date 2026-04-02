import { apiClient } from './client';

export interface Supplier {
  id: number
  code: string
  name: string
  contact?: string
  phone?: string
  email?: string
  address?: string
  bank_name?: string
  bank_account?: string
  tax_number?: string
  status: number
}

export const supplierApi = {
  list: (companyId?: number) =>
    apiClient.get('/suppliers', { params: { company_id: companyId } }),
  create: (data: any) => apiClient.post('/suppliers', data),
  update: (id: number, data: any) => apiClient.put(`/suppliers/${id}`, data),
  delete: (id: number) => apiClient.delete(`/suppliers/${id}`),
};
