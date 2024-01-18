import { keyBy } from 'lodash';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

import { convertObjectToFormData } from '@/helpers/request';

const types = {
  SET_REGISTERED_ICONS_BY_ID: 'SET_REGISTERED_ICONS_BY_ID',
  ADD_REGISTERED_ICON: 'ADD_REGISTERED_ICON',
  REMOVE_REGISTERED_ICON: 'REMOVE_REGISTERED_ICON',
};

export default createCRUDModule({
  route: API_ROUTES.icons,
  withFetchingParams: true,
  withWithoutStore: true,
}, {
  state: {
    registeredIconsById: {},
  },
  getters: {
    registeredIconsById: state => state.registeredIconsById,
  },
  mutations: {
    [types.ADD_REGISTERED_ICON](state, { id, icon }) {
      state.registeredIconsById[id] = icon;
    },

    [types.REMOVE_REGISTERED_ICON](state, { id }) {
      delete state.registeredIconsById[id];
    },

    [types.SET_REGISTERED_ICONS_BY_ID](state, { icons }) {
      state.registeredIconsById = keyBy(icons, '_id');
    },
  },
  actions: {
    addRegisteredIcon({ commit }, { id, icon } = {}) {
      commit(types.ADD_REGISTERED_ICON, { id, icon });
    },

    removeRegisteredIcon({ commit }, { id } = {}) {
      commit(types.REMOVE_REGISTERED_ICON, { id });
    },

    setRegisteredIcons({ commit }, { icons = [] } = {}) {
      commit(types.SET_REGISTERED_ICONS_BY_ID, { icons });
    },

    create(context, { data } = {}) {
      return request.post(API_ROUTES.icons, convertObjectToFormData(data), {
        headers: { 'Content-Type': 'multipart/form-data' },
      });
    },

    update(context, { id, data } = {}) {
      return request.put(`${API_ROUTES.icons}/${id}`, convertObjectToFormData(data), {
        headers: { 'Content-Type': 'multipart/form-data' },
      });
    },

    fetchItemWithoutStore(context, { id } = {}) {
      return request.get(`${API_ROUTES.icons}/${id}`, {
        headers: { 'Content-Type': 'multipart/form-data' },
      });
    },
  },
});
