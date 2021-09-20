import {
  convertDateToTimestampByTimezone,
  getLocaleTimezone,
  convertDateToDateObjectByTimezone,
} from '@/helpers/date/date';
import { addKeyInEntities, removeKeyFromEntities } from '@/helpers/entities';

/**
 * @typedef {Object} PbehaviorException
 * @property {string} _id
 * @property {number} created
 * @property {boolean} deletable
 * @property {string} description
 * @property {string} name
 * @property {PbehaviorExdate[]} exdates
 */

/**
 * @typedef {PbehaviorException} PbehaviorExceptionForm
 * @property {string} key
 */

/**
 * Convert pbehavior exception data to date exception form
 *
 * @param {Object} [exception = {}]
 * @param {string} [timezone = getLocaleTimezone()]
 * @return {Object}
 */
export function pbehaviorExceptionToForm(exception = {}, timezone = getLocaleTimezone()) {
  return {
    name: exception.name || '',
    description: exception.description || '',
    exdates: exception.exdates
      ? addKeyInEntities(exception.exdates.map(({ begin, end, type }) => ({
        begin: convertDateToDateObjectByTimezone(begin, timezone),
        end: convertDateToDateObjectByTimezone(end, timezone),
        type: { ...type },
      })))
      : [],
    _id: exception._id,
  };
}

/**
 * Convert exception form to pbehavior exception data
 *
 * @param {Object} [exceptionForm = {}]
 * @param {string} [timezone = getLocaleTimezone()]
 * @return {Object}
 */
export function formToPbehaviorException(exceptionForm = {}, timezone = getLocaleTimezone()) {
  const { exdates, ...form } = exceptionForm;

  return {
    exdates: removeKeyFromEntities(exdates).map(({ type, begin, end }) => ({
      type: type._id,
      begin: convertDateToTimestampByTimezone(begin, timezone),
      end: convertDateToTimestampByTimezone(end, timezone),
    })),
    ...form,
  };
}
