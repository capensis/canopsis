import { storeToRefs } from 'pinia';

import { useWebhookJobStore } from '@/stores/modules/webhook-job';

/**
 * Webhook job store facade.
 *
 * @returns {Object}
 */
export const useWebhookJob = () => {
  const webhookJobStore = useWebhookJobStore();
  const { items, meta, pending, fetchingParams } = storeToRefs(webhookJobStore);

  return {
    webhookJobs: items,
    webhookJobsMeta: meta,
    webhookJobsPending: pending,
    webhookJobsFetchingParams: fetchingParams,
    fetchWebhookJobsList: webhookJobStore.fetchList,
    fetchWebhookJobsListWithPreviousParams: webhookJobStore.fetchListWithPreviousParams,
    stopWebhookJobs: webhookJobStore.stop,
  };
};
