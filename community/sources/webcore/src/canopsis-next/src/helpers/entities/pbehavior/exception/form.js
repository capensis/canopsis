import {
  convertDateToTimestampByTimezone,
  getLocalTimezone,
  convertDateToDateObjectByTimezone,
  isStartOfDay,
} from '@/helpers/date/date';
import { addKeyInEntities, removeKeyFromEntities } from '@/helpers/array';

/**
 * @typedef {Object} PbehaviorException
 * @property {string} _id
 * @property {number} created
 * @property {boolean} deletable
 * @property {string} description
 * @property {string} name
 * @property {boolean} hidden
 * @property {PbehaviorExdate[]} exdates
 */

/**
 * @typedef {ObjectKey} PbehaviorExceptionForm
 * @property {string} _id
 * @property {number} created
 * @property {boolean} deletable
 * @property {string} description
 * @property {string} name
 * @property {boolean} visible
 * @property {PbehaviorExdate[]} exdates
 */

/**
 * Convert pbehavior exception data to date exception form
 *
 * @param {PbehaviorException} [exception = {}]
 * @param {string} [timezone = getLocalTimezone()]
 * @return {PbehaviorExceptionForm}
 */
export function pbehaviorExceptionToForm(exception = {}, timezone = getLocalTimezone()) {
  return {
    name: exception.name ?? '',
    description: exception.description ?? '',
    visible: !(exception.hidden ?? false),
    exdates: exception.exdates
      ? addKeyInEntities(exception.exdates.map(({ begin, end, type }) => ({
        begin: convertDateToDateObjectByTimezone(begin, timezone),
        end: convertDateToDateObjectByTimezone(isStartOfDay(end) ? end - 1 : end, timezone),
        type: { ...type },
      })))
      : [],
    _id: exception._id,
  };
}

/**
 * Convert pbehavior exception data to date exception form
 *
 * @param {Object} [exceptionImport = {}]
 * @return {Object}
 */
export function pbehaviorExceptionImportToForm(exceptionImport = {}) {
  return {
    name: exceptionImport.name ?? '',
    type: exceptionImport.type ?? '',
  };
}

/**
 * Convert exception form to pbehavior exception data
 *
 * @param {PbehaviorExceptionForm} [exceptionForm = {}]
 * @param {string} [timezone = getLocalTimezone()]
 * @return {PbehaviorException}
 */
export function formToPbehaviorException(exceptionForm = {}, timezone = getLocalTimezone()) {
  const { exdates, visible, ...form } = exceptionForm;

  return {
    exdates: removeKeyFromEntities(exdates).map(({ type, begin, end }) => ({
      type: type._id,
      begin: convertDateToTimestampByTimezone(begin, timezone),
      end: convertDateToTimestampByTimezone(end, timezone),
    })),
    hidden: !visible,
    ...form,
  };
}
