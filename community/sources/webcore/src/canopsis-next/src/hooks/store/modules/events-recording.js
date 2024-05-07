import { useStoreModuleHooks } from '@/hooks/store';

const useEventsRecordingStoreModule = () => useStoreModuleHooks('eventsRecording');

export const useEventsRecording = () => {
  const { useActions } = useEventsRecordingStoreModule();

  const actions = useActions({
    launchEventsRecording: 'launch',
    stopEventsRecording: 'stop',
    removeEventsRecording: 'remove',
    fetchEventsRecordingsListWithoutStore: 'fetchListWithoutStore',
    fetchEventsRecordingEventsListWithoutStore: 'fetchEventsListWithoutStore',
  });

  return {
    ...actions,
  };
};
