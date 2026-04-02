import { apiClient } from './client';

export interface Employee {
  id: number
  company_id: number
  username: string
  name: string
  email: string
  level: string
  status: number
  password?: string
}

export interface BankAccount {
  id: number
  employee_id: number
  bank_name: string
  bank_branch: string
  bank_account: string
  account_holder: string
  is_default: boolean
}

export const employeeApi = {
  list: (companyId?: number) =>
    apiClient.get<Employee[]>('/employees', {
      params: companyId ? { company_id: companyId } : undefined,
    }),

  search: (name: string, companyId?: number) =>
    apiClient.get<Employee[]>('/employees/search', {
      params: { name, ...(companyId ? { company_id: companyId } : {}) },
    }),

  create: (data: Partial<Employee>) => apiClient.post<Employee>('/employees', data),

  get: (id: number) => apiClient.get<Employee>(`/employees/${id}`),

  update: (id: number, data: Partial<Employee>) =>
    apiClient.put<Employee>(`/employees/${id}`, data),

  delete: (id: number) => apiClient.delete(`/employees/${id}`),

  listBankAccounts: (id: number) =>
    apiClient.get<BankAccount[]>(`/employees/${id}/bank-accounts`),

  createBankAccount: (id: number, data: Partial<BankAccount>) =>
    apiClient.post<BankAccount>(`/employees/${id}/bank-accounts`, data),

  updateBankAccount: (id: number, aid: number, data: Partial<BankAccount>) =>
    apiClient.put<BankAccount>(`/employees/${id}/bank-accounts/${aid}`, data),

  deleteBankAccount: (id: number, aid: number) =>
    apiClient.delete(`/employees/${id}/bank-accounts/${aid}`),
};
