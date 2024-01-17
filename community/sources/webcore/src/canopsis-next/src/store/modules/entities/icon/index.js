import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

import { convertObjectToFormData } from '@/helpers/request';

export default createCRUDModule({
  route: API_ROUTES.icons,
  withFetchingParams: true,
}, {
  actions: {
    create(context, { data } = {}) {
      return request.post(API_ROUTES.icons, convertObjectToFormData(data), {
        headers: { 'Content-Type': 'multipart/form-data' },
      });
    },

    update(context, { id, data } = {}) {
      return request.put(`${API_ROUTES.icons}/${id}`, convertObjectToFormData(data), {
        headers: { 'Content-Type': 'multipart/form-data' },
      });
    },
  },
});
