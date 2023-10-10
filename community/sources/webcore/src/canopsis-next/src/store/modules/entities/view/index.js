import Vue from 'vue';
import { schema, normalize, denormalize } from 'normalizr';

import { API_ROUTES } from '@/config';

import request, { useRequestCancelling } from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

import widgetModule from './widget';

const VIEW_ENTITIES = {
  viewTab: 'viewTab',
  view: 'view',
  group: 'group',
};

export const viewTabSchema = new schema.Entity(VIEW_ENTITIES.viewTab, {}, { idAttribute: '_id' });
export const viewSchema = new schema.Entity(VIEW_ENTITIES.view, {
  tabs: [viewTabSchema],
}, { idAttribute: '_id' });
export const groupSchema = new schema.Entity(VIEW_ENTITIES.group, {
  views: [viewSchema],
}, { idAttribute: '_id' });

const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',

  MERGE_VIEW_ENTITIES: 'MERGE_VIEW_ENTITIES',
};

export default createCRUDModule({
  types,
  route: API_ROUTES.view.groups,
  withWithoutStore: true,
}, {
  modules: {
    widget: widgetModule,
  },
  state: {
    data: {
      [VIEW_ENTITIES.viewTab]: {},
      [VIEW_ENTITIES.view]: {},
      [VIEW_ENTITIES.group]: {},
    },
    ids: [],
  },
  getters: {
    items: state => denormalize(state.ids, [groupSchema], state.data),
    getGroupById: state => id => denormalize(id, groupSchema, state.data),
    getViewById: state => id => denormalize(id, viewSchema, state.data),
    getViewTabById: state => id => denormalize(id, viewTabSchema, state.data),
  },
  mutations: {
    [types.FETCH_LIST_COMPLETED](state, { data, meta }) {
      const { entities, result } = normalize(data, [groupSchema]);

      state.data = entities;
      state.ids = result;
      state.meta = meta;
      state.pending = false;
    },

    [types.MERGE_VIEW_ENTITIES](state, { entities }) {
      Object.entries(entities).forEach(([type, typeEntities]) => {
        Vue.set(state.data, type, { ...state.data[type], ...typeEntities });
      });
    },
  },
  actions: {
    createPrivateGroup(context, { data } = {}) {
      return request.post(API_ROUTES.privateView.groups, data);
    },

    updatePrivateGroup(context, { id, data } = {}) {
      return request.put(`${API_ROUTES.privateView.groups}/${id}`, data);
    },

    removePrivateGroup(context, { id } = {}) {
      return request.delete(`${API_ROUTES.privateView.groups}/${id}`);
    },

    createView(context, { data } = {}) {
      return request.post(API_ROUTES.view.list, data);
    },

    cloneView(context, { data, id } = {}) {
      return request.post(`${API_ROUTES.view.list}/${id}/clone`, data);
    },

    updateView(context, { id, data } = {}) {
      return request.put(`${API_ROUTES.view.list}/${id}`, data);
    },

    updateViewWithoutStore(context, { id, data } = {}) {
      return request.put(`${API_ROUTES.view.list}/${id}`, data);
    },

    updateViewPositions(context, { data } = {}) {
      return request.put(API_ROUTES.view.positions, data);
    },

    removeView(context, { id } = {}) {
      return request.delete(`${API_ROUTES.view.list}/${id}`);
    },

    copyView(context, { id, data } = {}) {
      return request.post(`${API_ROUTES.view.copy}/${id}`, data);
    },

    exportViewWithoutStore(context, { data } = {}) {
      return request.post(API_ROUTES.view.export, data);
    },

    importViewWithoutStore(context, { data } = {}) {
      return request.post(API_ROUTES.view.import, data);
    },

    fetchView({ commit }, { id }) {
      return useRequestCancelling(async (source) => {
        const view = await request.get(`${API_ROUTES.view.list}/${id}`, { cancelToken: source.token });

        const { entities } = normalize(view, viewSchema);

        commit(types.MERGE_VIEW_ENTITIES, { entities });
      }, `view_${id}`);
    },

    createViewTab(context, { data } = {}) {
      return request.post(API_ROUTES.view.tabs, data);
    },

    cloneViewTab(context, { data, id } = {}) {
      return request.put(`${API_ROUTES.view.tabs}/${id}/clone`, data);
    },

    updateViewTab(context, { data, id } = {}) {
      return request.put(`${API_ROUTES.view.tabs}/${id}`, data);
    },

    removeViewTab(context, { id } = {}) {
      return request.delete(`${API_ROUTES.view.tabs}/${id}`);
    },

    copyViewTab(context, { id, data } = {}) {
      return request.post(`${API_ROUTES.view.tabCopy}/${id}`, data);
    },

    updateViewTabPositions(context, { data } = {}) {
      return request.put(API_ROUTES.view.tabPositions, data);
    },

    async fetchViewTab({ commit }, { id, params }) {
      const viewTab = await request.get(`${API_ROUTES.view.tabs}/${id}`, { params });

      const { entities } = normalize(viewTab, viewTabSchema);

      commit(types.MERGE_VIEW_ENTITIES, { entities });

      return viewTab;
    },
  },
});
