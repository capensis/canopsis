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

    fetchLlmHistoryWithoutStore(context, { id, params } = {}) {
      return request.get(`${API_ROUTES.llms.list}/${id}/history`, { params });
    },

    fetchLlmChatsWithoutStore(context, { id, params } = {}) {
      return request.get(`${API_ROUTES.llms.list}/${id}/chats`, { params });
    },

    fetchLlmUsersWithoutStore(context, { id, params } = {}) {
      return request.get(`${API_ROUTES.llms.list}/${id}/users`, { params });
    },

    fetchLlmMessagesWithoutStore(context, { id, params } = {}) {
      return request.get(`${API_ROUTES.llms.list}/${id}/messages`, { params });
    },

    bulkLinkLlmHistory(context, { data } = {}) {
      return request.put(API_ROUTES.llms.bulkHistoryLink, data);
    },
  },
});
