import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import request from '@/services/request';

export default {
  namespaced: true,
  getters: {
    getItemById: (state, getters, rootState, rootGetters) => id => rootGetters['entities/getItem'](
      ENTITIES_TYPES.widget,
      id,
    ),
  },
  actions: {
    create(context, { data } = {}) {
      return request.post(API_ROUTES.widget, data);
    },

    clone(context, { data, id } = {}) {
      return request.post(`${API_ROUTES.widget}/${id}/clone`, data);
    },

    update(context, { data, id } = {}) {
      return request.put(`${API_ROUTES.widget}/${id}`, data);
    },

    remove(context, { id } = {}) {
      return request.delete(`${API_ROUTES.widget}/${id}`);
    },

    copy(context, { id, data } = {}) {
      return request.post(`${API_ROUTES.widgetCopy}/${id}`, data);
    },

    updateGridPositions(context, { data } = {}) {
      return request.put(API_ROUTES.widgetGridPosition, data);
    },
  },
};
