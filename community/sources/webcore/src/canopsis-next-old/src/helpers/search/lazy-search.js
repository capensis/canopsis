import { isArray, keyBy, pick } from 'lodash';

/**
 * Prepare data for items by ID.
 *
 * @param {Object} [itemsById = {}] - The items indexed by ID.
 * @param {Array|String} [value = []] - The value or array of values to pick from itemsById.
 * @param {Object} [data = {}]- The data to be added to itemsById.
 * @param {String} [idKey = '_id'] - The key to use as the ID for indexing data.
 * @returns {Object} The updated itemsById object.
 */
export const prepareDataForItemsById = (itemsById = {}, value = [], data = {}, idKey = '_id') => ({
  ...itemsById,
  ...keyBy(data, idKey),
  ...pick(itemsById, isArray(value) ? value : [value]),
});
