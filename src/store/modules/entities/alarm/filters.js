import request from '@/services/request';

export default {
  namespaced: true,
  state: {
    filters: [],
    activeFilter: null,
  },
  getters: {
    filters(state) {
      return state.filters;
    },
    activeFilter(state) {
      return state.activeFilter;
    },
  },
  mutations: {
    setFilters(state, filters) {
      state.filters = filters;
    },
    setActiveFilter(state, filter) {
      state.activeFilter = filter;
    },
  },
  actions: {
    async loadFilters({ commit }) {
      const response = await request.get('http://localhost:28082/rest/userpreferences/userpreferences', {
        params: {
          limit: 1,
          filter: {
            crecord_name: 'root',
            widget_id: 'widget_listalarm_14e642d2-28d5-f2ba-99f7-2a7fcd62d6f6',
            _id: 'widget_listalarm_14e642d2-28d5-f2ba-99f7-2a7fcd62d6f6_root',
          },
        },
      });
      if (response[0].widget_preferences.user_filters) {
        commit('setFilters', response[0].widget_preferences.user_filters);
      }
    },
    async setActiveFilter({ commit, state }, selectedFilter) {
      try {
        await request.post('http://localhost:28082/rest/userpreferences/userpreferences', {
          widget_preferences: {
            user_filters: state.filters.map(filter => ({
              filter: filter.filter,
              isActive: filter.title === selectedFilter.title,
              title: filter.title,
            })),
          },
          crecord_name: 'root',
          widget_id: 'widget_listalarm_1a6df694-e985-66b7-82c7-6c3012915a88',
          widgetXtype: 'listalarm',
          title: null,
          id: 'bc2a19a5-8d79-ea2f-8172-e340017fbe9f_root',
          _id: 'widget_listalarm_1a6df694-e985-66b7-82c7-6c3012915a88_root',
          crecord_type: 'userpreferences',
        });
        commit('setActiveFilter', selectedFilter);
      } catch (e) {
        console.log(e);
      }
    },
  },
};
