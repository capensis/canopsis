import { useStoreModuleHooks } from '@/hooks/store';

const useEventsRecordCurrentStoreModule = () => useStoreModuleHooks('eventsRecord/current');

export const useEventsRecordCurrent = () => {
  const { useGetters, useActions } = useEventsRecordCurrentStoreModule();

  const getters = useGetters(['current', 'pending']);

  const actions = useActions({
    fetchEventsRecordCurrent: 'fetchCurrent',
    startEventsRecordCurrent: 'start',
    stopEventsRecordCurrent: 'stop',
  });

  return {
    ...getters,
    ...actions,
  };
};
