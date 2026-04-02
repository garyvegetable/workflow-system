import { apiClient } from './client';

export interface ExpenseCategory {
  id: number
  code: string
  name: string
  parent_id?: number | null
  status: number
  children?: ExpenseCategory[]
}

export const expenseCategoryApi = {
  list: (companyId?: number) =>
    apiClient.get<ExpenseCategory[]>('/expense-categories', {
      params: companyId ? { company_id: companyId } : {},
    }),

  create: (data: Partial<ExpenseCategory>) =>
    apiClient.post<ExpenseCategory>('/expense-categories', data),

  update: (id: number, data: Partial<ExpenseCategory>) =>
    apiClient.put<ExpenseCategory>(`/expense-categories/${id}`, data),

  delete: (id: number) =>
    apiClient.delete(`/expense-categories/${id}`),
};