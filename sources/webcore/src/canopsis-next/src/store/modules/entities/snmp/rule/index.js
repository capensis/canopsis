import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import request from '@/services/request';

import { snmpRuleSchema } from '@/store/schemas';
import { createEntityModule } from '@/store/plugins/entities';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
};

export default createEntityModule({
  types,
  route: API_ROUTES.snmpRule.list,
  entityType: ENTITIES_TYPES.snmpRule,
}, {
  actions: {
    async fetchList({ commit, dispatch }, { params } = {}) {
      try {
        commit(types.FETCH_LIST, { params });

        const { normalizedData } = await dispatch('entities/fetch', {
          body: params,
          method: 'POST',
          route: API_ROUTES.snmpRule.list,
          dataPreparer: d => d.data,
          schema: [snmpRuleSchema],
        }, { root: true });

        commit(types.FETCH_LIST_COMPLETED, {
          allIds: normalizedData.result,
        });
      } catch (err) {
        console.error(err);
        commit(types.FETCH_LIST_FAILED);

        throw err;
      }
    },
    create(context, { data } = {}) {
      return request.post(API_ROUTES.snmpRule.create, data);
    },
    remove(context, { data = {} } = {}) {
      return request.delete(API_ROUTES.snmpRule.list, { data });
    },
  },
});
