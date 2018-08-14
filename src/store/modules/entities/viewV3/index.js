import request from '@/services/request';
import { API_ROUTES } from '@/config';
import groupModule from './group';


export const types = {
};

export default {
  namespaced: true,
  modules: {
    group: groupModule,
  },
  state: {
  },
  getters: {
  },
  mutations: {
  },
  actions: {
    async create({ dispatch }, params = {}) {
      try {
        await request.post(API_ROUTES.viewV3.view, params);
      } catch (err) {
        console.warn(err);
        await dispatch('popup/add', { type: 'error', text: err.description }, { root: true });
      }
    },
  },
};
