import { API_ROUTES } from '@/config';
import request from '@/services/request';

export default {
  namespaced: true,
  actions: {
    fetchListWithoutStore(context, { id, params }) {
      return request.get(`${API_ROUTES.pbehavior.pbehaviors}/${id}/entities`, { params });
    },
  },
};
