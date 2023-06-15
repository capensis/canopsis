import { isArray } from 'lodash';

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
 * @property {string} bysetpos
 * @property {string} bymonthday
 * @property {string} byyearday
 * @property {string} byweekno
 * @property {string} byhour
 * @property {string} byminute
 * @property {string} bysecond
 */

/**
 * Prepare rrule options to form
 *
 * @param {RecurrenceRuleOptions} recurrenceRule
 * @return {RecurrenceRuleFormOptions}
 */
export const recurrenceRuleToFormOptions = recurrenceRule => ({
  freq: recurrenceRule.freq || '',
  count: recurrenceRule.count || '',
  interval: recurrenceRule.interval || '',
  byweekday: recurrenceRule.byweekday ? recurrenceRule.byweekday.map(v => v.weekday) : [],
  wkst: recurrenceRule.wkst ? recurrenceRule.wkst.weekday : '',
  bymonth: recurrenceRule.bymonth || [],
});

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
 * Prepare rrule options to advanced form
 *
 * @param {RecurrenceRuleOptions} recurrenceRule
 * @return {RecurrenceRuleFormAdvancedOptions}
 */
export const recurrenceRuleToFormAdvancedOptions = recurrenceRule => ({
  bysetpos: prepareRecurrenceRuleOption(recurrenceRule.bysetpos),
  bymonthday: prepareRecurrenceRuleOption(recurrenceRule.bymonthday),
  byyearday: prepareRecurrenceRuleOption(recurrenceRule.byyearday),
  byweekno: prepareRecurrenceRuleOption(recurrenceRule.byweekno),
  byhour: prepareRecurrenceRuleOption(recurrenceRule.byhour),
  byminute: prepareRecurrenceRuleOption(recurrenceRule.byminute),
  bysecond: prepareRecurrenceRuleOption(recurrenceRule.bysecond),
});
