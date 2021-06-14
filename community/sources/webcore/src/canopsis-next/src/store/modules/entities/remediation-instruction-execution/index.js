import Vue from 'vue';
import { API_ROUTES } from '@/config';

import request from '@/services/request';

export const types = {
  FETCH_ITEM: 'FETCH_ITEM',
  FETCH_ITEM_COMPLETED: 'FETCH_ITEM_COMPLETED',
  FETCH_ITEM_FAILED: 'FETCH_ITEM_FAILED',

  CREATE_ITEM_COMPLETED: 'CREATE_ITEM_COMPLETED',

  UPDATE_ITEM_COMPLETED: 'UPDATE_ITEM_COMPLETED',

  UPDATE_OPERATION_COMPLETED: 'UPDATE_OPERATION_COMPLETED',
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
    [types.FETCH_ITEM_COMPLETED]: (state, instructionExecution) => {
      Vue.set(state.byId, instructionExecution._id, instructionExecution);
      state.allIds.push(instructionExecution._id);

      state.pending = false;
    },
    [types.CREATE_ITEM_COMPLETED]: (state, instructionExecution) => {
      Vue.set(state.byId, instructionExecution._id, instructionExecution);
      state.allIds.push(instructionExecution._id);
    },
    [types.UPDATE_ITEM_COMPLETED]: (state, instructionExecution) => {
      Vue.set(state.byId, instructionExecution._id, instructionExecution);
    },
    [types.UPDATE_OPERATION_COMPLETED]: (state, { id, operation }) => {
      const execution = state.byId[id];

      execution.steps.forEach((step) => {
        const operationIndex = step.operations
          .findIndex(({ operation_id: operationId }) => operationId === operation.operation_id);

        if (operationIndex !== -1) {
          Vue.set(step.operations, operationIndex, operation);
        }
      });
    },
    [types.FETCH_ITEM_FAILED]: (state) => {
      state.pending = false;
    },
  },
  actions: {
    async fetchItem({ commit }, { id, params }) {
      try {
        commit(types.FETCH_ITEM);

        const instructionExecution = await request.get(`${API_ROUTES.remediation.executions}/${id}`, {
          params,
        });

        commit(types.FETCH_ITEM_COMPLETED, instructionExecution);
      } catch (err) {
        console.error(err);

        commit(types.FETCH_ITEM_FAILED);
      }
    },

    async create({ commit }, { data } = {}) {
      const instructionExecution = await request.post(API_ROUTES.remediation.executions, data);

      commit(types.CREATE_ITEM_COMPLETED, instructionExecution);

      return instructionExecution;
    },

    async update({ commit }, { path, id, data }) {
      const instructionExecution = await request.put(`${API_ROUTES.remediation.executions}/${id}/${path}`, data);

      commit(types.UPDATE_ITEM_COMPLETED, instructionExecution);

      return instructionExecution;
    },

    cancel({ dispatch }, { id }) {
      return dispatch('update', { path: 'cancel', id });
    },

    nextOperation({ dispatch }, { id }) {
      return dispatch('update', { path: 'next', id });
    },

    nextStep({ dispatch }, { id, data }) {
      return dispatch('update', { path: 'next-step', id, data });
    },

    pause({ dispatch }, { id }) {
      return dispatch('update', { path: 'pause', id });
    },

    previousOperation({ dispatch }, { id }) {
      return dispatch('update', { path: 'previous', id });
    },

    resume({ dispatch }, { id }) {
      return dispatch('update', { path: 'resume', id });
    },

    rate({ dispatch }, { id, data }) {
      return dispatch('update', { path: 'rate', id, data });
    },

    async ping({ commit }, { id }) {
      try {
        const operation = await request.put(`${API_ROUTES.remediation.executions}/${id}/ping`);

        commit(types.UPDATE_OPERATION_COMPLETED, { id, operation });
      } catch (err) {
        console.warn(err);
      }
    },

    fetchPausedExecutionsWithoutStore(context, { params }) {
      return request.get(API_ROUTES.remediation.pausedExecutions, { params });
    },
  },
};
