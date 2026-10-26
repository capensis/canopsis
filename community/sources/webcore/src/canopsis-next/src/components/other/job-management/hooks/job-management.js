import { storeToRefs } from 'pinia';
import { unref } from 'vue';

import { JOBS_TABS } from '@/constants';

import { useInstructionJobStore } from '@/stores/modules/instruction-job';
import { useTicketStatusJobStore } from '@/stores/modules/ticket-status-job';
import { useWebhookJobStore } from '@/stores/modules/webhook-job';
import { useWebhookTokenJobStore } from '@/stores/modules/webhook-token-job';

const JOB_STORES = {
  [JOBS_TABS.instructions]: useInstructionJobStore,
  [JOBS_TABS.webhooks]: useWebhookJobStore,
  [JOBS_TABS.ticketStatus]: useTicketStatusJobStore,
  [JOBS_TABS.authToken]: useWebhookTokenJobStore,
};

/**
 * Maps a job store to the common job management interface.
 *
 * @param {string | import('vue').Ref<string>} jobType - Job type or its reactive reference.
 * @returns {Object}
 */
export const useJobManagement = (jobType) => {
  const store = JOB_STORES[unref(jobType)]();
  const {
    items: jobs,
    meta: jobsMeta,
    pending: jobsPending,
    fetchingParams: jobsFetchingParams,
  } = storeToRefs(store);

  return {
    jobs,
    jobsMeta,
    jobsPending,
    jobsFetchingParams,
    fetchJobsList: store.fetchList,
    fetchJobsListWithPreviousParams: store.fetchListWithPreviousParams,
    stopJobs: store.stop,
  };
};
