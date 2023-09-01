import { API_ROUTES } from '@/config';

import request from '@/services/request';

export default {
  namespaced: true,
  actions: {
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
