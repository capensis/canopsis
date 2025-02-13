import { find } from 'lodash';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

export default createCRUDModule({
  route: API_ROUTES.roles.list,
  withWithoutStore: true,
}, {
  getters: {
    getItemById: state => id => find(state.items, { _id: id }),
  },
  actions: {
    fetchTemplatesListWithoutStore(context, { params }) {
      return request.get(API_ROUTES.roles.templates, { params });
    },

    bulkUpdatePermissions(context, { data }) {
      return request.put(API_ROUTES.roles.bulkPermissions, data);
    },
  },
});
