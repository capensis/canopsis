import { API_ROUTES } from '@/config';

import { createCRUDModule } from '@/store/plugins/entities';

export default createCRUDModule({
  route: API_ROUTES.entityInfosProperties,
  withWithoutStore: true,
}, {
  getters: {
    itemsWithAlias: state => state.items.filter(item => !!item.alias),
    itemsWithoutAlias: state => state.items.filter(item => !item.alias),
  },
});
