import { defineStore } from 'pinia';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

export const useExternalDataTableStore = defineStore('externalDataTable', () => {
  const fetchListWithoutStore = ({ params } = {}) => (
    request.get(API_ROUTES.externalDataTables, { params })
  );

  const fetchWithoutStore = ({ id } = {}) => (
    request.get(`${API_ROUTES.externalDataTables}/${id}`)
  );

  const fetchSchema = ({ id } = {}) => (
    request.get(`${API_ROUTES.externalDataTables}/${id}/schema`)
  );

  return {
    fetchListWithoutStore,
    fetchWithoutStore,
    fetchSchema,
  };
});
