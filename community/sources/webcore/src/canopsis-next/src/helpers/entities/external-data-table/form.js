/**
 * @typedef {0 | 1} ExternalDataTableTypes
 */

import { EXTERNAL_DATA_TABLES_TYPES } from '@/constants';

/**
 * @typedef {Object} ExternalDataTable
 * @property {ExternalDataTableTypes} type
 * @property {string} name
 * @property {string} description
 */

/**
 * Converts an ExternalDataTable object to a form representation.
 *
 * @param {ExternalDataTable} [externalDataTable = {}] - The external data table object to convert.
 * @returns {ExternalDataTable}.
 */
export const externalDataTableToForm = (externalDataTable = {}) => ({
  type: externalDataTable.type ?? EXTERNAL_DATA_TABLES_TYPES.mongo,
  name: externalDataTable.name ?? '',
  description: externalDataTable.description ?? '',
});
