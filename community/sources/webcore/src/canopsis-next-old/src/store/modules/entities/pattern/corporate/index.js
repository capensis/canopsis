import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

export const types = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
};

export default createCRUDModule({
  route: API_ROUTES.pattern.list,
  withFetchingParams: true,
}, {
  actions: {
    async fetchList({ commit }, { params } = {}) {
      try {
        const preparedParams = {
          ...params,

          corporate: true,
        };

        commit(types.FETCH_LIST, { params: preparedParams });

        const { data, meta } = await request.get(API_ROUTES.pattern.list, { params: preparedParams });

        commit(types.FETCH_LIST_COMPLETED, { data, meta });
      } catch (err) {
        console.error(err);
        commit(types.FETCH_LIST_FAILED);

        throw err;
      }
    },
  },
});
