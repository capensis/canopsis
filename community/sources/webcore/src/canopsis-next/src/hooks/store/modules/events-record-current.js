import { useStoreModuleHooks } from '@/hooks/store';

const useEventsRecordCurrentStoreModule = () => useStoreModuleHooks('eventsRecord/current');

/**
 * Retrieves getters and actions related to the events record current store module.
 *
 * @returns {Object} An object containing getters and actions for interacting with the events record current
 * store module.
 *
 * @example
 * // Usage example
 * import { useEventsRecordCurrentStoreModule } from './path/to/useEventsRecordCurrentStoreModule';
 *
 * const {
 *   current,
 *   pending,
 *   fetchEventsRecordCurrent,
 *   startEventsRecordCurrent,
 *   stopEventsRecordCurrent
 * } = useEventsRecordCurrent();
 */

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
