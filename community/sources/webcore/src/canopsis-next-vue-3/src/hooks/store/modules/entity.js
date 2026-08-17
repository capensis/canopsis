import { storeToRefs } from 'pinia';

import { useEntityStore } from '@/stores/modules/entity';

/**
 * Hook for managing entity-related operations and state
 *
 * @returns {Object} An object containing entity-related actions
 */
export const useEntity = () => {
  const entityStore = useEntityStore();
  const refs = storeToRefs(entityStore);

  return {
    ...refs,
    ...entityStore,
  };
};
