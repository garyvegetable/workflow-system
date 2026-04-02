import { apiClient } from './client';

export const approvalApi = {
  listPending: () => apiClient.get('/tasks/pending'),
  listHandled: () => apiClient.get('/tasks/handled'),
  approve: (taskId: number, data: { comment?: string }) =>
    apiClient.post(`/tasks/${taskId}/approve`, data),
  reject: (taskId: number, data: { comment: string }) =>
    apiClient.post(`/tasks/${taskId}/reject`, data),
  transfer: (taskId: number, data: { new_assignee_id: number; comment?: string }) =>
    apiClient.post(`/tasks/${taskId}/transfer`, data),
  getHistory: (instanceId: number) => apiClient.get('/tasks/history', { params: { instance_id: instanceId } }),
  batchApprove: (taskIds: number[], comment?: string) =>
    apiClient.post('/tasks/batch-approve', { task_ids: taskIds, comment }),
  batchReject: (taskIds: number[], comment?: string) =>
    apiClient.post('/tasks/batch-reject', { task_ids: taskIds, comment }),
  addApprover: (taskId: number, newApproverId: number) =>
    apiClient.post(`/tasks/${taskId}/add-approver`, { new_approver_id: newApproverId }),
  removeApprover: (taskId: number, targetAssigneeId: number) =>
    apiClient.post(`/tasks/${taskId}/remove-approver`, { target_assignee_id: targetAssigneeId }),
};
