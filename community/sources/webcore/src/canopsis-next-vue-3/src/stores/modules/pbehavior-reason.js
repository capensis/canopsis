import { API_ROUTES } from '@/config'

import request from '@/services/request'

import { createCRUDStore } from '@/stores/plugins/crud'

export const usePbehaviorReasonStore = createCRUDStore({
  id: 'pbehaviorReason',
  route: API_ROUTES.pbehavior.reasons,
  withFetchingParams: true,
  withWithoutStore: true,
  actions: {
    bulkHide({ data } = {}) {
      return request.put(API_ROUTES.pbehavior.bulkReasonsHide, data)
    },

    bulkUnhide({ data } = {}) {
      return request.put(API_ROUTES.pbehavior.bulkReasonsUnhide, data)
    },

    bulkRemove({ data } = {}) {
      return request.delete(API_ROUTES.pbehavior.bulkReasons, { data })
    },
  },
})
