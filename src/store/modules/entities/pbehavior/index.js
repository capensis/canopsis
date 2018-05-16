import request from '@/services/request';
import { API_ROUTES } from '@/config';

import entitiesTypes from '../types';

export default {
  namespaced: true,
  actions: {
    async create(context, data) {
      try {
        await request.post(API_ROUTES.pbehavior, data);
      } catch (err) {
        console.warn(err);
      }
    },
    async remove({ commit }, { id }) {
      try {
        await request.delete(`${API_ROUTES.pbehavior}/${id}`);

        commit(
          `entities/${entitiesTypes.ENTITIES_DELETE}`,
          { pbehavior: [id] },
          { root: true },
        );
      } catch (err) {
        console.warn(err);
      }
    },
  },
};
