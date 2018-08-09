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
    async create(context, params = {}) {
      try {
        await request.post(API_ROUTES.viewV3.view, params);
      } catch (err) {
        console.warn(err);
      }
    },
  },
};
