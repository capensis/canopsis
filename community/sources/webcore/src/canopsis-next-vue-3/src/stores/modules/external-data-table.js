import { API_ROUTES } from '@/config';
import request from '@/services/request';
import { createCRUDStore } from '@/stores/plugins/crud';

const convertObjectToFormData = (data = {}) =>
  Object.entries(data).reduce((formData, [key, value]) => {
    formData.append(key, value);

    return formData;
  }, new FormData());

export const useExternalDataTableStore = createCRUDStore({
  id: 'externalDataTable',
  route: API_ROUTES.externalDataTables,
  withWithoutStore: true,
  actions: {
    fetchItemWithoutStore({ id } = {}) {
      return request.get(`${API_ROUTES.externalDataTables}/${id}`);
    },

    fetchSchema({ id } = {}) {
      return request.get(`${API_ROUTES.externalDataTables}/${id}/schema`);
    },

    createExport({ id, data } = {}) {
      return request.post(`${API_ROUTES.externalDataExport}/${id}`, data);
    },

    fetchExportStatus({ id } = {}) {
      return request.get(`${API_ROUTES.externalDataExport}/${id}`);
    },

    createImport({ id, data } = {}) {
      return request.post(`${API_ROUTES.externalDataImport}/${id}`, convertObjectToFormData(data), {
        headers: { 'Content-Type': 'multipart/form-data' },
      });
    },

    fetchImportData({ id, params } = {}) {
      return request.get(`${API_ROUTES.externalDataImport}/${id}/data`, {
        params,
      });
    },

    fetchImportStatus({ id } = {}) {
      return request.get(`${API_ROUTES.externalDataImport}/${id}/status`);
    },

    previewImport({ id, data } = {}) {
      return request.put(`${API_ROUTES.externalDataImport}/${id}/preview`, data);
    },

    completeImport({ id, data } = {}) {
      return request.put(`${API_ROUTES.externalDataImport}/${id}/complete`, data);
    },

    fetchRecordsListWithoutStore({ id, params } = {}) {
      return request.get(`${API_ROUTES.externalDataTables}/${id}/data`, { params });
    },

    createRecord({ table, data } = {}) {
      return request.post(`${API_ROUTES.externalDataTables}/${table}/data`, data);
    },

    updateRecord({ table, id, data } = {}) {
      return request.put(`${API_ROUTES.externalDataTables}/${table}/data/${id}`, data);
    },

    removeRecord({ table, id } = {}) {
      return request.delete(`${API_ROUTES.externalDataTables}/${table}/data/${id}`);
    },

    bulkRemoveRecords({ table, data } = {}) {
      return request.delete(`${API_ROUTES.bulkExternalDataTables}/${table}/data`, { data });
    },
  },
});
