import { ref } from 'vue';

import { DATETIME_FORMATS } from '@/constants';

import { convertDateToTimezoneDateString, getLocalTimezone } from '@/helpers/date/date';

import { useI18n } from '@/hooks/i18n';
import { useInfo } from '@/hooks/store/modules/info';

/**
 * Hook for handling date formatting in pbehavior components with timezone support.
 * Provides utilities for formatting interval dates and recurring rule end dates while respecting
 * timezone settings from the application configuration.
 *
 * @returns {Object} An object containing date formatting utilities
 * @property {Ref<string>} timezone - Reactive reference to the current timezone
 * @property {Ref<boolean>} shownUserTimezone - Whether to display user timezone based on store settings
 * @property {Function} formatIntervalDate - Formats interval dates with timezone consideration
 * @property {Function} formatRruleEndDate - Formats recurring rule end dates
 */
export const usePbehaviorDateFormat = () => {
  const { t } = useI18n();
  const { shownUserTimezone } = useInfo();

  const timezone = ref(getLocalTimezone());

  /**
   * Formats an interval date string based on whether the item has a recurring rule.
   * Uses medium format for recurring items and long format for non-recurring items.
   *
   * @function formatIntervalDate
   * @param {Object} item - The item containing date information
   * @param {string|Date} item[field] - The date value to format
   * @param {boolean} [item.rrule] - Whether the item has a recurring rule
   * @param {string} field - The field name containing the date value
   * @returns {string} Formatted date string according to timezone and format rules
   */
  const formatIntervalDate = (item, field) => {
    const date = item[field];
    const format = item.rrule ? DATETIME_FORMATS.medium : DATETIME_FORMATS.long;

    return convertDateToTimezoneDateString(date, timezone.value, format);
  };

  /**
   * Formats the end date of a recurring rule if present.
   * Returns a hyphen for non-recurring items, the formatted date for items with an end date,
   * or a translated 'undefined' string for recurring items without an end date.
   *
   * @function formatRruleEndDate
   * @param {Object} item - The item containing recurring rule information
   * @param {boolean} item.rrule - Whether the item has a recurring rule
   * @param {string|null} item.rrule_end - The end date of the recurring rule
   * @returns {string} Formatted end date string, '-', or translated 'undefined'
   */
  const formatRruleEndDate = (item) => {
    if (!item.rrule) {
      return '-';
    }

    return item.rrule_end
      ? convertDateToTimezoneDateString(item.rrule_end, timezone.value, DATETIME_FORMATS.long)
      : t('common.undefined');
  };

  return {
    timezone,
    shownUserTimezone,

    formatIntervalDate,
    formatRruleEndDate,
  };
};
