import {
  convertDateToTimestampByTimezone,
  getLocalTimezone, convertTimestampToDateObjectByTimezone,
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
 * @param {string} [timezone = getLocalTimezone()]
 * @return {Object}
 */
export function pbehaviorExceptionToForm(exception = {}, timezone = getLocalTimezone()) {
  return {
    name: exception.name || '',
    description: exception.description || '',
    exdates: exception.exdates
      ? addKeyInEntities(exception.exdates.map(({ begin, end, type }) => ({
        begin: convertTimestampToDateObjectByTimezone(begin, timezone),
        end: convertTimestampToDateObjectByTimezone(end, timezone),
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
 * @param {string} [timezone = getLocalTimezone()]
 * @return {Object}
 */
export function formToPbehaviorException(exceptionForm = {}, timezone = getLocalTimezone()) {
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
