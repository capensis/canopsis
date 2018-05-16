import request from '@/services/request';
import { API_ROUTES } from '@/config';

export default {
  namespaced: true,
  actions: {
    async create(context, { data }) {
      try {
        await request.post(API_ROUTES.event, { event: data });
      } catch (e) {
        console.log(e);
      }
    },
  },
};
