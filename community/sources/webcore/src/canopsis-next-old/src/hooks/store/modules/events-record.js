import { useStoreModuleHooks } from '@/hooks/store';

const useEventsRecordStoreModule = () => useStoreModuleHooks('eventsRecord');

/**
 * Custom hook for accessing actions related to events record Vuex store module.
 *
 * @returns {Object} An object containing functions to access actions related to events record Vuex store module.
 *
 * @example
 * // Usage example
 * import { useEventsRecord } from './path/to/useEventsRecord';
 *
 * setup() {
 *   const {
 *     createEventsRecordExport,
 *     fetchEventsRecordExport,
 *     removeEventsRecordEvent,
 *     bulkRemoveEventsRecordEvent,
 *     playbackEventsRecordEvents,
 *     stopPlaybackEventsRecordEvents,
 *     removeEventsRecord,
 *     fetchEventsRecordsListWithoutStore,
 *     fetchEventsRecordEventsListWithoutStore,
 *   } = useEventsRecord();
 *
 *   // Access and use the actions as needed
 * }
 */
export const useEventsRecord = () => {
  const { useActions } = useEventsRecordStoreModule();

  const actions = useActions({
    createEventsRecordExport: 'createExport',
    fetchEventsRecordExport: 'fetchExport',
    removeEventsRecordEvent: 'removeEvent',
    bulkRemoveEventsRecordEvent: 'bulkRemoveEvent',
    playbackEventsRecordEvents: 'playback',
    stopPlaybackEventsRecordEvents: 'stopPlayback',
    removeEventsRecord: 'remove',
    fetchEventsRecordsListWithoutStore: 'fetchListWithoutStore',
    fetchEventsRecordEventsListWithoutStore: 'fetchEventsListWithoutStore',
  });

  return {
    ...actions,
  };
};
