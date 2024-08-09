import { useStoreModuleHooks } from '@/hooks/store';

const useEventsRecordStoreModule = () => useStoreModuleHooks('eventsRecord');

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
