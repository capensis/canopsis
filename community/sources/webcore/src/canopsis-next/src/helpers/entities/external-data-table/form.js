import { omit } from 'lodash';

import {
  EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES,
  EXTERNAL_DATA_TABLE_COLUMN_TAGS,
  EXTERNAL_DATA_TABLES_TYPES,
  CSV_SEPARATORS_TO_SYMBOLS,
  EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_SEPARATORS,
  EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_CUSTOM_SEPARATOR,
  EXTERNAL_DATA_TABLE_PRIORITY_COLUMN,
} from '@/constants';

/**
 * @typedef {0 | 1} ExternalDataTableTypes
 */

/**
 * @typedef {0 | 1 | 2} ExternalDataTableColumnTags
 */

/**
 * @typedef {Object} ExternalDataTableBaseColumnConfig
 * @property {string} name
 * @property {string} tag
 * @property {string} type
 */

/**
 * @typedef {Object} ExternalDataTable
 * @property {ExternalDataTableTypes} type
 * @property {string} name
 * @property {string} description
 * @property {string[]} [columns]
 * @property {ExternalDataTableBaseColumnConfig[]} [column_configs]
 */

/**
 * @typedef {Object} ExternalDataTableForm
 * @property {ExternalDataTableTypes} type
 * @property {string} name
 * @property {string} description
 * @property {ExternalDataTableColumnTags[]} [column_tags]
 */

/**
 * @typedef {Object} ExternalDataTableColumnConfig
 * @property {string} name
 * @property {string} type
 * @property {string} [decimal_separator]
 * @property {string} [thousands_separator]
 * @property {string} [string_array_type]
 * @property {string} [string_array_separator]
 * @property {number[]} [rows]
 * @property {string[]} [messages]
 */

/**
 * Converts an ExternalDataTable object to a form representation.
 *
 * @param {ExternalDataTable} [externalDataTable = {}] - The external data table object to convert.
 * @returns {ExternalDataTableForm}
 */
