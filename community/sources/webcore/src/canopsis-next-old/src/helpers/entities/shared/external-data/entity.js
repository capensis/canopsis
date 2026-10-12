import { EXTERNAL_DATA_TYPES } from '@/constants';

/**
 * Check external data type is table
 *
 * @param {string} type
 * @returns {boolean}
 */
export const isTableExternalDataType = type => type === EXTERNAL_DATA_TYPES.table;

/**
 * Check external data type is api
 *
 * @param {string} type
 * @returns {boolean}
 */
export const isApiExternalDataType = type => type === EXTERNAL_DATA_TYPES.api;
