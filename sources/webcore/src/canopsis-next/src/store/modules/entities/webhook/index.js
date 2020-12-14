import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import { createEntityModule } from '@/store/plugins/entities';
import { DEFAULT_ENTITY_MODULE_TYPES } from '@/store/plugins/entities/create-entity-module';

import schemas from '@/store/schemas';

export default createEntityModule({
  route: API_ROUTES.webhook,
  entityType: ENTITIES_TYPES.webhook,
  withMeta: true,
  withFetchingParams: true,
  dataPreparer: d => d.data,
}, {
  actions: {
    async fetchList({ commit, dispatch }, { params } = {}) {
      try {
        commit(DEFAULT_ENTITY_MODULE_TYPES.FETCH_LIST, { params });

        const { normalizedData, data } = await dispatch('entities/fetch', {
          route: API_ROUTES.webhook,
          params,
          dataPreparer: d => d.data,
          schema: [schemas.webhook],
        }, { root: true });

        commit(DEFAULT_ENTITY_MODULE_TYPES.FETCH_LIST_COMPLETED, {
          meta: {
            total_count: data.total_count,
            count: data.count,
          },
          allIds: normalizedData.result,
        });
      } catch (err) {
        console.error(err);
        commit(DEFAULT_ENTITY_MODULE_TYPES.FETCH_LIST_FAILED);

        throw err;
      }
    },
  },
});
