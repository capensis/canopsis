import { normalize } from 'normalizr';

import request from '@/services/request';
import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';
import schemas from '@/store/schemas';

import { types as entitiesTypes } from '@/store/plugins/entities';

export default {
  namespaced: true,
  actions: {
    async create({ commit }, { data, parents, parentsType }) {
      try {
        const parentSchema = schemas[parentsType];
        const id = await request.post(API_ROUTES.pbehavior, data);
        const pbehavior = {
          ...data,
          enabled: true,
          _id: id,
        };

        const parentEntities = parents.map(parent => ({ ...parent, pbehaviors: [...parent.pbehaviors, pbehavior] }));

        const { entities } = normalize(parentEntities, [parentSchema]);

        commit(entitiesTypes.ENTITIES_MERGE, entities, { root: true });
      } catch (err) {
        console.error(err);

        throw err;
      }
    },
    async remove({ dispatch }, { id }) {
      try {
        await request.delete(`${API_ROUTES.pbehavior}/${id}`);
        await dispatch('entities/removeFromStore', {
          id,
          type: ENTITIES_TYPES.pbehavior,
        }, { root: true });
      } catch (err) {
        console.warn(err);
      }
    },
  },
};
