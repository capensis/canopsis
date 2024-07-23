import { useStoreModuleHooks } from '@/hooks/store';

const useEventsRecordStoreModule = () => useStoreModuleHooks('eventsRecord');

export const useEventsRecord = () => {
  const { useActions } = useEventsRecordStoreModule();

  const actions = useActions({
    startEventsRecord: 'start',
    stopEventsRecord: 'stop',
    removeEventsRecord: 'remove',
    fetchEventsRecordsListWithoutStore: 'fetchListWithoutStore',
    fetchEventsRecordCurrentWithoutStore: 'fetchCurrentWithoutStore',
    fetchEventsRecordEventsListWithoutStore: 'fetchEventsListWithoutStore',
  });

  return {
    ...actions,
  };
};
