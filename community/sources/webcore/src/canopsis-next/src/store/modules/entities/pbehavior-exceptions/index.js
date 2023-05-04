import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import request from '@/services/request';

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
      const formData = Object.entries(data).reduce((acc, [key, value]) => {
        acc.append(key, value);
        return acc;
      }, new FormData());

      return request.post(API_ROUTES.pbehavior.exceptionImport, formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
      });
    },
  },
});
