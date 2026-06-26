import { isEmpty } from 'lodash';

import { DATETIME_FORMATS } from '@/constants';
import { convertDateToDateObject } from '@/helpers/date/date';

const getRuleValue = rule => (Array.isArray(rule) ? rule[0] : rule);

/**
 * Returns whether old vee-validate rules require a value.
 *
 * @param {Object} rules Legacy rules object.
 * @returns {boolean}
 */
export const isRequiredRule = rules => Boolean(rules?.required);

/**
 * Converts old date/time picker rules to field wrapper validation rules.
 *
 * @param {Object} rules Legacy rules object.
 * @param {Function} [dateObjectPreparer] Date parser.
 * @returns {Function[]}
 */
export const createDateTimeValidationRules = (
  rules = {},
  dateObjectPreparer = value => convertDateToDateObject(value, DATETIME_FORMATS.dateTimePicker),
) => {
  const validationRules = [];

  if (rules.date_format) {
    validationRules.push(value => {
      if (isEmpty(value)) {
        return true;
      }

      try {
        const date = dateObjectPreparer(value);

        return Number.isNaN(date.getTime()) ? 'Invalid date format' : true;
      } catch {
        return 'Invalid date format';
      }
    });
  }

  if (rules.after) {
    validationRules.push(value => {
      if (isEmpty(value)) {
        return true;
      }

      const afterValue = getRuleValue(rules.after);

      if (isEmpty(afterValue)) {
        return true;
      }

      try {
        const date = dateObjectPreparer(value);
        const afterDate = dateObjectPreparer(afterValue);

        return date > afterDate ? true : 'Date must be after start date';
      } catch {
        return true;
      }
    });
  }

  return validationRules;
};
