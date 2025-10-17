import {
  EXTERNAL_DATA_TYPES,
  EXTERNAL_DATA_CONDITION_TYPES,
  EXTERNAL_DATA_DEFAULT_CONDITION_FIELDS,
  EXTERNAL_DATA_TABLES_TYPES,
  EXTERNAL_DATA_TABLE_COLUMN_TAGS,
  EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES,
  EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_TYPES,
} from '@/constants';

export default {
  title: 'External data',
  add: 'Add external data',
  empty: 'No external data added yet',
  updatePreview: 'Update preview',
  loadingPreview: 'Loading preview',
  fields: {
    reference: 'Reference',
    collection: 'Collection',
    sort: 'Sort',
    sortBy: 'Sort by',
  },
  tooltips: {
    reference: 'Will be used in actions as <strong>.ExternalData.&lt;Reference&gt;</strong>',
  },
  types: {
    [EXTERNAL_DATA_TYPES.api]: 'API',
    [EXTERNAL_DATA_TYPES.table]: 'Table',
  },
  conditionTypes: {
    [EXTERNAL_DATA_CONDITION_TYPES.select]: 'Select',
    [EXTERNAL_DATA_CONDITION_TYPES.regexp]: 'Regexp',
  },
  conditionValues: {
    [EXTERNAL_DATA_DEFAULT_CONDITION_FIELDS.component]: 'Component',
    [EXTERNAL_DATA_DEFAULT_CONDITION_FIELDS.connector]: 'Connector',
    [EXTERNAL_DATA_DEFAULT_CONDITION_FIELDS.connectorName]: 'Connector name',
    [EXTERNAL_DATA_DEFAULT_CONDITION_FIELDS.resource]: 'Resource',
    [EXTERNAL_DATA_DEFAULT_CONDITION_FIELDS.output]: 'Output',
    [EXTERNAL_DATA_DEFAULT_CONDITION_FIELDS.extraInfos]: 'Extra infos',
  },

  tableTypes: {
    [EXTERNAL_DATA_TABLES_TYPES.mongo]: 'MongoDB',
    [EXTERNAL_DATA_TABLES_TYPES.postgres]: 'PostgreSQL',
  },

  tableNameTooltip: 'Supported symbols: latin letters, "_", numbers (not at the beginning)',

  importFileDescription: 'First row has to contain column names',
  exportTableStructure: 'Export table structure',

  tableField: 'Collection / table',

  andMore: 'and more...',
  fieldsHasError: '{count} field has errors|{count} fields have errors',
  linkedRules: {
    widgets: '<strong>Widgets</strong> that uses this table<br><ul>{rules}</ul>',
    eventFilters: '<strong>Event filters</strong>\n<ul>{rules}</ul>',
    links: '<strong>Links</strong>\n<ul>{rules}</ul>',
  },
  tableCanBeDeletedInConfig: 'Table can be deleted only in configuration file',
  tableCanBeDeletedAfter: 'Table can be deleted after deletion of \n{rules}',
  tableRemovedFromConfig: 'Table is removed from configuration file, but still used in the following items.\n<strong>Move table back to configuration file or delete all items that use it</strong>\n{rules}',
  tableEmptyColumns: 'Please choose at least 1 column in the settings',
  tableColumnTypes: {
    [EXTERNAL_DATA_TABLE_COLUMN_TAGS.noType]: 'No type',
    [EXTERNAL_DATA_TABLE_COLUMN_TAGS.filter]: 'Filter',
    [EXTERNAL_DATA_TABLE_COLUMN_TAGS.context]: 'Context',
  },
  tableColumnDataTypes: {
    [EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.string]: {
      text: '@:common.variableTypes.string',
    },
    [EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.number]: {
      text: '@:common.variableTypes.number',
      tooltip: 'Values in field must contain only:<br>• numbers<br>• commas<br>• dots<br>• minus signs<br>• space',
    },
    [EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.boolean]: {
      text: '@:common.variableTypes.boolean',
      tooltip: 'Values in field must contain any of the following values (case insensitive):<br>• yes / no<br>• y / n<br>• oui / non<br>• true / false<br>• 1 / 0<br><br>All these values will be converted to true / false',
    },
    [EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.timestamp]: {
      text: '@:common.timestamp',
    },
    [EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.datetime]: {
      text: 'Datetime',
      tooltip: 'Values in field must be in any of the following formats:<br>• 1990-12-31T00:00:00.000Z<br>• 1990-12-31T00:00:00:00:00<br>• 1990-12-31T00:00:00Z<br>• 1990-12-31T00:00:00+00:00',
    },
    [EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.stringArray]: {
      text: '@:common.variableTypes.array',
    },
  },
  selectDataType: 'Select data type',
  tableColumnDataTypesAdditionalChips: {
    number: {
      selectDecimalSeparator: 'Select decimal separator',
      selectThousandsSeparator: 'Select thousands separator',
      decimalSeparator: 'Decimal separator',
      thousandsSeparator: 'Thousands separator',
      decimalSeparatorDisabled: 'Is already used as thousands separator',
      thousandsSeparatorDisabled: 'Is already used as decimal separator',
      separatorDisabledByTableSeparator: 'Is already used as table separator',
    },
    stringArray: {
      separator: 'Separator',
      selectSeparator: 'Select separator',
      types: {
        [EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_TYPES.json]: {
          text: 'JSON',
          description: '[v1, v2, v3]',
        },
        [EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_TYPES.custom]: {
          text: 'Parse using separator',
          description: 'v1,v2,v3',
        },
      },
    },
    forbiddenSeparator: 'This separator cannot be used as it conflicts with the table separator',
  },
};
