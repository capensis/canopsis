import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

import { convertObjectToFormData } from '@/helpers/request';

export default createCRUDModule({
  route: API_ROUTES.pbehavior.exceptions,
  withFetchingParams: true,
  withWithoutStore: true,
}, {
  actions: {
    import(context, { data } = {}) {
      return request.post(API_ROUTES.pbehavior.exceptionImport, convertObjectToFormData(data), {
        headers: { 'Content-Type': 'multipart/form-data' },
      });
    },
  },
});
