import { computed, ref, unref } from 'vue';
import { isEmpty } from 'lodash';

import { PATTERN_DURATION_FORMAT, TIME_UNITS } from '@/constants';

import { sanitizeHtml } from '@/helpers/html';
import { convertDurationToString } from '@/helpers/date/duration';

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
  bothCounters,
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

  const checkFilterMessages = computed(() => {
    if (unref(hasError)) {
      return t('pattern.errors.required');
    }

    if (isEmpty(counters.value)) {
      return '';
    }

    const alarmsCount = counters.value?.all?.count ?? 0;
    const allDuration = convertDurationToString(
      counters.value?.all?.ms,
      PATTERN_DURATION_FORMAT,
      TIME_UNITS.millisecond,
    );
    const durationMessage = t('pattern.searchTime', { duration: allDuration });

    let message = '';

    if (unref(entityCountersType)) {
      const entitiesCount = counters.value?.entity_pattern?.count ?? 0;

      message = t('pattern.entitiesCount', { entitiesCount });
    } else if (unref(bothCounters)) {
      const entitiesCount = counters.value?.entities?.count ?? 0;

      message = t('pattern.alarmsEntitiesCount', {
        alarmsCount,
        entitiesCount,
      });
    } else {
      message = t('pattern.alarmsCount', { alarmsCount });
    }

    return sanitizeHtml(`${message} / ${durationMessage}`);
  });

  return {
    counters,
    pending,
    checkFilter,
    hasAllInCounter,
    allOverLimit,
    allCount,
    checkFilterMessages,
  };
};
