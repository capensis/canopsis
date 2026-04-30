import { computed, ref, unref } from 'vue';

import { useI18n } from '@/hooks/i18n';
import { usePendingHandler } from '@/hooks/query/pending';
import { usePattern } from '@/hooks/store/modules/pattern';

/**
 * Hook for fetching pattern counters
 *
 * @param {Object} options - Options for the hook
 * @param {Ref|Function} [options.counterMethod] - Custom counter method to use
 * @param {Ref|boolean} [options.entityCountersType] - Whether to use entity counters type
 * @param {Ref|boolean} [options.bothCounters] - Whether to show both alarms and entities counters
 * @param {ComputedRef|Ref|boolean} options.hasError - Whether pattern has validation error
 * @param {ComputedRef} options.patterns - Computed patterns to check
 * @returns {Object} An object containing counters, pending state, checkFilter function, and computed properties
 */
export const usePatternCounters = ({
  counterMethod,
  entityCountersType,
  hasError,
  patterns,
}) => {
  const { t } = useI18n();
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

  const hasAllInCounter = computed(() => (counters.value?.all?.count > 0));
  const allOverLimit = computed(() => (counters.value?.all?.over_limit ?? false));
  const allCount = computed(() => (counters.value?.all?.count ?? 0));

  const patternsCountMessageBind = computed(() => {
    let entityCounter;
    let alarmCounter;

    if (unref(hasError)) {
      return { errorMessage: t('pattern.errors.required') };
    }

    if (unref(entityCountersType)) {
      entityCounter = counters.value?.all;
    } else {
      alarmCounter = counters.value?.all;
      entityCounter = counters.value?.entities;
    }

    return {
      alarmCounter,
      entityCounter,
    };
  });

  return {
    counters,
    pending,
    checkFilter,
    hasAllInCounter,
    allOverLimit,
    allCount,
    patternsCountMessageBind,
  };
};
