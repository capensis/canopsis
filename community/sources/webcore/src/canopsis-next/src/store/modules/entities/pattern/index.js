import { API_ROUTES } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

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
});
