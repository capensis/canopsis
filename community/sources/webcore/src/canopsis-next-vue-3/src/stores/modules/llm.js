import { API_ROUTES } from '@/config'

import request from '@/services/request'

import { createCRUDStore } from '@/stores/plugins/crud'

export const useLlmStore = createCRUDStore({
  id: 'llm',
  route: API_ROUTES.llms.list,
  bulkRoute: API_ROUTES.llms.bulk,
  withFetchingParams: true,
  withWithoutStore: true,
  actions: {
    bulkEnable({ data } = {}) {
      return request.put(API_ROUTES.llms.bulkEnable, data)
    },

    bulkDisable({ data } = {}) {
      return request.put(API_ROUTES.llms.bulkDisable, data)
    },

    fetchModelsListWithoutStore({ params } = {}) {
      return request.get(API_ROUTES.llms.models, { params })
    },

    fetchLlmHistoryWithoutStore({ id, params } = {}) {
      return request.get(`${API_ROUTES.llms.list}/${id}/history`, { params })
    },

    fetchLlmChatsWithoutStore({ id, params } = {}) {
      return request.get(`${API_ROUTES.llms.list}/${id}/chats`, { params })
    },

    fetchLlmUsersWithoutStore({ id, params } = {}) {
      return request.get(`${API_ROUTES.llms.list}/${id}/users`, { params })
    },

    fetchLlmMessagesWithoutStore({ id, params } = {}) {
      return request.get(`${API_ROUTES.llms.list}/${id}/messages`, { params })
    },

    bulkLinkLlmHistory({ data } = {}) {
      return request.put(API_ROUTES.llms.bulkHistoryLink, data)
    },
  },
})
