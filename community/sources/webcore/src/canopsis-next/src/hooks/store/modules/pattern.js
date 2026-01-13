import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Creates hooks for accessing the pattern Vuex store module.
 * Provides access to getters and actions for managing pattern operations.
 *
 * @returns {Object} An object containing store module utilities:
 * @property {import('vuex').Store} store - The Vuex store instance
 * @property {import('vuex').Module} module - The pattern module instance
 * @property {Function} useGetters - Function to access module getters
 * @property {Function} useActions - Function to access module actions
 */
const usePatternStoreModule = () => useStoreModuleHooks('pattern');

/**
 * Custom hook for pattern operations.
 * Provides convenient access to pattern getters and actions.
 *
 * @returns {Object} An object containing pattern getters and actions:
 * @property {Object} patternsMeta - Getter for pattern metadata
 * @property {boolean} patternsPending - Getter for pattern loading state
 * @property {Array} patterns - Getter for pattern items
 * @property {Function} fetchPatternsList - Action to fetch patterns list
 * @property {Function} createPattern - Action to create a pattern
 * @property {Function} updatePattern - Action to update a pattern
 * @property {Function} removePattern - Action to remove a pattern
 * @property {Function} bulkRemovePatterns - Action to bulk remove patterns
 * @property {Function} fetchPatternsListWithPreviousParams - Action to fetch patterns list with previous params
 * @property {Function} fetchPatternsListWithoutStore - Action to fetch patterns list without store
 *
 * @example
 * // Usage in a component
 * const { patterns, patternsPending, fetchPatternsList } = usePattern();
 * await fetchPatternsList({ params: { page: 1, limit: 10 } });
 */
export const usePattern = () => {
  const { useGetters, useActions } = usePatternStoreModule();

  const getters = useGetters({
    patternsMeta: 'meta',
    patternsPending: 'pending',
    patterns: 'items',
  });

  const actions = useActions({
    fetchPatternsList: 'fetchList',
    createPattern: 'create',
    updatePattern: 'update',
    removePattern: 'remove',
    bulkRemovePatterns: 'bulkRemove',
    fetchPatternsListWithPreviousParams: 'fetchListWithPreviousParams',
    fetchPatternsListWithoutStore: 'fetchListWithoutStore',
  });

  return {
    ...getters,
    ...actions,
  };
};
