import moment from 'moment';
import { cloneDeep } from 'lodash';

import { addKeyInEntity, removeKeyFromEntity } from '@/helpers/entities';

/**
 * Convert pbehavior date exception data to date exception form
 *
 * @param {Object} dateException
 * @return {Object}
 */
export function pbehaviorDateExceptionToForm(dateException = {}) {
  return {
    name: dateException.name || '',
    description: dateException.description || '',
    exdates: dateException.exdates ? addKeyInEntity(cloneDeep(dateException.exdates)) : [],
  };
}

/**
 * Convert date exception form to pbehavior date exception data
 *
 * @param {Object} dateExceptionForm
 * @return {Object}
 */
export function formToPbehaviorDateException(dateExceptionForm = {}) {
  const { exdates, ...form } = dateExceptionForm;

  return {
    exdates: removeKeyFromEntity(exdates).map(({ type, begin, end }) => ({
      type: type._id,
      begin: moment(begin).unix(),
      end: moment(end).unix(),
    })),
    ...form,
  };
}
