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
 * @property {ExternalDataTableColumnTypes[]} [column_types]
 */

/**
 * Converts an ExternalDataTable object to a form representation.
 *
 * @param {ExternalDataTable} [externalDataTable = {}] - The external data table object to convert.
 * @returns {ExternalDataTable}
 */
export const externalDataTableToForm = (externalDataTable = {}) => ({
  type: externalDataTable.type ?? EXTERNAL_DATA_TABLES_TYPES.mongo,
  name: externalDataTable.name ?? '',
  description: externalDataTable.description ?? '',
  column_types: externalDataTable.column_types ?? [],
});

/**
 * Converts an array of columns and their corresponding types into an object mapping each column to its type.
 *
 * @param {Array} [columns=[]] - An array of column names.
 * @param {Array} [columnTypes=[]] - An array of column types corresponding to each column.
 * @returns {Object} An object where each key is a column name and its value is the column type.
 */
export const externalDataTableColumnsToForm = (columns = [], columnTypes = []) => (
  (columns ?? []).reduce((acc, column, index) => {
    acc[column] = columnTypes[index] ?? EXTERNAL_DATA_TABLE_COLUMN_TYPES.noType;

    return acc;
  }, {})
);
