import { API_ROUTES } from '@/config';

import request from '@/services/request';

import { createCRUDModule } from '@/store/plugins/entities';

export default createCRUDModule({
  route: API_ROUTES.llms.list,
  bulkRoute: API_ROUTES.llms.bulk,
  withWithoutStore: true,
}, {
  actions: {
    fetchModelsListWithoutStore(context, { params } = {}) {
      return request.get(API_ROUTES.llms.models, { params });
    },

    fetchPromptsHistoryWithoutStore(context, { id, params } = {}) {
      return request.get(`${API_ROUTES.llms.list}/${id}/prompts-history`, { params });
    },
  },
});
