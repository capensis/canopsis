import { isArray, keyBy, pick } from 'lodash';

/**
 * Prepare data for items by ID.
 *
 * @param {Object} [itemsById={}] - The items indexed by ID.
 * @param {Array|string} [value=[]] - The value or array of values to pick from itemsById.
 * @param {Array} [data=[]] - The data to be added to itemsById.
 * @param {string} [idKey='_id'] - The key to use as the ID for indexing data.
 * @returns {Object} The updated itemsById object.
 */
export const prepareDataForItemsById = (itemsById = {}, value = [], data = [], idKey = '_id') => ({
  ...itemsById,
  ...keyBy(data, idKey),
  ...pick(itemsById, isArray(value) ? value : [value]),
});
