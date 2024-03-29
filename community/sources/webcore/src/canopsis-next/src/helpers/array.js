import { omit, pick, isObject, map } from 'lodash';

import { uid } from './uid';

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
 * Return entities ids
 *
 * @param {Array} entities
 * @param {string} [idKey = '_id']
 * @return Array
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
 * Create range number array
 *
 * @param {number} length
 * @param {number} [start = 0]
 * @return {number[]}
 */
export const createRangeArray = (length, start = 0) => Array.from({ length }, (_, index) => start + index);
