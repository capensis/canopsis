import { API_ROUTES } from '@/config'

import { createCRUDStore } from '@/stores/plugins/crud'

export const useTagStore = createCRUDStore({
  id: 'tag',
  route: API_ROUTES.tags,
  withFetchingParams: true,
  withWithoutStore: true,
})
