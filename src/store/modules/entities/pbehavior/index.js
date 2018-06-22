// SERVICES
import request from '@/services/request';
// OTHERS
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
