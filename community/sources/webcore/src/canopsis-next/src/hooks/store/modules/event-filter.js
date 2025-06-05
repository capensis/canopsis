import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Hook to use the Event Filter Store Module.
 *
 * @returns {Object} An object containing getters and actions for event filters.
 */
const useEventFilterStoreModule = () => useStoreModuleHooks('eventFilter');

/**
 * Hook to access event filter store.
 *
 * @returns {Object} An object containing:
 * - Actions to fetch event filter rules and manage them.
 */
export const useEventFilterStore = () => {
  const { useActions } = useEventFilterStoreModule();

  const actions = useActions({
    fetchEventFilterRulesListWithoutStore: 'fetchRulesListWithoutStore',
  });

  return {
    ...actions,
  };
};
