import { schema, normalize, denormalize } from 'normalizr';
import { merge } from 'lodash';

import { API_ROUTES } from '@/config';

import request, { useRequestCancelling } from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

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
      state.data = merge({}, state.data, entities);
    },
  },
  actions: {
    fetchView({ commit }, { id }) {
      return useRequestCancelling(async (source) => {
        const view = await request.get(`${API_ROUTES.view.list}/${id}`, { cancelToken: source.token });

        const { entities } = normalize(view, viewSchema);

        commit(types.MERGE_VIEW_ENTITIES, { entities });
      }, `view_${id}`);
    },

    async fetchViewTab({ commit }, { id, params }) {
      const viewTab = await request.get(`${API_ROUTES.view.tabs}/${id}`, { params });

      const { entities } = normalize(viewTab, viewTabSchema);

      commit(types.MERGE_VIEW_ENTITIES, { entities });

      return viewTab;
    },
  },
});
