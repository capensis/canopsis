import { API_ROUTES } from '@/config'

import { createCRUDStore } from '@/stores/plugins/crud'

export const useCommentTemplateStore = createCRUDStore({
  id: 'commentTemplate',
  route: API_ROUTES.commentTemplate,
  withFetchingParams: true,
  withWithoutStore: true,
})
