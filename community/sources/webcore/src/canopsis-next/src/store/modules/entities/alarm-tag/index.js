import { API_ROUTES } from '@/config';

import { createCRUDModule } from '@/store/plugins/entities';

import labelModule from './label';

export default createCRUDModule({
  route: API_ROUTES.alarmTag.list,
  bulkRoute: API_ROUTES.alarmTag.bulkList,
  withWithoutStore: true,
}, {
  modules: {
    label: labelModule,
  },
});
