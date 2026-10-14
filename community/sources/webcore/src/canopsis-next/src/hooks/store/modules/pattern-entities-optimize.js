import { usePatternEntitiesOptimizeStore } from '@/stores/modules/pattern-entities-optimize';

/**
 * Pattern entities optimization store facade.
 *
 * @returns {Object}
 */
export const usePatternEntitiesOptimize = () => {
  const store = usePatternEntitiesOptimizeStore();

  return {
    optimizeEntities: store.optimize,
    fetchOptimizeEntitiesStatus: store.fetchOptimizeStatus,
    updateOptimization: store.update,
    removeOptimization: store.remove,
  };
};
