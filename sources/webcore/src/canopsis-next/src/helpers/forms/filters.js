import { cloneDeep } from 'lodash';

import { addKeyInEntity, removeKeyFromEntity } from '@/helpers/entities';

/**
 * Convert filters to filters form
 *
 * @param {Array} [filters = []]
 * @returns {Array}
 */
export const filtersToForm = (filters = []) => cloneDeep(addKeyInEntity(filters));

/**
 * Convert filters form to filters object
 *
 * @param {Array} [filters = []]
 * @returns {Array}
 */
export const formToFilters = (filters = []) => removeKeyFromEntity(filters);
