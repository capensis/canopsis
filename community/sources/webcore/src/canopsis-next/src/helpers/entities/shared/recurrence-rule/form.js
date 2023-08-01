import { RRule } from 'rrule';
import { isArray, isNumber, mapValues, pick } from 'lodash';

/**
 * @typedef {Object} Weekday
 * @property {number} weekday
 */

/**
 * @typedef {Weekday | string | number} ByWeekday
 */

/**
 * @typedef {Object} RecurrenceRuleOptions
 * @property {number} [freq]
 * @property {Date} [dtstart]
 * @property {number} [interval]
 * @property {number} [count]
 * @property {Date} [until]
 * @property {!Weekday} [wkst]
 * @property {string} [tzid]
 * @property {number | number[]} [bysetpos]
 * @property {number | number[]} [bymonth]
 * @property {number | number[]} [bymonthday]
 * @property {number | number[]} [bynmonthday]
 * @property {number[]} [bynmonthday]
 * @property {number | number[]} [byyearday]
 * @property {number | number[]} [byweekno]
 * @property {ByWeekday | ByWeekday[]} [byweekday]
 * @property {number[][]} [bynweekday]
 * @property {number | number[]} [byhour]
 * @property {number | number[]} [byminute]
 * @property {number | number[]} [bysecond]
 * @property {number} [byeaster]
 */

/**
 * @typedef {RecurrenceRuleOptions} RecurrenceRuleFormOptions
 * @property {Weekday | string} wkst
 */

/**
 * @typedef {Object} RecurrenceRuleFormAdvancedOptions
 */

/**
 * Convert string or array to string
 *
 * @param {number | string | string[] | number[]} value
 * @return {string}
 */
const prepareRecurrenceRuleOption = (value) => {
  if (value) {
    return (isArray(value) ? value.join(',') : String(value));
  }

  return '';
};

/**
 * Prepare rrule options to form
 *
 * @param {RecurrenceRuleOptions} recurrenceRule
 * @return {RecurrenceRuleFormOptions}
 */
export const recurrenceRuleToFormOptions = recurrenceRule => ({
  freq: recurrenceRule.freq || null,
  count: recurrenceRule.count || '',
  interval: recurrenceRule.interval || '',
  byweekday: recurrenceRule.byweekday ? recurrenceRule.byweekday.map(v => v.weekday) : [],
  wkst: recurrenceRule.wkst ? recurrenceRule.wkst.weekday : '',
  bymonth: recurrenceRule.bymonth || [],
  bysetpos: prepareRecurrenceRuleOption(recurrenceRule.bysetpos),
  bymonthday: prepareRecurrenceRuleOption(recurrenceRule.bymonthday),
  byyearday: prepareRecurrenceRuleOption(recurrenceRule.byyearday),
  byweekno: prepareRecurrenceRuleOption(recurrenceRule.byweekno),
  byhour: prepareRecurrenceRuleOption(recurrenceRule.byhour),
  byminute: prepareRecurrenceRuleOption(recurrenceRule.byminute),
  bysecond: prepareRecurrenceRuleOption(recurrenceRule.bysecond),
});

/**
 * Prepare form options to rrule form
 *
 * @param {RecurrenceRuleFormOptions} options
 * @param {string[]} advancedFields
 * @return {RecurrenceRuleOptions}
 */
export const formOptionsToRecurrenceRuleOptions = (options, advancedFields = []) => {
  const recurrenceRuleOptions = {
    freq: options.freq,
    ...mapValues(
      pick(options, advancedFields),
      o => o.split(',').filter(v => v),
    ),
  };

  if (isNumber(options.count)) {
    recurrenceRuleOptions.count = options.count;
  }

  if (isNumber(options.interval)) {
    recurrenceRuleOptions.interval = options.interval;
  }

  if (options.freq !== RRule.YEARLY && options.byweekday.length) {
    recurrenceRuleOptions.byweekday = options.byweekday;
  }

  if (options.bymonth.length) {
    recurrenceRuleOptions.bymonth = options.bymonth;
  }

  if (isNumber(options.wkst)) {
    recurrenceRuleOptions.wkst = options.wkst.weekday;
  }

  return recurrenceRuleOptions;
};
