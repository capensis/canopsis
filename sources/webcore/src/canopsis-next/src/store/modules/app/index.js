import { GROUPS_NAVIGATION_TYPES } from '@/constants';
import localStorageDataSource from '@/services/local-storage-data-source';

const GROUPS_NAVIGATION_TYPE_KEY = 'groups-navigation-type';

export const types = {
  SET_GROUPS_NAVIGATION_TYPE: 'SET_GROUPS_NAVIGATION_TYPE',
};

export default {
  namespaced: true,
  state: {
    groupsNavigationType: localStorageDataSource.getItem(GROUPS_NAVIGATION_TYPE_KEY) || GROUPS_NAVIGATION_TYPES.sideBar,
  },
  getters: {
    groupsNavigationType: state => state.groupsNavigationType,
  },
  mutations: {
    [types.SET_GROUPS_NAVIGATION_TYPE](state, type) {
      state.groupsNavigationType = type;
    },
  },
  actions: {
    setGroupsNavigationType({ commit }, type) {
      commit(types.SET_GROUPS_NAVIGATION_TYPE, type);
      localStorageDataSource.setItem(GROUPS_NAVIGATION_TYPE_KEY, type);
    },
  },
};
