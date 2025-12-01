import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Hook for accessing the pattern Vuex store module
 * Creates hooks for accessing pattern module's getters and actions
 *
 * @returns {Object} An object containing store module utilities
 * @property {Object} store - The Vuex store instance
 * @property {Object} module - The pattern Vuex module
 * @property {Function} useGetters - Function to access pattern module getters
 * @property {Function} useActions - Function to access pattern module actions
 */
const usePatternStoreModule = () => useStoreModuleHooks('pattern');

/**
 * Hook for managing pattern-related operations and state
 *
 * @returns {Object} An object containing pattern-related actions
 * @property {Function} checkPatternsEntitiesCount - Checks the count of entities matching patterns
 * @property {Function} checkPatternsAlarmsCount - Checks the count of alarms matching patterns
 */
export const usePattern = () => {
  const { useActions } = usePatternStoreModule();

  const actions = useActions({
    checkPatternsEntitiesCount: 'checkPatternsEntitiesCount',
    checkPatternsAlarmsCount: 'checkPatternsAlarmsCount',
  });

  return {
    ...actions,
  };
};
