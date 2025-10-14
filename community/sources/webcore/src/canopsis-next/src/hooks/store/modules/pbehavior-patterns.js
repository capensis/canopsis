import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Creates hooks for accessing the pbehavior patterns Vuex store module.
 * Provides access to getters and actions for managing pbehavior patterns operations.
 *
 * @returns {Object} An object containing store module utilities:
 * @property {import('vuex').Store} store - The Vuex store instance
 * @property {import('vuex').Module} module - The pbehavior patterns module instance
 * @property {Function} useGetters - Function to access module getters
 * @property {Function} useActions - Function to access module actions
 */
const usePbehaviorPatternsStoreModule = () => useStoreModuleHooks('pbehaviorPatterns');

/**
 * Custom hook for pbehavior patterns operations.
 * Provides convenient access to pbehavior patterns actions.
 *
 * @returns {Object} An object containing pbehavior patterns actions:
 * @property {Function} runAlarmFiltering - Action to run alarm filtering process
 *
 * @example
 * // Usage in a component
 * const { runAlarmFiltering } = usePbehaviorPatterns();
 * await runAlarmFiltering();
 */
export const usePbehaviorPatterns = () => {
  const { useActions } = usePbehaviorPatternsStoreModule();

  const actions = useActions({
    runAlarmFiltering: 'runAlarmFiltering',
  });

  return {
    ...actions,
  };
};
