import Vue from 'vue';
import { API_ROUTES } from '@/config';

import request from '@/services/request';

export const types = {
  FETCH_ITEM: 'FETCH_ITEM',
  FETCH_ITEM_COMPLETED: 'FETCH_ITEM_COMPLETED',
  FETCH_ITEM_FAILED: 'FETCH_ITEM_FAILED',

  CREATE_ITEM_COMPLETED: 'CREATE_ITEM_COMPLETED',
};

export default {
  namespaced: true,
  state: {
    allIds: [],
    byId: {},
    pending: true,
  },
  getters: {
    pending: state => state.pending,
    items: state => state.allIds.map(id => state.byId[id]),
    getItemById: state => id => state.byId[id],
  },
  mutations: {
    [types.FETCH_ITEM]: (state) => {
      state.pending = true;
    },
    [types.FETCH_ITEM_COMPLETED]: (state, jobExecution) => {
      Vue.set(state.byId, jobExecution._id, jobExecution);
      state.allIds.push(jobExecution._id);

      state.pending = false;
    },
    [types.FETCH_ITEM_FAILED]: (state) => {
      state.pending = false;
    },
    [types.CREATE_ITEM_COMPLETED]: (state, jobExecution) => {
      Vue.set(state.byId, jobExecution._id, jobExecution);
      state.allIds.push(jobExecution._id);
    },
  },
  actions: {
    async fetchItem({ commit }, { id, params }) {
      try {
        commit(types.FETCH_ITEM);

        const jobExecution = await request.get(`${API_ROUTES.remediation.jobExecutions}/${id}`, {
          params,
        });

        commit(types.FETCH_ITEM_COMPLETED, jobExecution);
      } catch (err) {
        console.error(err);

        commit(types.FETCH_ITEM_FAILED);
      }
    },

    async create({ commit }, { data } = {}) {
      const jobExecution = await request.post(API_ROUTES.remediation.jobExecutions, data);

      commit(types.CREATE_ITEM_COMPLETED, jobExecution);

      return jobExecution;
    },

    async cancel(context, { id } = {}) {
      await request.put(`${API_ROUTES.remediation.jobExecutions}/${id}/cancel`);
    },
  },
};
