import { API_ROUTES } from '@/config';
import { createCRUDStore } from '@/stores/plugins/crud';

export const useAlarmTagLabelStore = createCRUDStore({
  id: 'alarmTagLabel',
  route: API_ROUTES.alarmTag.label,
  withWithoutStore: true,
});
