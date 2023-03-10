import { API_ROUTES } from '@/config';

import { ENTITIES_TYPES } from '@/constants';

import request from '@/services/request';

import { createEntityModule } from '@/store/plugins/entities';

export default createEntityModule({
  route: API_ROUTES.linkRule,
  entityType: ENTITIES_TYPES.linkRule,
  dataPreparer: d => d.data,
  withMeta: true,
}, {
  actions: {
    bulkRemove(context, { data }) {
      return request.delete(API_ROUTES.linkRule, { data });
    },

    fetchLinkCategoriesWithoutStore(context, { params } = {}) {
      return request.get(API_ROUTES.linkCategories, { params });
    },
  },
});
