import { omit, isObject, groupBy } from 'lodash';

import i18n from '@/i18n';
import {
  ENTITIES_STATES,
  ENTITIES_STATUSES,
  DATETIME_FORMATS,
} from '@/constants';

import { convertDateToString } from '@/helpers/date/date';

import uid from './uid';

/**
 * @typedef {Object} Interval
 * @property {number} interval
 * @property {string} unit
 */

/**
 * Convert default columns from constants to columns with prepared by i18n label
 *
 * @param {{ labelKey: string, value: string }[]} [columns = []]
 * @returns {{ label: string, value: string }[]}
 */
export function defaultColumnsToColumns(columns = []) {
  return columns.map(({ labelKey, value }) => ({
    label: i18n.t(labelKey),
    value,
  }));
}

/**
 * Checks if alarm is resolved
 * @param alarm - alarm entity
 * @returns {boolean}
 */
export const isResolvedAlarm = alarm => [ENTITIES_STATUSES.closed, ENTITIES_STATUSES.cancelled]
  .includes(alarm.v.status.val);

/**
 * Checks if alarm have critical state
 *
 * @param alarm - alarm entity
 * @returns {boolean}
 */
export const isWarningAlarmState = alarm => ENTITIES_STATES.ok !== alarm.v.state.val;

/**
 * Add uniq key field in each entity.
 *
 * @param {Array} entities
 * @return {Array}
 */
export const addKeyInEntities = (entities = []) => entities.map(entity => ({
  ...entity,
  key: uid(),
}));

/**
 * Remove key field from each entity.
 *
 * @param {Array} entities
 * @return {Array}
 */
export const removeKeyFromEntities = (entities = []) => entities.map(entity => omit(entity, ['key']));

/**
 * Get id from entity
 *
 * @param {Object} entity
 * @param {string} idField
 * @return {string}
 */
export const getIdFromEntity = (entity, idField = '_id') => (isObject(entity) ? entity[idField] : entity);

/**
 * Get grouped steps by date
 *
 * @param {AlarmEvent[]} steps
 * @return {Object.<string, AlarmEvent[]>}
 */
export const groupAlarmSteps = (steps) => {
  const orderedSteps = [...steps].reverse();

  return groupBy(orderedSteps, step => convertDateToString(step.t, DATETIME_FORMATS.short));
};
