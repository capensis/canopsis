import {
  DATETIME_FORMATS,
  DEFAULT_USER_INTERFACE_DATE_FORMAT,
  USER_INTERFACE_DATE_FORMATS,
} from '@/constants';

const defaultDateFormat = USER_INTERFACE_DATE_FORMATS.dateTime;
let dateFormatProvider = () => DEFAULT_USER_INTERFACE_DATE_FORMAT;

/**
 * Register the provider owned by the reactive `$system` context.
 *
 * @param {Function} provider
 */
export const registerUserInterfaceDateFormatProvider = (provider) => {
  dateFormatProvider = provider;
};

/**
 * Return the frontend formats for the value supplied by `$system`.
 *
 * @param {string} [value]
 * @returns {{ format: string, dateFormat: string }}
 */
export const getUserInterfaceDateFormat = (value = dateFormatProvider()) =>
  Object.values(USER_INTERFACE_DATE_FORMATS).find(format => format.value === value) ??
  defaultDateFormat;

/**
 * Replace a standard display/picker format with the current UI preference.
 * Internal transport formats such as YYYY-MM-DD are intentionally preserved.
 *
 * @param {string} format
 * @returns {string}
 */
export const resolveUserInterfaceDateFormat = (format) => {
  const currentFormat = getUserInterfaceDateFormat();

  if (format === DATETIME_FORMATS.longWithTimezone) {
    return `${currentFormat.format} z`;
  }

  if (
    format === DATETIME_FORMATS.short ||
    format === DATETIME_FORMATS.shortWithDayOfWeek
  ) {
    return currentFormat.dateFormat;
  }

  if (
    format === DATETIME_FORMATS.long ||
    format === DATETIME_FORMATS.longWithDayOfWeek ||
    format === DATETIME_FORMATS.dateTimePicker ||
    format === DATETIME_FORMATS.dateTimePickerWithSeconds
  ) {
    return currentFormat.format;
  }

  return format;
};
