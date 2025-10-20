import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { convertObjectToFormData } from '@/helpers/request';

export default {
  namespaced: true,
  actions: {
    create(context, { id, data } = {}) {
      return request.post(`${API_ROUTES.externalDataImport}/${id}`, convertObjectToFormData(data), {
        headers: { 'Content-Type': 'multipart/form-data' },
      });
    },

    fetchData(context, { id, params } = {}) {
      return request.get(`${API_ROUTES.externalDataImport}/${id}/data`, { params });
    },

    fetchStatus(context, { id } = {}) {
      return request.get(`${API_ROUTES.externalDataImport}/${id}/status`);
    },

    preview(context, { id, data } = {}) {
      return request.put(`${API_ROUTES.externalDataImport}/${id}/preview`, data);
    },

    complete(context, { id, data } = {}) {
      return request.put(`${API_ROUTES.externalDataImport}/${id}/complete`, data);
    },
  },
};
