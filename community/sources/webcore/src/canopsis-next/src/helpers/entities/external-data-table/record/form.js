import { isArray } from 'lodash';

import { EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES } from '@/constants';

import { convertDateToTimestamp } from '@/helpers/date/date';
import { primitiveArrayToForm, formToPrimitiveArray } from '@/helpers/entities/shared/form';

/**
 * Converts an external data table record to a form representation.
 *
 * @param {Object} [externalDataTableRecord = {}] - The external data table record object with column values.
 * @param {Array<{name: string, type: number}>} [columnConfigs = []] - Array of column configuration objects.
 * @returns {Object<string, string|string[]>} Form object with column names as keys and form values as values.
 *                                           String arrays are cloned for stringArray type columns.
 */
export const externalDataTableRecordToForm = (
  externalDataTableRecord = {},
  columnConfigs = [],
) => columnConfigs.reduce((acc, column) => {
  acc[column.name] = externalDataTableRecord[column.name] ?? '';

  if (column.type === EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.stringArray) {
    const stringArray = acc[column.name];

    acc[column.name] = (isArray(stringArray) ? primitiveArrayToForm(stringArray) : stringArray) || [];
  } else if (column.type === EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.boolean) {
    acc[column.name] = !!acc[column.name];
  }

  if (!acc[column.name] && [
    EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.timestamp,
    EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.datetime,
  ].includes(column.type)) {
    acc[column.name] = null;
  }

  return acc;
}, {});

/**
 * Converts a form representation to an external data table record.
 *
 * @param {Object} [form = {}] - The form object with column names as keys and form values as values.
 * @param {Array<{name: string, type: number}>} [columnConfigs = []] - Array of column configuration objects.
 * @returns {Object<string, *>} External data table record object with column names as keys.
 *                              Timestamp and datetime values are converted to Unix timestamps.
 */
export const formToExternalDataTableRecord = (form = {}, columnConfigs = []) => columnConfigs.reduce((acc, column) => {
  let value = form[column.name];

  if (column.type === EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.stringArray) {
    value = isArray(value) ? formToPrimitiveArray(value) : value;
  }

  acc[column.name] = value;

  if ([
    EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.timestamp,
    EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.datetime,
  ].includes(column.type)) {
    acc[column.name] = convertDateToTimestamp(acc[column.name]);
  }

  return acc;
}, {});
