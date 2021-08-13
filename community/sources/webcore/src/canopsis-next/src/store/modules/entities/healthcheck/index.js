import { API_ROUTES } from '@/config';

import request from '@/services/request';

const types = {
  FETCH_STATUS: 'FETCH_STATUS',
  FETCH_STATUS_COMPLETED: 'FETCH_STATUS_COMPLETED',
  FETCH_STATUS_FAILED: 'FETCH_STATUS_FAILED',
};

export default {
  namespaced: true,
  state: {
    pending: true,
    services: [],
    engines: {
      edges: [],
      nodes: [],
    },
    maxQueueLength: 0,
    hasInvalidEnginesOrder: false,
    error: null,
  },
  getters: {
    pending: state => state.pending,
    services: state => state.services,
    engines: state => ({
      edges: state.engines.edges,
      nodes: state.engines.nodes.map(node => ({
        ...node,
        max_queue_length: state.maxQueueLength,
      })),
    }),
    maxQueueLength: state => state.maxQueueLength,
    hasInvalidEnginesOrder: state => state.hasInvalidEnginesOrder,
    error: state => state.error,
  },
  mutations: {
    [types.FETCH_STATUS](state) {
      state.pending = true;
    },

    [types.FETCH_STATUS_COMPLETED](state, {
      services,
      engines,
      maxQueueLength,
      hasInvalidEnginesOrder,
    }) {
      state.services = services;
      state.engines = engines;
      state.maxQueueLength = maxQueueLength;
      state.hasInvalidEnginesOrder = hasInvalidEnginesOrder;
      state.pending = false;
    },

    [types.FETCH_STATUS_FAILED](state, { error }) {
      state.pending = false;
      state.error = error;
    },
  },
  actions: {
    async fetchStatus({ commit }) {
      try {
        commit(types.FETCH_STATUS);

        const {
          services = [],
          engines = {},
          max_queue_length: maxQueueLength,
          has_invalid_engines_order: hasInvalidEnginesOrder,
        } = await request.get(API_ROUTES.healthcheck);

        commit(types.FETCH_STATUS_COMPLETED, {
          services,
          engines,
          maxQueueLength,
          hasInvalidEnginesOrder,
        });
      } catch (error) {
        commit(types.FETCH_STATUS_COMPLETED, { error });
      }
    },
  },
};
