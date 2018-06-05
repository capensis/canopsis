export const types = {
  ADD_LIVE_REPORTING_FILTER: 'ADD_LIVE_REPORTING_FILTER',
  REMOVE_LIVE_REPORTING_FILTER: 'DELETE_LIVE_REPORTING_FILTER',
};

export default {
  namespaced: true,
  state: {
    liveReportingFilter: '',
  },
  getters: {
    liveReportingFilter: state => state.liveReportingFilter,
  },
  mutations: {
    [types.ADD_LIVE_REPORTING_FILTER](state, filter) {
      state.liveReportingFilter = filter;
    },
    [types.REMOVE_LIVE_REPORTING_FILTER](state) {
      state.liveReportingFilter = '';
    },
  },
  actions: {
    addLiveReportingFilter({ commit }, filter) {
      commit(types.ADD_LIVE_REPORTING_FILTER, filter);
    },
    removeLiveReportingFilter({ commit }) {
      commit(types.REMOVE_LIVE_REPORTING_FILTER);
    },
  },
};
