import { API_ROUTES } from '@/config';

import { createCRUDStore } from '@/stores/plugins/crud';

export const useWebhookTokenRuleStore = createCRUDStore({
  id: 'webhookTokenRule',
  route: API_ROUTES.webhookTokenRule,
  withWithoutStore: true,
});
