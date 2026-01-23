import { PATTERN_DURATION_FORMAT, TIME_UNITS } from '@/constants';

import { sanitizeHtml } from '@/helpers/html';
import { convertDurationToString } from '@/helpers/date/duration';

import { useI18n } from '@/hooks/i18n';

/**
 * Hook for generating count message from alarm and entity counters.
 *
 * @returns {Function} Function that generates count message from alarm and entity counters.
 */
export const usePatternCountMessage = () => {
  const { t, tc } = useI18n();

  /**
   * Generates count message from alarm and entity counters.
   *
   * @param {Object} alarmCounter - Object containing alarm counter data.
   * @param {Object} entityCounter - Object containing entity counter data.
   * @returns {string} Formatted count message.
   */
  const getCountMessage = (alarmCounter, entityCounter) => {
    const messages = [];

    if (alarmCounter) {
      const alarmCount = alarmCounter?.count ?? 0;

      messages.push(tc('pattern.alarmFound', alarmCount, { count: alarmCount }));
    }

    if (entityCounter) {
      const entityCount = entityCounter?.count ?? 0;

      messages.push(tc('pattern.entityFound', entityCount, { count: entityCount }));
    }

    if (!messages.length) {
      return '';
    }

    const searchTime = convertDurationToString(
      Math.max(alarmCounter?.ms ?? 0, entityCounter?.ms ?? 0),
      PATTERN_DURATION_FORMAT,
      TIME_UNITS.millisecond,
    );

    return sanitizeHtml(t('pattern.found', { message: messages.join(', '), searchTime }));
  };

  return { getCountMessage };
};
