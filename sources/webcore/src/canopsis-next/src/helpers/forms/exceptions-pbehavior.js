import moment from 'moment';

import { addKeyInEntity, removeKeyFromEntity } from '@/helpers/entities';

/**
 * Convert pbehavior exception data to date exception form
 *
 * @param {Object} exception
 * @return {Object}
 */
export function pbehaviorExceptionToForm(exception = {}) {
  return {
    name: exception.name || '',
    description: exception.description || '',
    exdates: exception.exdates
      ? addKeyInEntity(exception.exdates.map(({ begin, end, type }) => ({
        begin: moment.unix(begin).toDate(),
        end: moment.unix(end).toDate(),
        type: { ...type },
      })))
      : [],
  };
}

/**
 * Convert exception form to pbehavior exception data
 *
 * @param {Object} exceptionForm
 * @return {Object}
 */
export function formToPbehaviorException(exceptionForm = {}) {
  const { exdates, ...form } = exceptionForm;

  return {
    exdates: removeKeyFromEntity(exdates).map(({ type, begin, end }) => ({
      type: type._id,
      begin: moment(begin).unix(),
      end: moment(end).unix(),
    })),
    ...form,
  };
}
