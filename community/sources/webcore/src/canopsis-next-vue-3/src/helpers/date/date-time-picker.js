/**
 * Returns time string from Vuetify item/value shapes.
 *
 * @param {*} time Time value.
 * @param {string} [defaultValue='00:00'] Default time.
 * @returns {string}
 */
export const normalizeTimeString = (time, defaultValue = '00:00') => {
  if (time?.raw) {
    return normalizeTimeString(time.raw, defaultValue);
  }

  if (time?.value) {
    return normalizeTimeString(time.value, defaultValue);
  }

  if (time?.title) {
    return normalizeTimeString(time.title, defaultValue);
  }

  return String(time ?? defaultValue);
};

/**
 * Immutable update time in the Date value.
 *
 * @param {Date} value Date value.
 * @param {*} [time='00:00'] Time in HH:mm format or Vuetify item.
 * @returns {Date}
 */
export const getDateObjectByTime = (value, time = '00:00') => {
  const newValue = value ? new Date(value.getTime()) : new Date();
  const [hours = 0, minutes = 0] = normalizeTimeString(time).split(':');

  newValue.setHours(parseInt(hours, 10) || 0, parseInt(minutes, 10) || 0, 0, 0);

  return newValue;
};

/**
 * Immutable update date in the Date value.
 *
 * @param {Date} value Date value.
 * @param {string} date Date in YYYY-MM-DD format.
 * @returns {Date}
 */
export const getDateObjectByDate = (value, date) => {
  const newValue = value ? new Date(value.getTime()) : new Date();
  const [year, month, day] = date.split('-');

  newValue.setFullYear(parseInt(year, 10), parseInt(month, 10) - 1, parseInt(day, 10));

  if (!value) {
    newValue.setHours(0, 0, 0, 0);
  } else {
    newValue.setSeconds(0, 0);
  }

  return newValue;
};
