/**
 * @typedef {0 | 1} ExternalDataTableTypes
 */

/**
 * @typedef {0 | 1 | 2} ExternalDataTableColumnTypes
 */

import { EXTERNAL_DATA_TABLE_COLUMN_TAGS, EXTERNAL_DATA_TABLES_TYPES } from '@/constants';

/**
 * @typedef {Object} ExternalDataTable
 * @property {ExternalDataTableTypes} type
 * @property {string} name
 * @property {string} description
 * @property {string[]} [columns]
 * @property {ExternalDataTableColumnTypes[]} [column_types]
 */

/**
 * @typedef {Object} ExternalDataTableColumnConfig
 * @property {string} name
 * @property {string} type
 * @property {string} [decimal_separator]
 * @property {string} [thousands_separator]
 * @property {string} [string_array_type]
 * @property {string} [string_array_separator]
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
 * @param {Array<ExternalDataTableColumnConfig>} [columnsConfigs=[]] - An array of column configs.
 * @param {boolean} [isImport=false] - Whether the conversion is for import mode. When true,
 *                                     preserves the original type, when false sets type to null.
 * @returns {Object<string, ExternalDataTableColumnConfig>} An object where each key is a column name
 *                                                          and its value is the column config with all
 *                                                          formatting properties.
 * @example
 * // Basic usage with string and number columns
 * const columnsConfigs = [
 *   { name: 'user_name', type: 'string' },
 *   { name: 'age', type: 'number', decimal_separator: '.', thousands_separator: ',' },
 *   { name: 'tags', type: 'string_array', string_array_separator: '|', string_array_type: 'set' }
 * ];
 *
 * // For import mode (preserves types)
 * const importForm = externalDataTableColumnsConfigToForm(columnsConfigs, true);
 * // Result:
 * // {
 * //   user_name: {
 * //     name: 'user_name',
 * //     type: 'string',
 * //     tag: 'no-type',
 * //     decimal_separator: null,
 * //     thousands_separator: null,
 * //     string_array_type: null,
 * //     string_array_separator: null
 * //   },
 * //   age: {
 * //     name: 'age',
 * //     type: 'number',
 * //     tag: 'no-type',
 * //     decimal_separator: '.',
 * //     thousands_separator: ',',
 * //     string_array_type: null,
 * //     string_array_separator: null
 * //   },
 * //   tags: {
 * //     name: 'tags',
 * //     type: 'string_array',
 * //     tag: 'no-type',
 * //     decimal_separator: null,
 * //     thousands_separator: null,
 * //     string_array_type: 'set',
 * //     string_array_separator: '|'
 * //   }
 * // }
 *
 * // For edit mode (resets types to null)
 * const editForm = externalDataTableColumnsConfigToForm(columnsConfigs, false);
 * // Same structure but all type properties will be null
 */
export const externalDataTableColumnsConfigToForm = (columnsConfigs = [], isImport = false) => (
  (columnsConfigs ?? []).reduce((acc, columnConfig) => {
    const name = columnConfig.name ?? '';

    acc[name] = {
      name,
      tag: columnConfig.tag ?? EXTERNAL_DATA_TABLE_COLUMN_TAGS.noType,
      type: isImport ? null : columnConfig.type,
      decimal_separator: columnConfig.decimal_separator ?? null,
      thousands_separator: columnConfig.thousands_separator ?? null,
      string_array_type: columnConfig.string_array_type ?? null,
      string_array_separator: columnConfig.string_array_separator ?? null,
    };

    return acc;
  }, {})
);
