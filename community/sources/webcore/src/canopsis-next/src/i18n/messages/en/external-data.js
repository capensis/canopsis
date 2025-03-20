import {
  EXTERNAL_DATA_TYPES,
  EXTERNAL_DATA_CONDITION_TYPES,
  EXTERNAL_DATA_DEFAULT_CONDITION_FIELDS,
  EXTERNAL_DATA_TABLES_TYPES,
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
    [EXTERNAL_DATA_TYPES.mongo]: 'MongoDB collection',
    [EXTERNAL_DATA_TYPES.api]: 'API',
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

  tableField: 'Collection / table',
};
