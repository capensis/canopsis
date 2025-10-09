import {
  EXTERNAL_DATA_TYPES,
  EXTERNAL_DATA_CONDITION_TYPES,
  EXTERNAL_DATA_DEFAULT_CONDITION_FIELDS,
  EXTERNAL_DATA_TABLES_TYPES,
  EXTERNAL_DATA_TABLE_COLUMN_TYPES,
} from '@/constants';

export default {
  title: 'External data',
  add: 'Add external data',
  empty: 'No external data added yet',
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

  tableNameTooltip: 'Supported symbols: latin letters, “_”, numbers (not at the beginning)',

  importFileDescription: 'First row has to contain column names',
  exportTableStructure: 'Export table structure',

  tableField: 'Collection / table',

  tableCanBeDeletedInConfig: 'Table can be deleted only in configuration file',
  tableCanBeDeletedAfter: 'Table can be deleted after deletion of \n{rules}',
  tableRemovedFromConfig: 'Table is removed from configuration file, but still used in the following items.\n<strong>Move table back to configuration file or delete all items that use it</strong>\n{rules}',
  tableEmptyColumns: 'Please choose at least 1 column in the settings',
  tableColumnTypes: {
    [EXTERNAL_DATA_TABLE_COLUMN_TYPES.noType]: 'No type',
    [EXTERNAL_DATA_TABLE_COLUMN_TYPES.filter]: 'Filter',
    [EXTERNAL_DATA_TABLE_COLUMN_TYPES.context]: 'Context',
  },
};
