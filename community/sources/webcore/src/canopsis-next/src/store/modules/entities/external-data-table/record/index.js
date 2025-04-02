import { API_ROUTES } from '@/config';

import request from '@/services/request';

export default {
  namespaced: true,
  actions: {
    fetchListWithoutStore(context, { id, params } = {}) {
      return request.get(`${API_ROUTES.externalDataTables}/${id}/data`, { params });
    },

    create(context, { table, data }) {
      return request.post(`${API_ROUTES.externalDataTables}/${table}/data`, data);
    },

    update(context, { table, id, data }) {
      return request.put(`${API_ROUTES.externalDataTables}/${table}/data/${id}`, data);
    },

    remove(context, { table, id } = {}) {
      return request.delete(`${API_ROUTES.externalDataTables}/${table}/data/${id}`);
    },

    bulkRemove(context, { table, data }) {
      return request.delete(`${API_ROUTES.bulkExternalDataTables}/${table}/data`, { data });
    },
  },
};
