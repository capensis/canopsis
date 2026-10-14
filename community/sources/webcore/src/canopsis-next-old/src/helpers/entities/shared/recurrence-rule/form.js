import { RRule, rrulestr } from 'rrule';
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
  until: recurrenceRule.until,
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
    recurrenceRuleOptions.wkst = options.wkst;
  }

  if (options.until) {
    recurrenceRuleOptions.until = Date.UTC(
      options.until.getFullYear(),
      options.until.getMonth(),
      options.until.getDate(),
      options.until.getHours(),
      options.until.getMinutes(),
    );
  }

  return recurrenceRuleOptions;
};

/**
 * @return {RecurrenceRuleFormOptions}
 */
export const emptyRecurrenceRuleFormOptions = () => (
  recurrenceRuleToFormOptions(new RRule().origOptions)
);

/**
 * @param {RecurrenceRuleFormOptions} options
 * @return {string[]}
 */
export const getRecurrenceAdvancedFieldNames = (options) => {
  const { freq } = options;

  if (freq == null) {
    return [];
  }

  const fields = ['bysetpos'];

  if (freq !== RRule.MONTHLY) {
    fields.push('byyearday');
  }

  if (freq !== RRule.YEARLY) {
    fields.push('bymonthday');
  }

  if (freq !== RRule.MONTHLY && freq !== RRule.YEARLY) {
    fields.push('byweekno');
  }

  if (freq === RRule.HOURLY) {
    fields.push('byhour');
  }

  return fields;
};

/**
 * @param {string} [rruleBody]
 * @return {RecurrenceRuleFormOptions}
 */
export const rruleBodyStringToRecurrenceRuleFormOptions = (rruleBody) => {
  if (!rruleBody || typeof rruleBody !== 'string') {
    return emptyRecurrenceRuleFormOptions();
  }

  try {
    return recurrenceRuleToFormOptions(rrulestr(rruleBody).origOptions);
  } catch {
    return emptyRecurrenceRuleFormOptions();
  }
};

/**
 * API / legacy: string body, or already-normalized form options
 *
 * @param {*} value
 * @return {RecurrenceRuleFormOptions}
 */
export const normalizeRruleFormFieldForForm = (value) => {
  if (value == null || value === '') {
    return emptyRecurrenceRuleFormOptions();
  }

  if (typeof value === 'string') {
    return rruleBodyStringToRecurrenceRuleFormOptions(value);
  }

  if (typeof value === 'object' && 'freq' in value) {
    return value;
  }

  return emptyRecurrenceRuleFormOptions();
};

/**
 * @param {RecurrenceRuleFormOptions} options
 * @return {string}
 */
export const recurrenceRuleFormOptionsToRruleBodyString = (options) => {
  if (!options || options.freq == null) {
    return '';
  }

  try {
    const names = getRecurrenceAdvancedFieldNames(options);
    const rrule = new RRule(formOptionsToRecurrenceRuleOptions(options, names));

    if (!rrule.isFullyConvertibleToText()) {
      return '';
    }

    return rrule.toString().replace(/.*RRULE:/, '');
  } catch {
    return '';
  }
};

/**
 * Normalize API / config RRULE value into recurrence form state (`form.rrule`).
 *
 * @param {*} rrule
 * @return {RecurrenceRuleFormOptions}
 */
export const rruleToForm = rrule => normalizeRruleFormFieldForForm(rrule);

/**
 * Serialize `form.rrule` for the API: unchanged string (e.g. event filter before edit), or RRULE body from options.
 *
 * @param {*} formRrule
 * @return {string}
 */
export const formToRrule = (formRrule) => {
  if (formRrule == null || formRrule === '') {
    return '';
  }

  if (typeof formRrule === 'string') {
    return formRrule;
  }

  return recurrenceRuleFormOptionsToRruleBodyString(formRrule);
};

/**
 * Create-recurrence-rule modal: build internal state from `config` (API string or legacy shape).
 *
 * @param {{ rrule: *, exdates?: Array, exceptions?: Array }} modalConfig
 * @return {{ rrule: RecurrenceRuleFormOptions, exdates: Array, exceptions: Array }}
 */
export const recurrenceRuleModalConfigToForm = ({ rrule, exdates, exceptions }) => ({
  rrule: rruleToForm(rrule),
  exdates: exdates ?? [],
  exceptions: exceptions ?? [],
});

/**
 * Create-recurrence-rule modal: serialize state for `action` callback / parent merge.
 *
 * @param {{ rrule: *, exdates: Array, exceptions: Array }} modalForm
 * @return {{ rrule: string, exdates: Array, exceptions: Array }}
 */
export const formToReccurenceRuleModalConfig = ({ rrule, exdates, exceptions }) => ({
  rrule: formToRrule(rrule),
  exdates,
  exceptions,
});
