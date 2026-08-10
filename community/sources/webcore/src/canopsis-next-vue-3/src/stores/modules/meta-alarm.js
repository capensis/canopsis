import { defineStore } from 'pinia';

import { API_ROUTES } from '@/config';
import request from '@/services/request';

export const useMetaAlarmStore = defineStore('metaAlarm', () => {
  const fetchListWithoutStore = ({ params } = {}) => request.get(API_ROUTES.metaAlarm, { params });
  const create = ({ data } = {}) => request.post(API_ROUTES.metaAlarm, data);
  const addAlarms = ({ id, data } = {}) => request.put(`${API_ROUTES.metaAlarm}/${id}/add`, data);
  const removeAlarms = ({ id, data } = {}) =>
    request.put(`${API_ROUTES.metaAlarm}/${id}/remove`, data);

  return {
    fetchListWithoutStore,
    create,
    addAlarms,
    removeAlarms,
  };
});
