import { useStoreModuleHooks } from '@/hooks/store';

const useTicketStatusJobStoreModule = () => useStoreModuleHooks('ticketStatusJob');

/**
 * Custom hook for accessing getters and actions related to ticket status job Vuex store module.
 *
 * @returns {Object} An object containing getters and actions for the ticket status job store module.
 */
export const useTicketStatusJob = () => {
  const { useGetters, useActions } = useTicketStatusJobStoreModule();

  const getters = useGetters({
    ticketStatusJobs: 'items',
    ticketStatusJobsMeta: 'meta',
    ticketStatusJobsPending: 'pending',
  });

  const actions = useActions({
    fetchTicketStatusJobsList: 'fetchList',
    fetchTicketStatusJobsListWithPreviousParams: 'fetchListWithPreviousParams',
    updateTicketStatusJob: 'update',
    playTicketStatusJob: 'play',
    pauseTicketStatusJob: 'pause',
  });

  return {
    ...getters,
    ...actions,
  };
};
