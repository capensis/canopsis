import { API_ROUTES } from '@/config';

import request from '@/services/request';

import groupModule from './group';
import tabModule from './tab';
import widgetModule from './widget';

export default {
  namespaced: true,
  modules: {
    group: groupModule,
    tab: tabModule,
    widget: widgetModule,
  },
  actions: {
    create(context, { data } = {}) {
      return request.post(API_ROUTES.view.list, data);
    },

    clone(context, { data, id } = {}) {
      return request.post(`${API_ROUTES.view.list}/${id}/clone`, data);
    },

    update(context, { id, data } = {}) {
      return request.put(`${API_ROUTES.view.list}/${id}`, data);
    },

    updateWithoutStore(context, { id, data } = {}) {
      return request.put(`${API_ROUTES.view.list}/${id}`, data);
    },

    updatePositions(context, { data } = {}) {
      return request.put(API_ROUTES.view.positions, data);
    },

    remove(context, { id } = {}) {
      return request.delete(`${API_ROUTES.view.list}/${id}`);
    },

    copy(context, { id, data } = {}) {
      return request.post(`${API_ROUTES.view.copy}/${id}`, data);
    },

    exportWithoutStore(context, { data } = {}) {
      return request.post(API_ROUTES.view.export, data);
    },

    importWithoutStore(context, { data } = {}) {
      return request.post(API_ROUTES.view.import, data);
    },
  },
};
