import { defineStore } from 'pinia';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

export const useNotificationSettingsStore = defineStore('notificationSettings', () => {
  const fetchItemWithoutStore = () => request.get(API_ROUTES.notificationSettings);
  const update = ({ data } = {}) => request.put(API_ROUTES.notificationSettings, data);

  return {
    fetchItemWithoutStore,
    update,
  };
});