export const externalDataTableToForm = (externalDataTable = {}) => ({
  type: externalDataTable.type ?? EXTERNAL_DATA_TABLES_TYPES.mongo,
  name: externalDataTable.name ?? '',
  description: externalDataTable.description ?? '',
  column_tags: (externalDataTable.column_configs ?? []).map(columnConfig => columnConfig.tag),
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
 * const importForm = externalDataTableColumnConfigsToForm(columnsConfigs, true);
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
 * const editForm = externalDataTableColumnConfigsToForm(columnsConfigs, false);
 * // Same structure but all type properties will be null
 */
export const externalDataTableColumnConfigsToForm = (columnsConfigs = [], isImport = false) => {
  let hasRegexpColumns = false;

  return (columnsConfigs ?? []).reduce((acc, columnConfig, index) => {
    const name = columnConfig.name ?? '';

    acc[name] = {
      name,
      rows: [],
      messages: [],
      tag: columnConfig.tag ?? EXTERNAL_DATA_TABLE_COLUMN_TAGS.noTag,
      type: isImport ? null : columnConfig.type,
      decimal_separator: columnConfig.decimal_separator ?? null,
      thousands_separator: columnConfig.thousands_separator ?? null,
      string_array_type: columnConfig.string_array_type ?? null,
      string_array_separator: columnConfig.string_array_separator ?? null,
    };

    if (columnConfig.type === EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.regexp) {
      hasRegexpColumns = true;
    }

    if (hasRegexpColumns && index === columnsConfigs.length - 1) {
      acc[EXTERNAL_DATA_TABLE_PRIORITY_COLUMN] = {
        name: EXTERNAL_DATA_TABLE_PRIORITY_COLUMN,
      };
    }
    return acc;
  }, {});
};

/**
 * Converts a form representation of external data table columns to a config representation.
 *
 * @param {Object<string, ExternalDataTableColumnConfig>} form
 * @returns {ExternalDataTableColumnConfig[]}
 */
export const formToExternalDataTableColumnConfigs = (form = {}) => Object.values(form).reduce((acc, columnConfig) => {
  if (columnConfig.name === EXTERNAL_DATA_TABLE_PRIORITY_COLUMN) {
    return acc;
  }

  acc.push({
    ...omit(columnConfig, ['rows', 'messages']),

    type: columnConfig.type ?? EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.string,
  });

  return acc;
}, []);

/**
 * Converts a form representation of external data table columns to a tags representation.
 *
 * @param {Object<string, ExternalDataTableColumnConfig>} form
 * @returns {number[]}
 */
export const formToExternalDataTableColumnTags = (form = {}) => Object.values(form).reduce((acc, columnConfig) => {
  if (columnConfig.name === EXTERNAL_DATA_TABLE_PRIORITY_COLUMN) {
    return acc;
  }

  acc.push(columnConfig.tag);

  return acc;
}, []);

/**
 * Checks if a separator is a predefined standard separator
 * @param {string} separator - The separator to check
 * @returns {boolean} True if the separator is predefined, false otherwise
 */
const isStandardSeparator = separator => Object.values(EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_SEPARATORS)
  .includes(separator);

/**
 * Gets the fallback separator when no separator is defined
 * Avoids using the same separator as the table's CSV separator
 * @param {string} tableSeparator - The table's CSV separator
 * @returns {string} A safe fallback separator
 */
const getFallbackSeparator = (tableSeparator) => {
  const tableSymbol = CSV_SEPARATORS_TO_SYMBOLS[tableSeparator];
  const { comma, semicolon } = EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_SEPARATORS;

  return tableSymbol !== comma ? comma : semicolon;
};

/**
 * Gets the default separator value based on the current value and table separator
 *
 * Logic:
 * 1. If value has a custom separator (not in predefined list) → return 'custom'
 * 2. If value has no separator → return a safe fallback separator
 * 3. If value has a standard separator → return it as-is
 *
 * @param {Object} value - The current column data type value
 * @param {string} tableSeparator - The table's CSV separator
 * @returns {string} The default separator value or 'custom' if custom separator is needed
 */
export const getDefaultSeparator = (value, tableSeparator) => {
  const { string_array_separator: currentSeparator } = value;

  if (currentSeparator && !isStandardSeparator(currentSeparator)) {
    return EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_CUSTOM_SEPARATOR;
  }

  if (!currentSeparator) {
    return getFallbackSeparator(tableSeparator);
  }

  return currentSeparator;
};

/**
 * Adds priority column to columns array when regexp columns are detected.
 *
 * @param {Array} columns - Array of column objects with type property
 * @returns {Array} Columns array with priority column added if needed
 */
export const addPriorityColumnToColumnsArray = (columns = []) => {
  const hasRegexpColumns = columns.some(column => column.type === EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.regexp);

  if (!hasRegexpColumns) {
    return columns;
  }

  return [
    ...columns,
    {
      text: EXTERNAL_DATA_TABLE_PRIORITY_COLUMN,
      value: EXTERNAL_DATA_TABLE_PRIORITY_COLUMN,
      name: EXTERNAL_DATA_TABLE_PRIORITY_COLUMN,
    },
  ];
};

/**
 * Adds priority column to form when regexp columns are detected.
 *
 * @param {Object<string, ExternalDataTableColumnConfig>} form - The form object to process
 * @returns {Object<string, ExternalDataTableColumnConfig>} Form with priority column added if needed
 */
export const addPriorityColumnForRegexpTypes = (form) => {
  const hasRegexpColumns = Object.values(form)
    .some(column => column.type === EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.regexp);

  if (!hasRegexpColumns) {
    return form;
  }

  return {
    ...form,

    [EXTERNAL_DATA_TABLE_PRIORITY_COLUMN]: {
      name: EXTERNAL_DATA_TABLE_PRIORITY_COLUMN,
    },
  };
};
