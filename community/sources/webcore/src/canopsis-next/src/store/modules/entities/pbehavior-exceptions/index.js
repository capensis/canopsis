import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import request from '@/services/request';

import { convertObjectToFormData } from '@/helpers/request';

import { createEntityModule } from '@/store/plugins/entities';

export default createEntityModule({
  route: API_ROUTES.pbehavior.exceptions,
  entityType: ENTITIES_TYPES.pbehaviorExceptions,
  dataPreparer: d => d.data,
  withFetchingParams: true,
  withMeta: true,
}, {
  actions: {
    fetchListWithoutStore(context, { params }) {
      return request.get(API_ROUTES.pbehavior.exceptions, { params });
    },

    import(context, { data } = {}) {
      return request.post(API_ROUTES.pbehavior.exceptionImport, convertObjectToFormData(data), {
        headers: { 'Content-Type': 'multipart/form-data' },
      });
    },
  },
});
