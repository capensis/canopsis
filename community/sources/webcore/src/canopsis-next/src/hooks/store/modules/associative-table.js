import { useAssociativeTableStore } from '@/stores/modules/associative-table';

/**
 * Hook for dispatching associative table API actions.
 *
 * @returns {Object} Fetch, create, update and remove associative table actions.
 */
export const useAssociativeTable = () => {
  const associativeTableStore = useAssociativeTableStore();

  return {
    fetchAssociativeTable: associativeTableStore.fetch,
    createAssociativeTable: associativeTableStore.create,
    updateAssociativeTable: associativeTableStore.update,
    removeAssociativeTable: associativeTableStore.remove,
  };
};
