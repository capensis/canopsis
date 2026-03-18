import { useStoreModuleHooks } from '@/hooks/store';

const useTicketStatusJobStoreModule = () => useStoreModuleHooks('ticketStatusJob');

/**
 * Custom hook for accessing actions related to ticket status job Vuex store module.
 *
 * @returns {Object} An object containing functions to access actions related to ticket status job Vuex store module.
 */
export const useTicketStatusJob = () => {
  const { useActions } = useTicketStatusJobStoreModule();

  const actions = useActions({
    fetchTicketStatusJobsListWithoutStore: 'fetchListWithoutStore',
    updateTicketStatusJob: 'update',
    playTicketStatusJob: 'play',
    pauseTicketStatusJob: 'pause',
  });

  return {
    ...actions,
  };
};
