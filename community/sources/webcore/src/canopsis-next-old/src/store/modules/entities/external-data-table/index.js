import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

import recordModule from './record';
import importModule from './import';

export default createCRUDModule({
  route: API_ROUTES.externalDataTables,
  withWithoutStore: true,
}, {
  modules: {
    record: recordModule,
    import: importModule,
  },
  actions: {
    fetchItemWithoutStore(context, { id }) {
      return request.get(`${API_ROUTES.externalDataTables}/${id}`);
    },

    /**
     * SCHEMA
     */
    fetchSchema(context, { id } = {}) {
      return request.get(`${API_ROUTES.externalDataTables}/${id}/schema`);
    },

    /**
     * EXPORT
     */
    createExport(context, { id, data } = {}) {
      return request.post(`${API_ROUTES.externalDataExport}/${id}`, data);
    },

    fetchExportStatus(context, { id } = {}) {
      return request.get(`${API_ROUTES.externalDataExport}/${id}`);
    },
  },
});
