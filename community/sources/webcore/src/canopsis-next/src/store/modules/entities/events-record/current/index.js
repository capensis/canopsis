import Vue from 'vue';

import { API_ROUTES } from '@/config';

import request, { useRequestCancelling } from '@/services/request';

const types = {
  FETCH_CURRENT: 'FETCH_CURRENT',
  FETCH_CURRENT_COMPLETED: 'FETCH_CURRENT_COMPLETED',
  FETCH_CURRENT_FAILED: 'FETCH_CURRENT_FAILED',

  SET_RESENDING: 'SET_RESENDING',

  RESET: 'RESET',
};

export default {
  namespaced: true,
  state: {
    pending: false,
    current: {},
  },
  getters: {
    pending: state => state.pending,
    current: state => state.current,
  },
  mutations: {
    [types.FETCH_CURRENT]: (state) => {
      state.pending = true;
    },
    [types.FETCH_CURRENT_COMPLETED]: (state, current) => {
      state.current = current;
      state.pending = false;
    },
    [types.FETCH_CURRENT_FAILED]: (state) => {
      state.pending = false;
    },
    [types.SET_RESENDING]: (state, isResending = false) => {
      Vue.set(state.current, 'isResending', isResending);
    },
    [types.RESET]: (state) => {
      state.current = {
        count: state.current.count ?? 0,
        t: null,
        isResending: false,
        isRecording: false,
        _id: '',
      };
    },
  },
  actions: {
    reset({ commit }) {
      commit(types.RESET);
    },
    setCurrentResending({ commit }, isResending) {
      commit(types.SET_RESENDING, isResending);
    },

    setCurrent({ commit }, current) {
      commit(types.FETCH_CURRENT_COMPLETED, current);
    },

    async fetchCurrent({ commit }) {
      return useRequestCancelling(async (source) => {
        try {
          const current = await request.get(API_ROUTES.eventsRecord.current, { cancelToken: source.token });

          commit(types.FETCH_CURRENT_COMPLETED, current);
        } catch (err) {
          console.warn(err);

          commit(types.FETCH_CURRENT_FAILED);
        }
      }, 'event-records-current');
    },

    start(context, { data } = {}) {
      return request.post(API_ROUTES.eventsRecord.current, data);
    },

    async stop({ commit }) {
      const response = await request.delete(API_ROUTES.eventsRecord.current);

      commit(types.RESET);

      return response;
    },
  },
};
