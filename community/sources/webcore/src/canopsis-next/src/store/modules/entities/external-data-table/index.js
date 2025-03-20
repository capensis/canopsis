import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

import { convertObjectToFormData } from '@/helpers/request';

export default createCRUDModule({
  route: API_ROUTES.externalDataTables,
  withWithoutStore: true,
}, {
  actions: {
    fetchItemWithoutStore(context, { id }) {
      return request.get(`${API_ROUTES.externalDataTables}/${id}`);
    },

    /**
     * DATA
     */
    fetchData(context, { id } = {}) {
      return request.get(`${API_ROUTES.externalDataTables}/${id}/data`);
    },

    createData(context, { id, data } = {}) {
      return request.post(`${API_ROUTES.externalDataTables}/${id}/data`, data);
    },

    /**
     * SCHEMA
     */
    fetchSchema(context, { id } = {}) {
      return request.get(`${API_ROUTES.externalDataTables}/${id}/schema`);
    },

    /**
     * IMPORTS
     */
    createImport(context, { id, data } = {}) {
      return request.post(`${API_ROUTES.externalDataImport}/${id}`, convertObjectToFormData(data), {
        headers: { 'Content-Type': 'multipart/form-data' },
      });
    },

    completeImport(context, { id, data } = {}) {
      return request.put(`${API_ROUTES.externalDataImport}/${id}/complete`, data);
    },

    fetchImportData(context, { id } = {}) {
      return request.get(`${API_ROUTES.externalDataImport}/${id}/data`);
    },

    fetchImportStatus(context, { id } = {}) {
      return request.get(`${API_ROUTES.externalDataImport}/${id}/status`);
    },
  },
});
