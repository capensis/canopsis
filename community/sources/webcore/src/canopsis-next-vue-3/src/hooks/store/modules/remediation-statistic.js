import { storeToRefs } from 'pinia';

import { useRemediationStatisticStore } from '@/stores/modules/remediation-statistic';

/**
 * Remediation statistic store facade.
 *
 * @returns {Object}
 */
export const useRemediationStatistic = () => {
  const store = useRemediationStatisticStore();
  const { metrics, meta, pending, fetchingParams } = storeToRefs(store);

  return {
    remediationStatisticMetrics: metrics,
    remediationStatisticMeta: meta,
    remediationStatisticPending: pending,
    remediationStatisticFetchingParams: fetchingParams,
    fetchRemediationStatisticList: store.fetchList,
    fetchRemediationStatisticListWithPreviousParams: store.fetchListWithPreviousParams,
  };
};
