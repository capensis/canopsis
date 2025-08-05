import { EVENT_FILTER_TYPES, EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES, EVENT_FILTER_FAILURE_TYPES } from '@/constants';

export default {
  externalData: 'External data',
  actionsRequired: 'Please add at least one action',
  configRequired: 'No configuration defined. Please add at least one config parameter',
  idHelp: 'If no id is specified, a unique id will be generated automatically on rule creation',
  editPattern: 'Edit pattern',
  advanced: 'Advanced',
  addAField: 'Add a field',
  simpleEditor: 'Simple editor',
  field: 'Field',
  value: 'Value',
  advancedEditor: 'Advanced editor',
  comparisonRules: 'Comparison rules',
  editActions: 'Edit actions',
  addAction: 'Add an action',
  editAction: 'Edit an action',
  actions: 'Actions',
  onSuccess: 'On success',
  onFailure: 'On failure',
  configuration: 'Configuration',
  resource: 'Resource ID or template',
  component: 'Component ID or template',
  connector: 'Connector ID or template',
  connectorName: 'Connector name or template',
  duringPeriod: 'Applied during this period only',
  enrichmentOptions: 'Enrichment options',
  changeEntityOptions: 'Change entity options',
  eventsFilteredSinceLastUpdate: 'Events filtered since last update',
  errorsSinceLastUpdate: 'Errors since last update',
  markAsRead: 'Mark as read',
  filterByType: 'Filter by type',
  copyEventToClipboard: 'Copy event to clipboard',
  event: 'Event',
  eventCopied: 'Event copied to clipboard',
  syntaxIsValid: 'Syntax is valid',
  types: {
    [EVENT_FILTER_TYPES.drop]: 'Drop',
    [EVENT_FILTER_TYPES.break]: 'Break',
    [EVENT_FILTER_TYPES.enrichment]: 'Enrichment',
    [EVENT_FILTER_TYPES.changeEntity]: 'Change entity',
  },
  failureTypes: {
    [EVENT_FILTER_FAILURE_TYPES.invalidPattern]: 'Invalid pattern',
    [EVENT_FILTER_FAILURE_TYPES.invalidTemplate]: 'Invalid template',
    [EVENT_FILTER_FAILURE_TYPES.externalDataMongo]: 'Mongo',
    [EVENT_FILTER_FAILURE_TYPES.externalDataApi]: 'External API',
    [EVENT_FILTER_FAILURE_TYPES.other]: 'Other',
  },
  tooltips: {
    addValueRuleField: 'Add value rule field',
    editValueRuleField: 'Edit value rule field',
    addObjectRuleField: 'Add object rule field',
    editObjectRuleField: 'Edit object rule field',
    removeRuleField: 'Remove rule field',
  },
  validation: {
    incorrectRegexOnSetTagsValue: 'Invalid value: the value for the set_tags action must contain regex to extract groups <name> and <value>',
  },
  actionsTypes: {
    [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.copy]: {
      text: 'Copy a value from a field of event to another',
      message: 'This action is used used to copy the value or a pair key+value of a control in an event.',
      description: '<h3 class="text-subtitle-1 font-weight-bold">Action parameters</h3>'
      + '<ul>'
        + '<li><strong>value</strong>: the name of the control whose value must be copied. It can be an event field, a subgroup of a regular expression, or an external data</li>'
        + '<li><strong>description</strong> (optional): the description</li>'
        + '<li><strong>name</strong>: the name of the event field into which the value must be copied</li>'
      + '</ul>',
    },
    [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.copyToEntityInfo]: {
      text: 'Copy a value from a field of an event to an info of an entity',
      message: 'This action is used to copy the field value of an event to the field of an entity.',
      description: '<h3 class="text-subtitle-1 font-weight-bold">Action parameters</h3>'
      + '<ul>'
        + '<li><strong>description</strong> (optional): the description</li>'
        + '<li><strong>name</strong>: the name of the field of an entity</li>'
        + '<li><strong>value</strong>: the name of the control whose value must be copied. It can be an event field, a subgroup of a regular expression, or an external data</li>'
      + '</ul>',
    },
    [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfo]: {
      text: 'Set an info of an entity to a constant',
      message: 'This action is used to set the dynamic information from an entity corresponding to the event.',
      description: '<h3 class="text-subtitle-1 font-weight-bold">Action parameters</h3>'
      + '<ul>'
        + '<li><strong>description</strong> (optional): the description</li>'
        + '<li><strong>name</strong>: the name of the field</li>'
        + '<li><strong>value</strong>: the value of a field</li>'
      + '</ul>',
    },
    [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfoFromTemplate]: {
      text: 'Set a string info of an entity using a template',
      message: 'This action is used to modify the dynamic information from an entity corresponding to the event.',
      description: '<h3 class="text-subtitle-1 font-weight-bold">Action parameters</h3>'
      + '<ul>'
        + '<li><strong>description</strong> (optional): the description</li>'
        + '<li><strong>name</strong>: the name of the field</li>'
        + '<li><strong>value</strong>: the template used to determine the value of the data item. Templates {{.Event.NomDuChamp}}, regular expressions or external data can be used</li>'
      + '</ul>',
    },
    [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setField]: {
      text: 'Set a field of an event to a constant',
      message: 'This action can be used to modify a field of the event.',
      description: '<h3 class="text-subtitle-1 font-weight-bold">Action parameters</h3>'
      + '<ul>'
        + '<li><strong>description</strong> (optional): the description</li>'
        + '<li><strong>name</strong>: the name of the field</li>'
        + '<li><strong>value</strong>: the new value of the field</li>'
      + '</ul>',
    },
    [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setFieldFromTemplate]: {
      text: 'Set a string field of an event using a template',
      message: 'This action allows you to modify an event field from a template.',
      description: '<h3 class="text-subtitle-1 font-weight-bold">Action parameters</h3>'
      + '<ul>'
        + '<li><strong>description</strong> (optional): the description</li>'
        + '<li><strong>name</strong>: the name of the field</li>'
        + '<li><strong>value</strong>: the template used to determine the value of the field. Templates {{.Event.NomDuChamp}}, regular expressions or external data can be used</li>'
      + '</ul>',
    },
    [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfoFromDictionary]: {
      text: 'Set entity info from a dictionary',
      message: 'This action can be used for setting entity infos from event fields with a dictionary type node.',
      description: '<h3 class="text-subtitle-1 font-weight-bold">Action parameters</h3>'
      + '<ul>'
        + '<li><strong>description</strong> (optional): the description which is used for the entity infos description. If not defined, the entity infos description is left empty</li>'
        + '<li><strong>value</strong>: the event field from which the infos are retrieved. The value must contain an array of name: value pairs</li>'
      + '</ul>',
    },
    [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setTags]: {
      text: 'Set tags from a field using regexp match',
      message: 'This action can be used for setting tags from other filtered events using regexp match.',
      description: '<p>'
        + 'The <strong>set_tags</strong> action allows for the dynamic creation of tags in the format '
        + '<strong>"Name: Value"</strong> from a field of the <strong>event currently being processed</strong>, '
        + 'using <strong>capture groups</strong> defined in a regular expression.'
        + '</p>'

        + '<p>'
        + 'This regular expression must be applied upstream, in an <strong>event filter</strong>, '
        + 'on a text field containing the information to be transformed into tags.'
        + '</p>'

        + '<p>'
        + 'The regular expression must include two named groups:'
        + '</p>'
        + '<ul>'
        + '<li><code>(?P&lt;name&gt;...)</code> to extract the tag <strong>name</strong></li>'
        + '<li><code>(?P&lt;value&gt;...)</code> to extract the tag <strong>value</strong></li>'
        + '</ul>'

        + '<p>'
        + 'Once the filter is applied and the groups are detected, the <strong>set_tags</strong> action '
        + 'uses these values to automatically generate the corresponding tags.'
        + '</p>'

        + '<h3 class="text-subtitle-1 font-weight-bold mt-4 mb-2">Regular expressions examples</h3>'
        + '<table>'
        + '<thead>'
        + '<tr>'
        + '<th class="pa-2">Format attendu dans le champ source</th>'
        + '<th class="pa-2">Expression régulière</th>'
        + '</tr>'
        + '</thead>'
        + '<tbody>'
        + '<tr>'
        + '<td class="pa-2"><code>value name;</code></td>'
        + '<td class="pa-2"><code>(?P&lt;value&gt;[a-zA-Z]+)\\s+(?P&lt;name&gt;[a-zA-Z]+);</code></td>'
        + '</tr>'
        + '<tr>'
        + '<td class="pa-2"><code>name value;</code></td>'
        + '<td class="pa-2"><code>(?P&lt;name&gt;[a-zA-Z]+)\\s+(?P&lt;value&gt;[a-zA-Z]+);</code></td>'
        + '</tr>'
        + '<tr>'
        + '<td class="pa-2"><code>name: value;</code></td>'
        + '<td class="pa-2"><code>(?P&lt;name&gt;[a-zA-Z]+):\\s+(?P&lt;value&gt;[a-zA-Z]+);</code></td>'
        + '</tr>'
        + '</tbody>'
        + '</table>'

        + '<h3 class="text-subtitle-1 font-weight-bold mt-4 mb-2">Action parameters</h3>'
        + '<ul>'
        + '<li><strong>description</strong> (optional): comment or free description of the action.</li>'
        + '<li><strong>value</strong> (required): name of the event field to which the capture groups <code>name</code> and <code>value</code> have been applied.</li>'
        + '</ul>',
    },
    [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setTagsFromTemplate]: {
      text: 'Set tags from a field using a template',
      message: 'This action can be used for setting tags from other event fields using a template.',
      description: '<p>'
        + 'This action can be used for setting tags from other event fields using a template.'
        + '</p>'

        + '<p>'
        + 'The action <strong>set_tags_from_template</strong> allows you to add a <strong>unique tag</strong> in the format '
        + '<strong>"Name: Value"</strong>, dynamically constructed from a template based on the fields of the event.'
        + '</p>'

        + '<p>'
        + 'This action is useful when you want to define a tag whose <strong>value</strong> is calculated from the content of one or more fields '
        + 'of the event, using Go templating syntax (<code>{{.Event.Field}}</code>).'
        + '</p>'

        + '<p>'
        + 'This action only allows the creation of <strong>one tag at a time</strong>.'
        + '</p>'

        + '<h3 class="text-subtitle-1 font-weight-bold mt-4 mb-2">Action parameters</h3>'
        + '<ul>'
        + '<li><strong>description</strong> (optional): a comment or free description of the action.</li>'
        + '<li><strong>name</strong> (required): the name of the tag to create.</li>'
        + '<li><strong>value</strong> (required): the template used to generate the tag\'s value.'
        + '<br>It can contain:'
        + '<ul>'
        + '<li>references to the fields of the event (<code>{{.Event.field}}</code>)</li>'
        + '<li>regular expressions if the field has been previously filtered</li>'
        + '<li>or data from an external source if it has been injected into the context.</li>'
        + '</ul>'
        + '</li>'
        + '</ul>',
    },
  },
};
