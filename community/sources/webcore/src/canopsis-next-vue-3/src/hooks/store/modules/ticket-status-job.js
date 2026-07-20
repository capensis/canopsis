import { storeToRefs } from 'pinia'

import { useTicketStatusJobStore } from '@/stores/modules/ticket-status-job'

/**
 * Ticket status job store facade.
 *
 * @returns {Object}
 */
export const useTicketStatusJob = () => {
  const ticketStatusJobStore = useTicketStatusJobStore()
  const { items, meta, pending, fetchingParams } = storeToRefs(ticketStatusJobStore)

  return {
    ticketStatusJobs: items,
    ticketStatusJobsMeta: meta,
    ticketStatusJobsPending: pending,
    ticketStatusJobsFetchingParams: fetchingParams,

    fetchTicketStatusJobsList: ticketStatusJobStore.fetchList,
    fetchTicketStatusJobsListWithPreviousParams: ticketStatusJobStore.fetchListWithPreviousParams,
    createTicketStatusJob: ticketStatusJobStore.create,
    updateTicketStatusJob: ticketStatusJobStore.update,
    removeTicketStatusJob: ticketStatusJobStore.remove,
    playTicketStatusJobs: ticketStatusJobStore.play,
    pauseTicketStatusJobs: ticketStatusJobStore.pause,
  }
}
