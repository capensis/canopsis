import { map, omit, pick } from 'lodash'
import { isArray } from 'lodash'

import { uid } from './uid'

export const addKeyInEntity = (entity = {}) => ({ ...entity, key: uid() })

export const addKeyInEntities = (entities = []) => entities.map(addKeyInEntity)

export const removeKeyFromEntity = (entity = {}) => omit(entity, ['key'])

export const removeKeyFromEntities = (entities = []) => entities.map(removeKeyFromEntity)

export const getIdFromEntity = (entity, idField = '_id') => (
  typeof entity === 'object' && entity !== null ? entity[idField] : entity
)

export const mapIds = (entities, idKey = '_id') => map(entities, idKey)

export const pickIds = (entities = [], idKey = '_id') => entities.map(entity => pick(entity, [idKey]))

export const filterById = (items, item, idKey = '_id') => (
  items.filter(({ [idKey]: itemId }) => item[idKey] !== itemId)
)

export const filterValue = (items, removingValue) => items.filter(item => item !== removingValue)

export const createRangeArray = (length, start = 0) => Array.from({ length }, (_, index) => start + index)

/**
 * Recursively transforms an array into an object keyed by a specified property.
 *
 * @param {Array<Object>} array - The array to transform.
 * @param {string} [idKey = '_id'] - The key to use for object keys.
 * @param {string} [childrenKey = 'items'] - The key for nested children arrays.
 * @param {string} [prefix = ''] - Prefix for the key correctness
 * @returns {Object} An object keyed by the specified idKey.
 */
export const deepKeyBy = (array, idKey = '_id', childrenKey = 'items', prefix = '') => array.reduce((acc, item) => {
  const key = prefix && !String(item[idKey]).startsWith(prefix) ? `${prefix}${item[idKey]}` : item[idKey]

  acc[key] = item

  if (isArray(item[childrenKey])) {
    Object.assign(acc, deepKeyBy(item[childrenKey], idKey, childrenKey, `${key}.`))
  }

  return acc
}, {})

/**
 * Creates an object mapping each array element key to its index.
 *
 * @param {Array<Object>} items
 * @param {string} [key = 'value']
 * @returns {Object}
 */
export const indexesByKey = (items = [], key = 'value') => items.reduce((acc, item, index) => {
  acc[item[key]] = index

  return acc
}, {})
