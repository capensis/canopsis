import axios from 'axios';

export default {
  namespaced: true,
  state: {
    filters: [],
  },
  getters: {
    filters(state) {
      return state.filters;
    },
  },
  mutations: {
    setFilters(state, filters) {
      state.filters = filters;
    },
  },
  actions: {
    async loadFilters({ commit }) {
      return new Promise((resolve) => {
        axios.get('http://localhost:28082/rest/userpreferences/userpreferences', {
          params: {
            limit: 1,
            filter: {
              crecord_name: 'root',
              widget_id: 'widget_listalarm_14e642d2-28d5-f2ba-99f7-2a7fcd62d6f6',
              _id: 'widget_listalarm_14e642d2-28d5-f2ba-99f7-2a7fcd62d6f6_root',
            },
          },
        })
          .then(({ data }) => {
            if (data.data[0].widget_preferences.user_filters) {
              commit('setFilters', data.data[0].widget_preferences.user_filters);
            }
            resolve();
          });
      });
    },
  },
};
