import moment from 'moment-timezone';

import { convertTimestampToMomentByTimezone, convertDateToTimestampByTimezone } from '@/helpers/date/date';
import { addKeyInEntity, removeKeyFromEntity } from '@/helpers/entities';

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
 * @param {string} [timezone = moment.tz.guess()]
 * @return {Object}
 */
export function pbehaviorExceptionToForm(exception = {}, timezone = moment.tz.guess()) {
  return {
    name: exception.name || '',
    description: exception.description || '',
    exdates: exception.exdates
      ? addKeyInEntity(exception.exdates.map(({ begin, end, type }) => ({
        begin: convertTimestampToMomentByTimezone(begin, timezone).toDate(),
        end: convertTimestampToMomentByTimezone(end, timezone).toDate(),
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
 * @param {string} [timezone = moment.tz.guess()]
 * @return {Object}
 */
export function formToPbehaviorException(exceptionForm = {}, timezone = moment.tz.guess()) {
  const { exdates, ...form } = exceptionForm;

  return {
    exdates: removeKeyFromEntity(exdates).map(({ type, begin, end }) => ({
      type: type._id,
      begin: convertDateToTimestampByTimezone(begin, timezone),
      end: convertDateToTimestampByTimezone(end, timezone),
    })),
    ...form,
  };
}
