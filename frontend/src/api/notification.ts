import { apiClient } from './client';

export const notificationApi = {
  list: () => apiClient.get('/notifications'),
  markAsRead: (id: number) => apiClient.put(`/notifications/${id}/read`),
  markAllRead: () => apiClient.put('/notifications/read-all'),
};
