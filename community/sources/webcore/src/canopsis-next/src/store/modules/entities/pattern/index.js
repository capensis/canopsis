import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import request from '@/services/request';

import { createEntityModule } from '@/store/plugins/entities';

import corporatePatternModule from './corporate';

export default createEntityModule({
  route: API_ROUTES.pattern.list,
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
      return request.get(API_ROUTES.pattern.list, { params });
    },

    bulkRemove(context, { data }) {
      return request.delete(API_ROUTES.pattern.bulkList, { data });
    },

    checkPatternsCount(context, { data }) {
      return request.post(API_ROUTES.pattern.count, data);
    },

    checkPatternsEntitiesCount(context, { data }) {
      return request.post(API_ROUTES.pattern.entitiesCount, data);
    },

    checkPatternsAlarmsCount(context, { data }) {
      return request.post(API_ROUTES.pattern.alarmsCount, data);
    },
  },
});
