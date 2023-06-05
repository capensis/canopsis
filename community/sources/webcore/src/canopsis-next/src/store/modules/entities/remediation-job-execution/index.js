import Vue from 'vue';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

export const types = {
  FETCH_OUTPUT: 'FETCH_OUTPUT',
  FETCH_OUTPUT_COMPLETED: 'FETCH_OUTPUT_COMPLETED',
  FETCH_OUTPUT_FAILED: 'FETCH_OUTPUT_FAILED',
};

export default {
  namespaced: true,
  state: {
    outputs: {},
  },
  getters: {
    getOutputById: state => id => state.outputs[id] ?? '',
  },
  mutations: {
    [types.FETCH_OUTPUT_COMPLETED]: (state, { id, output }) => {
      Vue.set(state.outputs, id, output);
    },
  },
  actions: {
    create(context, { data } = {}) {
      return request.post(API_ROUTES.remediation.jobExecutions, data);
    },

    cancel(context, { id }) {
      return request.put(`${API_ROUTES.remediation.jobExecutions}/${id}/cancel`);
    },

    async fetchOutput({ commit }, { id }) {
      const output = await request.get(`${API_ROUTES.remediation.jobExecutions}/${id}/output`);

      commit(types.FETCH_OUTPUT_COMPLETED, { id, output });
    },
  },
};
