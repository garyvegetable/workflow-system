import { apiClient } from './client';

export interface Department {
  id: number
  name: string
  parent_id: number | null
  company_id: number
  leader_id?: number | null
  status: number
  children?: Department[]
}

export interface ApprovalChainStep {
  id?: number
  employee_id: number
  step_order: number
  employee_name?: string
}

export const departmentApi = {
  list: (companyId?: number) =>
    apiClient.get<Department[]>('/departments', {
      params: companyId ? { company_id: companyId } : {},
    }),

  create: (data: Partial<Department>) =>
    apiClient.post<Department>('/departments', data),

  update: (id: number, data: Partial<Department>) =>
    apiClient.put<Department>(`/departments/${id}`, data),

  delete: (id: number, transferToDeptId?: number) =>
    apiClient.delete(`/departments/${id}`, transferToDeptId ? { data: { transfer_to_dept_id: transferToDeptId } } : {}),

  getApprovalChain: (id: number) =>
    apiClient.get<ApprovalChainStep[]>(`/departments/${id}/approval-chain`),

  setApprovalChain: (id: number, steps: ApprovalChainStep[]) =>
    apiClient.put(`/departments/${id}/approval-chain`, { steps }),
};
