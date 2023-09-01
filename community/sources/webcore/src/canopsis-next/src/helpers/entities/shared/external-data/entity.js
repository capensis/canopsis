import { EXTERNAL_DATA_TYPES } from '@/constants';

/**
 * Check external data type is mongo
 *
 * @param {string} type
 * @returns {boolean}
 */
export const isMongoExternalDataType = type => type === EXTERNAL_DATA_TYPES.mongo;

/**
 * Check external data type is api
 *
 * @param {string} type
 * @returns {boolean}
 */
export const isApiExternalDataType = type => type === EXTERNAL_DATA_TYPES.api;
