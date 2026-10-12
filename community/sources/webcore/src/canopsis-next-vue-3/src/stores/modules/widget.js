import { defineStore } from 'pinia';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

export const useWidgetStore = defineStore('widget', () => {
  const fetchItemWithoutStore = ({ id } = {}) => request.get(`${API_ROUTES.widget.list}/${id}`);
  const create = ({ data } = {}) => request.post(API_ROUTES.widget.list, data);
  const clone = ({ data, id } = {}) => request.post(`${API_ROUTES.widget.list}/${id}/clone`, data);
  const update = ({ data, id } = {}) => request.put(`${API_ROUTES.widget.list}/${id}`, data);
  const remove = ({ id } = {}) => request.delete(`${API_ROUTES.widget.list}/${id}`);
  const updateGridPositions = ({ data } = {}) => request.put(API_ROUTES.widget.gridPositions, data);
  const fetchWidgetFilters = ({ params } = {}) =>
    request.get(API_ROUTES.widget.filters, { params });
  const createWidgetFilter = ({ data } = {}) => request.post(API_ROUTES.widget.filters, data);
  const updateWidgetFilter = ({ id, data } = {}) =>
    request.put(`${API_ROUTES.widget.filters}/${id}`, data);
  const removeWidgetFilter = ({ id, data } = {}) =>
    request.delete(`${API_ROUTES.widget.filters}/${id}`, { data });
  const fetchWidgetFilter = ({ id, params } = {}) =>
    request.get(`${API_ROUTES.widget.filters}/${id}`, { params });
  const updateWidgetFiltersPositions = ({ data } = {}) =>
    request.put(API_ROUTES.widget.filterPositions, data);

  return {
    fetchItemWithoutStore,
    create,
    clone,
    update,
    remove,
    updateGridPositions,
    fetchWidgetFilters,
    createWidgetFilter,
    updateWidgetFilter,
    removeWidgetFilter,
    fetchWidgetFilter,
    updateWidgetFiltersPositions,
  };
});
