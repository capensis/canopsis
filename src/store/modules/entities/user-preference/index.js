// SERVICES
import request from '@/services/request';
// OTHERS
import { API_ROUTES } from '@/config';
import { userPreferenceSchema } from '@/store/schemas';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
  SET_ACTIVE_FILTER: 'SET_ACTIVE_FILTER',
};

export default {
  namespaced: true,
  state: {
    activeFilter: null,
    pending: false,
    allIds: [],
  },
  getters: {
    filters: (state, getters, rootState, rootGetters) => {
      const userPreferences = rootGetters['entities/getList']('userPreference', state.allIds);
      let filters = [];

      userPreferences.forEach((userPreferenceObject) => {
        filters = filters.concat(userPreferenceObject.widget_preferences.user_filters);
      });

      return filters;
    },
    activeFilter: state => state.activeFilter,
  },
  mutations: {
    [types.SET_ACTIVE_FILTER]: (state, filter) => state.activeFilter = filter,
    [types.FETCH_LIST]: state => state.pending = true,
    [types.FETCH_LIST_COMPLETED]: (state, { allIds }) => {
      state.allIds = allIds;
      state.pending = false;
    },
    [types.FETCH_LIST_FAILED]: (state) => {
      state.pending = false;
    },
  },
  actions: {
    async fetchList({ commit, dispatch }, { params }) {
      try {
        commit(types.FETCH_LIST);

        const { normalizedData } = await dispatch('entities/fetch', {
          route: API_ROUTES.userPreferences,
          schema: [userPreferenceSchema],
          params,
          dataPreparer: d => d,
        }, { root: true });

        commit(types.FETCH_LIST_COMPLETED, {
          allIds: normalizedData.result,
        });
      } catch (e) {
        commit(types.FETCH_LIST_FAILED);
        console.warn(e);
      }
    },
    async setActiveFilter({ commit, getters }, { data, selectedFilter }) {
      try {
        await request.post(API_ROUTES.userPreferences, {
          widget_preferences: {
            user_filters: getters.filters.map(filter => ({
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
        console.warn(e);
      }
    },
  },
};
