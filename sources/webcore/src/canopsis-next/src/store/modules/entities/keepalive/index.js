import request from '@/services/request';
import { API_ROUTES } from '@/config';


export default {
  namespaced: true,
  actions: {
    async keepAlive(context, payload) {
      try {
        await request.post(API_ROUTES.keepalive, payload);
      } catch (err) {
        console.warn(err);
      }
    },
    async sessionHide(context, payload) {
      try {
        await request.post(API_ROUTES.sessionhide, payload);
      } catch (err) {
        console.warn(err);
      }
    },
  },
};
