/**
 * @typedef {0 | 1} ExternalDataTableTypes
 */

/**
 * @typedef {0 | 1 | 2} ExternalDataTableColumnTypes
 */

import { EXTERNAL_DATA_TABLE_COLUMN_TYPES, EXTERNAL_DATA_TABLES_TYPES } from '@/constants';

/**
 * @typedef {Object} ExternalDataTable
 * @property {ExternalDataTableTypes} type
 * @property {string} name
 * @property {string} description
 * @property {string[]} [columns]
 * @property {ExternalDataTableColumnTypes[]} [columns]
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

export const externalDataTableColumnsToForm = (columns = [], columnTypes = []) => (
  columns.reduce((acc, column, index) => {
    acc[column] = columnTypes[index] ?? EXTERNAL_DATA_TABLE_COLUMN_TYPES.noType;

    return acc;
  }, {})
);
