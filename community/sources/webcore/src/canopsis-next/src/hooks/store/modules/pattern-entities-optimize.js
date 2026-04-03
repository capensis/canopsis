import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Creates hooks for accessing the pattern entities optimize Vuex store module.
 * Provides access to getters and actions for managing pattern entities optimize operations.
 *
 * @returns {Object} An object containing store module utilities:
 * @property {import('vuex').Store} store - The Vuex store instance
 * @property {import('vuex').Module} module - The pattern entities optimize module instance
 * @property {Function} useGetters - Function to access module getters
 * @property {Function} useActions - Function to access module actions
 */
const usePatternEntitiesOptimizeStoreModule = () => useStoreModuleHooks('pattern/entitiesOptimize');

/**
 * Custom hook for pattern entities optimize operations.
 * Provides convenient access to entities optimize actions.
 *
 * @returns {Object} An object containing entities optimize actions:
 * @property {Function} optimize - Action to optimize entities
 * @property {Function} fetchOptimizeStatus - Action to fetch optimize status by id
 * @property {Function} remove - Action to remove optimize by id
 *
 * @example
 * // Usage in a component
 * const { optimize, fetchOptimizeStatus, remove } = usePatternEntitiesOptimize();
 * const result = await optimize({ data: { patterns: [...] } });
 * const status = await fetchOptimizeStatus({ id: 'task-id' });
 */
export const usePatternEntitiesOptimize = () => {
  const { useActions } = usePatternEntitiesOptimizeStoreModule();

  const actions = useActions({
    optimizeEntities: 'optimize',
    fetchOptimizeEntitiesStatus: 'fetchOptimizeStatus',
    updateOptimization: 'update',
    removeOptimization: 'remove',
  });

  return {
    ...actions,
  };
};
