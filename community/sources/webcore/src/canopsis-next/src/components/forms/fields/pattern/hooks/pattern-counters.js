import { ref, unref } from 'vue';

import { usePendingHandler } from '@/hooks/query/pending';
import { usePattern } from '@/hooks/store/modules/pattern';

/**
 * Hook for fetching pattern counters
 *
 * @param {Object} options - Options for the hook
 * @param {Ref|Function} [options.counterMethod] - Custom counter method to use
 * @param {Ref|boolean} [options.entityCountersType] - Whether to use entity counters type
 * @param {ComputedRef} options.patterns - Computed patterns to check
 * @returns {Object} An object containing counters, pending state, and checkFilter function
 */
export const usePatternCounters = ({ counterMethod, entityCountersType, patterns }) => {
  const counters = ref({});
  const { checkPatternsEntitiesCount, checkPatternsAlarmsCount } = usePattern();

  const { pending, handler: checkFilter } = usePendingHandler(async () => {
    const unwrappedCounterMethod = unref(counterMethod);
    const unwrappedEntityCountersType = unref(entityCountersType);
    const unwrappedPatterns = unref(patterns);

    const method = unwrappedCounterMethod ?? {
      [true]: checkPatternsAlarmsCount,
      [unwrappedEntityCountersType]: checkPatternsEntitiesCount,
    }.true;

    try {
      counters.value = await method({ data: unwrappedPatterns });
    } catch (err) {
      console.error(err);

      counters.value = {};
    }
  });

  return {
    counters,
    pending,
    checkFilter,
  };
};
