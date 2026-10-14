import { defineStore } from 'pinia';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

export const useDataStorageStore = defineStore('dataStorage', () => {
  const fetchItemWithoutStore = () => request.get(API_ROUTES.dataStorage);
  const update = ({ data } = {}) => request.put(API_ROUTES.dataStorage, data);

  return {
    fetchItemWithoutStore,
    update,
  };
});
