import { usePbehaviorPatternsStore } from '@/stores/modules/pbehavior-patterns';

/**
 * Pbehavior patterns store facade.
 *
 * @returns {Object}
 */
export const usePbehaviorPatterns = () => {
  const store = usePbehaviorPatternsStore();

  return {
    runAlarmFiltering: store.runAlarmFiltering,
    checkPatternsPbehaviorsCount: store.checkPatternsPbehaviorsCount,
  };
};
