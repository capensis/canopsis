import { useWebhookTokenRuleStore } from '@/stores/modules/webhook-token-rule';

/**
 * Webhook token rule store hook.
 *
 * @returns {Object}
 */
export const useWebhookTokenRule = () => {
  const webhookTokenRuleStore = useWebhookTokenRuleStore();

  return {
    fetchWebhookTokenRulesListWithoutStore: webhookTokenRuleStore.fetchListWithoutStore,
    createWebhookTokenRule: webhookTokenRuleStore.create,
    updateWebhookTokenRule: webhookTokenRuleStore.update,
    removeWebhookTokenRule: webhookTokenRuleStore.remove,
  };
};
