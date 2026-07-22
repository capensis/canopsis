import moment from 'moment-timezone'

import {
  DATETIME_FORMATS,
  DATETIME_INTERVAL_TYPES,
  QUICK_RANGES,
  SAMPLINGS,
  TIME_UNITS,
} from '@/constants'

import {
  convertDateToEndOfUnitMoment,
  convertDateToMomentByTimezone,
  convertDateToStartOfUnitMoment,
  getLocalTimezone,
  getNowTimestamp,
} from './date'

/**
 * @typedef {Object} IntervalForm
 * @property {number | string} from
 * @property {number | string} to
 */

/**
 * Convert a date interval string to moment date object.
 *
 * @param {string} string
 * @param {string} type
 * @returns {Moment}
 */
export const convertStringToDateInterval = (string, type) => {
  const matches = string.match(
    /^(?<initial>now|today)((?<operator>[+--])(?<deltaValue>\d+)(?<deltaUnit>[hdwmMy]{1}))?(\/(?<roundUnit>[hdwmMy]{1}))?$/,
  )

  if (matches?.groups) {
    const { initial, operator, deltaValue, deltaUnit, roundUnit } = matches.groups
    const preparedRoundUnit = roundUnit === TIME_UNITS.week ? 'isoWeek' : roundUnit
    const isToday = initial === 'today'
    const isStart = type === DATETIME_INTERVAL_TYPES.start
    let result = moment()

    if (isToday) {
      result =
        isStart || deltaValue
          ? convertDateToStartOfUnitMoment(getNowTimestamp(), TIME_UNITS.day)
          : convertDateToEndOfUnitMoment(getNowTimestamp(), TIME_UNITS.day)
    }

    if (preparedRoundUnit) {
      if (isStart) {
        result.startOf(preparedRoundUnit)
      } else {
        result.endOf(preparedRoundUnit)
      }
    }

    if (deltaValue && deltaUnit) {
      if (operator === '+') {
        result.add(deltaValue, deltaUnit)
      } else {
        result.subtract(deltaValue, deltaUnit)
      }
    }

    return result
  }

  throw new Error('Date string pattern not recognized')
}

/**
 * Parse date in every format to moment object.
 *
 * @param {number|string} date
 * @param {string} type
 * @param {string} [format=DATETIME_FORMATS.datePicker]
 * @param {string} [unit]
 * @returns {Moment}
 */
export const convertDateIntervalToMoment = (
  date,
  type,
  format = DATETIME_FORMATS.datePicker,
  unit,
) => {
  const momentDate = Number.isFinite(date) ? moment.unix(date) : moment(date, format, true)

  if (!momentDate.isValid()) {
    return convertStringToDateInterval(date, type)
  }

  if (unit) {
    const method = type === DATETIME_INTERVAL_TYPES.start ? 'startOf' : 'endOf'

    return momentDate[method](unit)
  }

  return momentDate
}

/**
 * Parse date in every format to moment object by timezone.
 *
 * @param {number|string} date
 * @param {string} type
 * @param {string} [format=DATETIME_FORMATS.datePicker]
 * @param {string} [unit=SAMPLINGS.hour]
 * @param {string} [timezone=getLocalTimezone()]
 * @returns {Moment}
 */
export const convertDateIntervalToMomentByTimezone = (
  date,
  type,
  format = DATETIME_FORMATS.datePicker,
  unit = SAMPLINGS.hour,
  timezone = getLocalTimezone(),
) =>
  convertDateToMomentByTimezone(
    convertDateIntervalToMoment(date, type, format, unit),
    getLocalTimezone(),
    timezone,
  )

/**
 * Convert date interval to unix timestamp.
 *
 * @param {number|string} date
 * @param {string} type
 * @param {string} format
 * @param {string} unit
 * @returns {number}
 */
export const convertDateIntervalToTimestamp = (date, type, format, unit) =>
  convertDateIntervalToMoment(date, type, format, unit).unix()

/**
 * Convert date interval to unix timestamp by timezone.
 *
 * @param {number|string} date
 * @param {string} type
 * @param {string} format
 * @param {string} unit
 * @param {string} [timezone=getLocalTimezone()]
 * @returns {number}
 */
export const convertDateIntervalToTimestampByTimezone = (
  date,
  type,
  format,
  unit,
  timezone = getLocalTimezone(),
) => convertDateIntervalToMomentByTimezone(date, type, format, unit, timezone).unix()

/**
 * Convert from value to timestamp.
 *
 * @param {number|string} date
 * @param {string} [format=DATETIME_FORMATS.datePicker]
 * @param {string} [unit]
 * @returns {number}
 */
export const convertStartDateIntervalToTimestamp = (
  date,
  format = DATETIME_FORMATS.datePicker,
  unit,
) => convertDateIntervalToTimestamp(date, DATETIME_INTERVAL_TYPES.start, format, unit)

/**
 * Convert to value to timestamp.
 *
 * @param {number|string} date
 * @param {string} [format=DATETIME_FORMATS.datePicker]
 * @param {string} [unit]
 * @returns {number}
 */
export const convertStopDateIntervalToTimestamp = (
  date,
  format = DATETIME_FORMATS.datePicker,
  unit,
) => convertDateIntervalToTimestamp(date, DATETIME_INTERVAL_TYPES.stop, format, unit)

/**
 * Convert metric interval form to timestamp interval.
 *
 * @param {Object} options
 * @param {{ from: number|string, to: number|string }} options.interval
 * @param {string} [options.timezone=getLocalTimezone()]
 * @returns {{ from: number, to: number }}
 */
export const convertMetricIntervalToTimestamp = ({
  interval,
  timezone = getLocalTimezone(),
} = {}) => {
  const { from, to } = interval

  return {
    from: convertDateIntervalToTimestampByTimezone(
      from,
      DATETIME_INTERVAL_TYPES.start,
      DATETIME_FORMATS.dateTimePicker,
      undefined,
      timezone,
    ),
    to: convertDateIntervalToTimestampByTimezone(
      to,
      DATETIME_INTERVAL_TYPES.stop,
      DATETIME_FORMATS.dateTimePicker,
      undefined,
      timezone,
    ),
  }
}

/**
 * Find range object.
 *
 * @param {string|number} start
 * @param {string|number} stop
 * @param {Object|Array} ranges
 * @param {Object} defaultValue
 * @returns {Object}
 */
export const findQuickRangeValue = (
  start,
  stop,
  ranges = QUICK_RANGES,
  defaultValue = QUICK_RANGES.custom,
) =>
  Object.values(ranges).find((range) => start === range.start && stop === range.stop) ??
  defaultValue
