import { API_ROUTES } from '@/config';
import request from '@/services/request';

export default {
  namespaced: true,
  actions: {
    fetchListWithoutStore(context, { id, params }) {
      return request.get(`${API_ROUTES.planning.pbehaviors}/${id}/eids`, { params });
    },
  },
};
