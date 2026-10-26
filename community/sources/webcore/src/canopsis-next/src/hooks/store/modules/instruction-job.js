import { storeToRefs } from 'pinia';

import { useInstructionJobStore } from '@/stores/modules/instruction-job';

/**
 * Instruction job store facade.
 *
 * @returns {Object}
 */
export const useInstructionJob = () => {
  const instructionJobStore = useInstructionJobStore();
  const { items, meta, pending, fetchingParams } = storeToRefs(instructionJobStore);

  return {
    instructionJobs: items,
    instructionJobsMeta: meta,
    instructionJobsPending: pending,
    instructionJobsFetchingParams: fetchingParams,
    fetchInstructionJobsList: instructionJobStore.fetchList,
    fetchInstructionJobsListWithPreviousParams: instructionJobStore.fetchListWithPreviousParams,
    stopInstructionJobs: instructionJobStore.stop,
  };
};
