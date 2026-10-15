import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

import corporatePatternModule from './corporate';
import entitiesOptimizePatternModule from './entities-optimize';

export default createCRUDModule({
  route: API_ROUTES.pattern.list,
  bulkRoute: API_ROUTES.pattern.bulkList,
  withFetchingParams: true,
  withWithoutStore: true,
}, {
  modules: {
    corporate: corporatePatternModule,
    entitiesOptimize: entitiesOptimizePatternModule,
  },
  actions: {
    checkPatternsEntitiesCount(context, { data }) {
      return request.post(API_ROUTES.pattern.entitiesCount, data);
    },

    checkPatternsAlarmsCount(context, { data }) {
      return request.post(API_ROUTES.pattern.alarmsCount, data);
    },
  },
});
