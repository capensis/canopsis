import { storeToRefs } from 'pinia';

import { useWebhookTokenJobStore } from '@/stores/modules/webhook-token-job';

/**
 * Webhook token job store facade.
 *
 * @returns {Object}
 */
export const useWebhookTokenJob = () => {
  const webhookTokenJobStore = useWebhookTokenJobStore();
  const { items, meta, pending, fetchingParams } = storeToRefs(webhookTokenJobStore);

  return {
    webhookTokenJobs: items,
    webhookTokenJobsMeta: meta,
    webhookTokenJobsPending: pending,
    webhookTokenJobsFetchingParams: fetchingParams,
    fetchWebhookTokenJobsList: webhookTokenJobStore.fetchList,
    fetchWebhookTokenJobsListWithPreviousParams: webhookTokenJobStore.fetchListWithPreviousParams,
    stopWebhookTokenJobs: webhookTokenJobStore.stop,
  };
};
