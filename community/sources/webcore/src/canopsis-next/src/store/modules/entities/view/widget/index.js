import { API_ROUTES } from '@/config';
import request from '@/services/request';

export default {
  namespaced: true,
  actions: {
    create(context, { data }) {
      return request.post(API_ROUTES.widget, data);
    },

    update(context, { data, id }) {
      return request.put(`${API_ROUTES.widget}/${id}`, data);
    },

    remove(context, { id }) {
      return request.delete(`${API_ROUTES.widget}/${id}`);
    },

    updatePositions(context, { data }) {
      return request.delete(API_ROUTES.widgetPosition, data);
    },
  },
};
