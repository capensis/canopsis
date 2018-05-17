import request from '@/services/request';
import { API_ROUTES } from '@/config';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  SET_ACTIVE_FILTER: 'SET_ACTIVE_FILTER',
};

export default {
  namespaced: true,
  state: {
    filters: [],
    activeFilter: null,
    pending: false,
  },
  getters: {
    filters: state => state.filters,
    activeFilter: state => state.activeFilter,
  },
  mutations: {
    [types.SET_ACTIVE_FILTER]: (state, filter) => state.activeFilter = filter,
    [types.FETCH_LIST]: state => state.pending = true,
    [types.FETCH_LIST_COMPLETED]: (state, filters) => {
      state.filters = filters;
      state.pending = false;
    },
  },
  actions: {
    async fetchList({ commit }, { params }) {
      try {
        commit(types.FETCH_LIST);
        const response = await request.get(API_ROUTES.userPreferences, {
          params: {
            limit: params.limit,
            filter: {
              crecord_name: params.filter.crecord_name,
              widget_id: params.filter.widget_id,
              _id: params.filter._id,
            },
          },
        });

        const data = response.shift();
        let filters = [];

        if (data) {
          filters = data.widget_preferences.user_filters;
        }

        commit(types.FETCH_LIST_COMPLETED, filters);
      } catch (e) {
        commit(types.FETCH_LIST_COMPLETED, []);
        console.log(e);
      }
    },
    async setActiveFilter({ commit, state }, { data, selectedFilter }) {
      try {
        await request.post(API_ROUTES.userPreferences, {
          widget_preferences: {
            user_filters: state.filters.map(filter => ({
              filter: filter.filter,
              title: filter.title,
            })),
            selected_filter: {
              filter: selectedFilter.filter,
              isActive: true,
              title: selectedFilter.title,
            },
          },
          crecord_name: data.crecord_name,
          widget_id: data.widget_id,
          widgetXtype: data.widgetXtype,
          title: data.title,
          id: data.id,
          _id: data._id,
          crecord_type: data.crecord_type,
        });
        commit(types.SET_ACTIVE_FILTER, selectedFilter);
      } catch (e) {
        console.log(e);
      }
    },
  },
};
