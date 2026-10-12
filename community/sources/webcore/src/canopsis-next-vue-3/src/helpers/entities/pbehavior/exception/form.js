import { addKeyInEntities, removeKeyFromEntities } from '@/helpers/array';
import {
  convertDateToDateObjectByTimezone,
  convertDateToTimestampByTimezone,
  getLocalTimezone,
  isStartOfDay,
} from '@/helpers/date/date';

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
export const pbehaviorExceptionToForm = (exception = {}, timezone = getLocalTimezone()) => ({
  name: exception.name ?? '',
  description: exception.description ?? '',
  visible: !(exception.hidden ?? false),
  exdates: exception.exdates
    ? addKeyInEntities(
        exception.exdates.map(({ begin, end, type }) => ({
          begin: convertDateToDateObjectByTimezone(begin, timezone),
          end: convertDateToDateObjectByTimezone(isStartOfDay(end) ? end - 1 : end, timezone),
          type: { ...type },
        })),
      )
    : [],
  _id: exception._id,
});

/**
 * Convert pbehavior exception data to date exception form
 *
 * @param {Object} [exceptionImport = {}]
 * @return {Object}
 */
export const pbehaviorExceptionImportToForm = (exceptionImport = {}) => ({
  name: exceptionImport.name ?? '',
  type: exceptionImport.type ?? '',
  file: exceptionImport.file,
});

/**
 * Convert pbehavior exception import form to request data.
 *
 * @param {Object} [exceptionImportForm = {}]
 * @return {Object}
 */
export const formToPbehaviorExceptionImport = (exceptionImportForm = {}) => ({
  name: exceptionImportForm.name,
  type: exceptionImportForm.type?._id ?? exceptionImportForm.type,
  file: exceptionImportForm.file,
});

/**
 * Convert exception form to pbehavior exception data
 *
 * @param {PbehaviorExceptionForm} [exceptionForm = {}]
 * @param {string} [timezone = getLocalTimezone()]
 * @return {PbehaviorException}
 */
export const formToPbehaviorException = (exceptionForm = {}, timezone = getLocalTimezone()) => {
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
};
