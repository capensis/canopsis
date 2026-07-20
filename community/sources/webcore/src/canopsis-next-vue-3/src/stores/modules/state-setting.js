import { API_ROUTES } from '@/config'

import { createCRUDStore } from '@/stores/plugins/crud'

export const useStateSettingStore = createCRUDStore({
  id: 'stateSetting',
  route: API_ROUTES.stateSetting,
  withFetchingParams: true,
})
