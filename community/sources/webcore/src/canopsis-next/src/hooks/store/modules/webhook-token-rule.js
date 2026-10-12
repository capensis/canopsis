import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Hook for accessing webhook token rule store module helpers
 *
 * @returns {Object} Store module hooks for external auth token management
 */
export const useWebhookTokenRuleStoreModule = () => useStoreModuleHooks('webhookTokenRule');

/**
 * Hook for webhook token rule management operations
 * Provides actions for CRUD operations on webhook token rules
 *
 * @returns {Object} Object containing webhook token rule actions:
 * - fetchWebhookTokenRulesListWithoutStore: Fetch tokens list without storing in state
 * - createWebhookTokenRule: Create new external auth token
 * - updateWebhookTokenRule: Update existing external auth token
 * - removeWebhookTokenRule: Remove external auth token
 */
export const useWebhookTokenRule = () => {
  const { useActions } = useWebhookTokenRuleStoreModule();

  const actions = useActions({
    fetchWebhookTokenRulesListWithoutStore: 'fetchListWithoutStore',
    createWebhookTokenRule: 'create',
    updateWebhookTokenRule: 'update',
    removeWebhookTokenRule: 'remove',
  });

  return actions;
};
