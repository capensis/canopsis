import request from '@/services/request';
import { API_ROUTES } from '@/config';

import { types as entitiesTypes } from '@/store/plugins/entities';

export default {
  namespaced: true,
  actions: {
    async create(context, data) {
      try {
        await request.post(API_ROUTES.pbehavior, data);
      } catch (err) {
        console.error(err);

        throw err;
      }
    },
    async remove({ commit }, { id }) {
      try {
        await request.delete(`${API_ROUTES.pbehavior}/${id}`);

        commit(
          entitiesTypes.ENTITIES_DELETE,
          { pbehavior: [id] },
          { root: true },
        );
      } catch (err) {
        console.warn(err);
      }
    },
  },
};
