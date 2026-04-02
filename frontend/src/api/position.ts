import { apiClient } from './client';

export interface Position {
  id: number
  company_id: number
  name: string
  code: string
  status: number
}

export const positionApi = {
  list: (companyId?: number) =>
    apiClient.get<Position[]>('/positions', {
      params: companyId ? { company_id: companyId } : undefined,
    }),
  create: (data: Partial<Position>) => apiClient.post<Position>('/positions', data),
  get: (id: number) => apiClient.get<Position>(`/positions/${id}`),
  update: (id: number, data: Partial<Position>) =>
    apiClient.put<Position>(`/positions/${id}`, data),
  delete: (id: number) => apiClient.delete(`/positions/${id}`),
};
