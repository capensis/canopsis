import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

import corporatePatternModule from './corporate';

export default createCRUDModule({
  route: API_ROUTES.pattern.list,
  withFetchingParams: true,
  withWithoutStore: true,
}, {
  modules: {
    corporate: corporatePatternModule,
  },
  actions: {
    bulkRemove(context, { data }) {
      return request.delete(API_ROUTES.pattern.bulkList, { data });
    },

    checkPatternsEntitiesCount(context, { data }) {
      return request.post(API_ROUTES.pattern.entitiesCount, data);
    },

    checkPatternsAlarmsCount(context, { data }) {
      return request.post(API_ROUTES.pattern.alarmsCount, data);
    },
  },
});
