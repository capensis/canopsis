import { defineStore } from 'pinia';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

export const useNotificationStore = defineStore('notification', () => {
  const fetchNotificationsListWithoutStore = ({ params } = {}) => request.get(API_ROUTES.notifications, { params });

  return {
    fetchNotificationsListWithoutStore,
  };
});
