import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import request from '@/services/request';

import { createEntityModule } from '@/store/plugins/entities';

import corporatePatternModule from './corporate';

export default createEntityModule({
  route: API_ROUTES.patterns,
  entityType: ENTITIES_TYPES.pattern,
  dataPreparer: d => d.data,
  withFetchingParams: true,
  withMeta: true,
}, {
  modules: {
    corporate: corporatePatternModule,
  },
  actions: {
    fetchListWithoutStore(context, { params } = {}) {
      return request.get(API_ROUTES.patterns, { params });
    },

    bulkRemove(context, { data }) {
      return request.delete(API_ROUTES.bulkPatterns, { data });
    },

    checkPatternsCount(context, { data }) {
      return request.post(API_ROUTES.patternsCount, data);
    },
  },
});
