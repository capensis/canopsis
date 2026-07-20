import { storeToRefs } from 'pinia'

import { useEventsRecordStore } from '@/stores/modules/events-record'

/**
 * Events record store facade.
 *
 * @returns {Object}
 */
export const useEventsRecord = () => {
  const eventsRecordStore = useEventsRecordStore()
  const {
    currentPending,
    currentStatus,
    recordings,
    resendings,
    limit,
    recordingsById,
    resendingsById,
  } = storeToRefs(eventsRecordStore)

  return {
    eventsRecordCurrentPending: currentPending,
    eventsRecordCurrentStatus: currentStatus,
    eventsRecordRecordings: recordings,
    eventsRecordResendings: resendings,
    eventsRecordLimit: limit,
    eventsRecordRecordingsById: recordingsById,
    eventsRecordResendingsById: resendingsById,

    setEventsRecordCurrentStatus: eventsRecordStore.setCurrentStatus,
    fetchEventsRecordsListWithoutStore: eventsRecordStore.fetchListWithoutStore,
    fetchEventsRecordEventsListWithoutStore: eventsRecordStore.fetchEventsListWithoutStore,
    createEventsRecordExport: eventsRecordStore.createExport,
    fetchEventsRecordExport: eventsRecordStore.fetchExport,
    playbackEventsRecord: eventsRecordStore.playback,
    stopEventsRecordPlayback: eventsRecordStore.stopPlayback,
    removeEventsRecord: eventsRecordStore.remove,
    removeEventsRecordEvent: eventsRecordStore.removeEvent,
    bulkRemoveEventsRecordEvents: eventsRecordStore.bulkRemoveEvent,
    fetchCurrentEventsRecord: eventsRecordStore.fetchCurrent,
    startCurrentEventsRecord: eventsRecordStore.startCurrent,
    stopCurrentEventsRecord: eventsRecordStore.stopCurrent,
  }
}
