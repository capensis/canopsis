import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Creates an instance of the associative table Vuex store module hooks.
 *
 * @returns {Object} Store module hooks for the 'associativeTable' namespace
 */
export const useAssociativeTableStoreModule = () => useStoreModuleHooks('associativeTable');

/**
 * Hook for dispatching associative table API actions.
 *
 * @returns {Object} Fetch, create, update, and remove associative table actions
 */
export const useAssociativeTable = () => {
  const { useActions } = useAssociativeTableStoreModule();

  return useActions({
    fetchAssociativeTable: 'fetch',
    createAssociativeTable: 'create',
    updateAssociativeTable: 'update',
    removeAssociativeTable: 'remove',
  });
};
