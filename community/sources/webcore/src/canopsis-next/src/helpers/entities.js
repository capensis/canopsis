import { omit, pick, isObject, groupBy, map } from 'lodash';

import {
  ENTITIES_STATES,
  ENTITIES_STATUSES,
  DATETIME_FORMATS,
} from '@/constants';

import { uid } from './uid';
import { convertDateToString } from './date/date';

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
 * Add uniq key field in entity object.
 *
 * @param {Object} entity
 * @return {Object}
 */
export const addKeyInEntity = (entity = {}) => ({ ...entity, key: uid() });

/**
 * Add uniq key field in each entity.
 *
 * @param {Array} entities
 * @return {Array}
 */
export const addKeyInEntities = (entities = []) => entities.map(addKeyInEntity);

/**
 * Remove key field from each entity.
 *
 * @param {Object} entity
 * @return {Object}
 */
export const removeKeyFromEntity = (entity = {}) => omit(entity, ['key']);

/**
 * Remove key field from each entity.
 *
 * @param {Array} entities
 * @return {Array}
 */
export const removeKeyFromEntities = (entities = []) => entities.map(removeKeyFromEntity);

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
export const groupAlarmSteps = steps => (
  groupBy(steps, step => convertDateToString(step.t, DATETIME_FORMATS.short))
);

/**
 * Return entities ids
 *
 * @param {Array} entities
 * @param {string} [idKey = '_id']
 */
export const mapIds = (entities, idKey = '_id') => map(entities, idKey);

/**
 * Pick id field from entities
 *
 * @param {U[]} entities
 * @param {string} idKey
 * @returns {PartialObject<U>[]}
 */
export const pickIds = (entities = [], idKey = '_id') => entities.map(entity => pick(entity, [idKey]));

/**
 * Filter entities by ids
 *
 * @param {Object[]} items
 * @param {Object} item
 * @param {string} [idKey = '_id']
 */
export const filterById = (items, item, idKey = '_id') => items
  .filter(({ [idKey]: itemId }) => item[idKey] !== itemId);

/**
 * Filter entities by value
 *
 * @param {string[] | number[]} items
 * @param {string | number} removingValue
 */
export const filterValue = (items, removingValue) => items.filter(item => item !== removingValue);

/**
 * Revert grouped values by key
 *
 * @example
 *  revertGroupBy({
 *    'key1': ['value1', 'value2'],
 *    'key2': ['value1', 'value2', 'value3'],
 *    'key3': ['value3'],
 *  }) -> {
 *    'value1': ['key1', 'key2'],
 *    'value2': ['key1', 'key2'],
 *    'value3': ['key2', 'key3'],
 *  }
 * @param {Object.<string, string[]>} obj
 * @returns {Object.<string, string[]>}
 */
export const revertGroupBy = obj => Object.entries(obj).reduce((acc, [id, values]) => {
  values.forEach((value) => {
    if (acc[value]) {
      acc[value].push(id);
    } else {
      acc[value] = [id];
    }
  });

  return acc;
}, {});

/**
 * Generate alarm details id by widgetId
 *
 * @param {string} alarmId
 * @param {string} widgetId
 * @returns {string}
 */
export const generateAlarmDetailsId = (alarmId, widgetId) => `${alarmId}_${widgetId}`;

/**
 * Get dataPreparer for alarmDetails entity
 *
 * @param {string} widgetId
 * @returns {Function}
 */
export const getAlarmDetailsDataPreparer = widgetId => data => (
  data.map(item => ({
    ...item,

    /**
     * We are generating new id based on alarmId and widgetId to avoiding collision with two widgets
     * on the same view with opened expand panel on the same alarm
     */
    _id: generateAlarmDetailsId(item._id, widgetId),
  }))
);
