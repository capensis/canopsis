import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import request from '@/services/request';

import { widgetSchema } from '@/store/schemas';

export default {
  namespaced: true,
  getters: {
    getItemById: (state, getters, rootState, rootGetters) => id => rootGetters['entities/getItem'](
      ENTITIES_TYPES.widget,
      id,
    ),
  },
  actions: {
    fetchItem({ dispatch }, { id }) {
      return dispatch('entities/fetch', {
        route: `${API_ROUTES.widget.list}/${id}`,
        schema: widgetSchema,
      }, { root: true });
    },

    fetchItemWithoutStore(context, { id }) {
      return request.get(`${API_ROUTES.widget.list}/${id}`);
    },

    create(context, { data } = {}) {
      return request.post(API_ROUTES.widget.list, data);
    },

    clone(context, { data, id } = {}) {
      return request.post(`${API_ROUTES.widget.list}/${id}/clone`, data);
    },

    update(context, { data, id } = {}) {
      return request.put(`${API_ROUTES.widget.list}/${id}`, data);
    },

    remove(context, { id } = {}) {
      return request.delete(`${API_ROUTES.widget.list}/${id}`);
    },

    copy(context, { id, data } = {}) {
      return request.post(`${API_ROUTES.widget.copy}/${id}`, data);
    },

    updateGridPositions(context, { data } = {}) {
      return request.put(API_ROUTES.widget.gridPositions, data);
    },

    fetchWidgetFilters(context, { params } = {}) {
      return request.get(API_ROUTES.widget.filters, { params });
    },

    createWidgetFilter(context, { data } = {}) {
      return request.post(API_ROUTES.widget.filters, data);
    },

    updateWidgetFilter(context, { id, data } = {}) {
      return request.put(`${API_ROUTES.widget.filters}/${id}`, data);
    },

    removeWidgetFilter(context, { id, data } = {}) {
      return request.delete(`${API_ROUTES.widget.filters}/${id}`, data);
    },

    fetchWidgetFilter(context, { id, data } = {}) {
      return request.get(`${API_ROUTES.widget.filters}/${id}`, data);
    },

    updateWidgetFiltersPositions(context, { data } = {}) {
      return request.put(API_ROUTES.widget.filterPositions, data);
    },
  },
};
