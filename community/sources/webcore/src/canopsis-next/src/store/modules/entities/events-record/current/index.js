import { API_ROUTES } from '@/config';

import request, { useRequestCancelling } from '@/services/request';

const DEFAULT_STATUS = {
  limit: 1,
  recording: [],
  resending: [],
};

const types = {
  FETCH_CURRENT: 'FETCH_CURRENT',
  FETCH_CURRENT_COMPLETED: 'FETCH_CURRENT_COMPLETED',
  FETCH_CURRENT_FAILED: 'FETCH_CURRENT_FAILED',

  SET_STATUS: 'SET_STATUS',
};

export default {
  namespaced: true,
  state: {
    pending: false,
    status: { ...DEFAULT_STATUS },
  },
  getters: {
    pending: state => state.pending,
    status: state => state.status,
    recordings: state => state.status?.recording ?? [],
    resendings: state => state.status?.resending ?? [],
    limit: state => state.status?.limit ?? 1,
    recordingsById: (state, getters) => Object.fromEntries(
      getters.recordings.map(item => [item._id, true]),
    ),
    resendingsById: (state, getters) => Object.fromEntries(
      getters.resendings.map(item => [item._id, true]),
    ),
  },
  mutations: {
    [types.FETCH_CURRENT]: (state) => {
      state.pending = true;
    },
    [types.FETCH_CURRENT_COMPLETED]: (state, status) => {
      state.status = status ?? { ...DEFAULT_STATUS };
      state.pending = false;
    },
    [types.FETCH_CURRENT_FAILED]: (state) => {
      state.pending = false;
    },
    [types.SET_STATUS]: (state, status) => {
      state.status = status ?? { ...DEFAULT_STATUS };
    },
  },
  actions: {
    reset({ commit }) {
      commit(types.RESET);
    },

    setStatus({ commit }, status) {
      commit(types.SET_STATUS, status);
    },

    async fetchCurrent({ commit }) {
      return useRequestCancelling(async (source) => {
        try {
          const status = await request.get(API_ROUTES.eventsRecord.current, { cancelToken: source.token });

          commit(types.FETCH_CURRENT_COMPLETED, status);
        } catch (err) {
          console.warn(err);

          commit(types.FETCH_CURRENT_FAILED);
        }
      }, 'event-records-current');
    },

    start(context, { data } = {}) {
      return request.post(API_ROUTES.eventsRecord.current, data);
    },

    stop(context, { id } = {}) {
      return request.delete(`${API_ROUTES.eventsRecord.current}/${id}`);
    },
  },
};
