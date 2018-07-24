import omit from 'lodash/omit';

import request from '@/services/request';
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
  },
  mutations: {
    [types.SET_ACTIVE_FILTER]: (state, filter) => state.activeFilter = filter,
    [types.FETCH_LIST]: state => state.pending = true,
    [types.FETCH_LIST_COMPLETED]: (state) => {
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

        await dispatch('entities/fetch', {
          route: API_ROUTES.userPreferences,
          schema: [userPreferenceSchema],
          params,
          dataPreparer: d => d.data,
        }, { root: true });

        commit(types.FETCH_LIST_COMPLETED);
      } catch (e) {
        commit(types.FETCH_LIST_FAILED);
        console.warn(e);
      }
    },
    async save(context, { userPreference }) {
      try {
        const data = omit(userPreference, ['crecord_creation_time', 'crecord_write_time', 'enable']);
        const response = await request.post(API_ROUTES.userPreferences, JSON.stringify(data));

        console.warn(response);
      } catch (err) {
        console.warn(err);
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
