/**
 * Immutable update time in the Date value
 *
 * @param {Date} value
 * @param {string} [time = '00:00'];
 * @returns {Date}
 */
export const getDateObjectByTime = (value, time = '00:00') => {
  const newValue = value ? new Date(value.getTime()) : new Date();
  const [hours = 0, minutes = 0] = time.split(':');

  newValue.setHours(parseInt(hours, 10) || 0, parseInt(minutes, 10) || 0, 0, 0);

  return newValue;
};

/**
 * Immutable update date in the Date value
 *
 * @param {Date} value
 * @param {string} date
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
