import request from '@/services/request';
import { API_ROUTES } from '@/config';

import { types as entitiesTypes } from '@/store/plugins/entities';

export default {
  namespaced: true,
  actions: {
    async create({ commit }, { data, parents }) {
      try {
        const id = await request.post(API_ROUTES.pbehavior, data);
        const pbehaviorEntities = { [id]: { ...data, enabled: true, _id: id } };
        const alarmEntities = parents.reduce((acc, parent) => {
          const parentPbehaviors = parent.pbehaviors.map(v => v._id);

          acc[parent._id] = { pbehaviors: [...parentPbehaviors, id] };

          return acc;
        }, {});

        commit(entitiesTypes.ENTITIES_MERGE, {
          pbehavior: pbehaviorEntities,
          alarm: alarmEntities,
        }, { root: true });
      } catch (err) {
        console.error(err);

        throw err;
      }
    },
    /**
     * TODO: finish it
     */
    async remove({ dispatch }, { id }) {
      try {
        await dispatch('entities/remove', { id, type: 'pbehavior' }, { root: true });
        // await request.delete(`${API_ROUTES.pbehavior}/${id}`);

        // console.log(`${parentId}`id);

        /*        commit(
          entitiesTypes.ENTITIES_DELETE,
          { pbehavior: [id] },
          { root: true },
        ); */
      } catch (err) {
        console.warn(err);
      }
    },
  },
};
