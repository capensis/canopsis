import Vue from 'vue';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

const types = {
  FETCH_FIELD_LIST: 'FETCH_FIELD_LIST',
  FETCH_FIELD_LIST_COMPLETED: 'FETCH_FIELD_LIST_COMPLETED',
  FETCH_FIELD_LIST_FAILED: 'FETCH_FIELD_LIST_FAILED',
};

export default createCRUDModule({
  route: API_ROUTES.pbehavior.types,
  withFetchingParams: true,
  withWithoutStore: true,
}, {
  state: {
    field: {},
  },
  getters: {
    fieldItems: state => state.field.items ?? [],
    fieldPending: state => state.field.pending ?? false,
  },
  mutations: {
    [types.FETCH_FIELD_LIST]: (state) => {
      Vue.set(state.field, 'pending', true);
    },
    [types.FETCH_FIELD_LIST_COMPLETED]: (state, { items = [] } = {}) => {
      state.field = { items, pending: false };
    },
    [types.FETCH_FIELD_LIST_FAILED]: (state) => {
      Vue.set(state.field, 'pending', false);
    },
  },
  actions: {
    fetchNextPriority() {
      return request.get(API_ROUTES.pbehavior.nextTypesPriority);
    },

    async fetchFieldList({ commit }, { params = { paginate: false, with_hidden: true } } = {}) {
      try {
        commit(types.FETCH_FIELD_LIST);

        const { data: items } = await request.get(API_ROUTES.pbehavior.types, { params });

        commit(types.FETCH_FIELD_LIST_COMPLETED, { items });
      } catch (err) {
        console.error(err);

        commit(types.FETCH_FIELD_LIST_FAILED);
      }
    },
  },
});
