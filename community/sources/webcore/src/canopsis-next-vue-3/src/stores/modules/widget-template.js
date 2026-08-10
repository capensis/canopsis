import { API_ROUTES } from '@/config';
import { createCRUDStore } from '@/stores/plugins/crud';

export const useWidgetTemplateStore = createCRUDStore({
  id: 'widgetTemplate',
  route: API_ROUTES.widgetTemplate,
  withFetchingParams: true,
  withWithoutStore: true,
});
