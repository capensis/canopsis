import { normalize } from 'normalizr';

import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import request, { useRequestCancelling } from '@/services/request';

import { viewSchema } from '@/store/schemas';
import { types as entitiesTypes } from '@/store/plugins/entities';

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
  getters: {
    getItemById: (state, getters, rootState, rootGetters) => id => rootGetters['entities/getItem'](
      ENTITIES_TYPES.view,
      id,
    ),
  },
  actions: {
    fetchItem({ dispatch }, { id }) {
      return useRequestCancelling(source => dispatch('entities/fetch', {
        route: `${API_ROUTES.view.list}/${id}`,
        schema: viewSchema,
        cancelToken: source.token,
      }, { root: true }), 'view');
    },

    create(context, { data } = {}) {
      return request.post(API_ROUTES.view.list, data);
    },

    clone(context, { data, id } = {}) {
      return request.post(`${API_ROUTES.view.list}/${id}/clone`, data);
    },

    async update({ commit }, { id, data } = {}) {
      const result = await request.put(`${API_ROUTES.view.list}/${id}`, data);

      const { entities } = normalize(result, viewSchema);

      commit(entitiesTypes.ENTITIES_UPDATE, entities, { root: true });

      return result;
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
