import { META_ALARMS_RULE_TYPES } from '@/constants';

export default {
  outputTemplate: 'Output template',
  threshold: 'Threshold',
  thresholdType: 'Threshold type',
  thresholdRate: 'Threshold rate',
  thresholdHelpText: 'After achieving this threshold rate, alarms that fit the patterns and created during the defined time interval are grouped\n\n'
    + 'a - alarms matched alarm and entity patterns\n'
    + 'b - entities in entity pattern\n'
    + 'c - entities in total entity pattern\n'
    + 'd - all entities in Canopsis\n'
    + '------------------------------------\n'
    + 'Case 1 - If only alarm pattern is filled\n'
    + 'Threshold rate =  a * 100 / d\n\n'
    + 'Case 2 - If entity pattern is filled, total entity pattern not filled\n'
    + 'Threshold rate = a * 100 / b\n\n'
    + 'Case 3 - If total entity pattern is filled\n'
    + 'Threshold rate = a * 100 / c',
  thresholdCount: 'Threshold count',
  thresholdCountHelpText: 'After achieving this threshold, alarms that fit the patterns and created during the defined time interval are grouped',
  timeInterval: 'Time interval',
  timeIntervalHelpText: 'Alarms created during this time interval are grouped',
  childInactiveDelay: 'Child inactive delay',
  childInactiveDelayHelpText: 'The alarm matched with this rule is activated only after the inactivity delay',
  valuePath: 'Value path | Value paths',
  addValuePath: 'Add value path',
  autoResolve: 'Auto resolve',
  idHelp: 'If no id is specified, a unique id will be generated automatically on rule creation',
  corelId: 'Corel ID',
  corelIdHelpText: 'Alarms with the same selected attribute are grouped',
  corelStatus: 'Corel status',
  corelStatusHelpText: 'By this parameter alarms are divided into parents and children',
  corelParent: 'Corel parent',
  corelParentHelpText: 'Alarms with this value of the Corel Status field are defined as parents',
  corelChild: 'Corel child',
  corelChildHelpText: 'Alarms with this value of the Corel Status field are defined as children',
  corelGroupingSummary: 'Alarms that have the same value in <strong>{corelID}</strong> field and\n'
    + '<strong>{corelChild} = {corel}</strong>\n'
    + 'will be grouped under alarm that have\n'
    + '<strong>{corelParent} = {corel}</strong>',
  outputTemplateHelp: '<p>The accessible variables are:</p>'
    + '<p><strong>.Count</strong>: The number of consequence alarms attached to the meta alarm.</p>'
    + '<p><strong>.Children</strong>: The set of variables of the last consequence alarm attached to the meta alarm.</p>'
    + '<p><strong>.Rule</strong>: The administrative information of the meta alarm itself.</p>'
    + '<p>For example:</p>'
    + '<p>Count: <strong>{{ .Count }};</strong> Children: <strong>{{ .Children.Alarm.Value.State.Message }};</strong> Rule: <strong>{{ .Rule.Name }};</strong></p>'
    + '<p>A static informative message</p>'
    + '<p>Correlated by the rule <strong>{{ .Rule.Name }}</strong></p>',
  removeConfirmationText: 'When deleting a meta alarm rule, all corresponding meta alarms will be deleted as well.\n'
    + 'Are you sure to proceed with it?\n',
  selectType: 'Select meta alarm rule type',
  grouping: 'Grouping',
  groupingTabs: {
    existing: 'Group under existing alarm',
    createNew: 'Group and create new meta alarm',
  },
  groupingLabels: {
    groupUnder: 'Group under',
    groupBy: 'Group by',
  },
  valuePathHelpText: 'Alarms with the same combination of values in these fields will be grouped under the same meta alarm ',
  componentTemplate: 'Component template',
  resourceTemplate: 'Resource template',
  copyTagsFromChildren: 'Copy tags from children alarms',
  filterByLabelEnabled: 'Filter by label',
  filterByLabelEnabledTooltip: 'Some tags can be defined in the format Tag:Value, e.g. Env:Prod.\nWith filter by label, only tags with the defined label will be copied from children alarms to the meta alarm.',
  copyFromLastChild: 'Copy from last child',
  copyFromLastChildTooltip: 'When enabled, the infos value is copied from an infos of a last child alarm. The child alarm infos name shall be defined in this case.',
  massRemove: 'Remove meta alarm rules',
  massEnable: 'Enable meta alarm rules',
  massDisable: 'Disable meta alarm rules',
  steps: {
    basics: 'Basics',
    defineType: 'Define type',
    addParameters: 'Add parameters',
  },
  types: {
    [META_ALARMS_RULE_TYPES.relation]: {
      label: 'Parent component',
      text: 'Group under parent component',
      helpText: 'Define the patterns that should have all alarms raised on its dependencies grouped',
    },
    [META_ALARMS_RULE_TYPES.timebased]: {
      label: 'Time interval (when severity changed)',
      text: 'Group by creation date interval',
      helpText: 'All alarms that match the patterns and raised during a defined time interval are grouped',
    },
    [META_ALARMS_RULE_TYPES.attribute]: {
      label: 'Parent only',
      text: 'Group by pattern only',
      helpText: 'All alarms that have the attribute(s) defined by filter patterns are grouped',
    },
    [META_ALARMS_RULE_TYPES.complex]: {
      label: 'With trigger threshold or rate',
      text: 'Group with trigger threshold or rate',
      helpText: 'All alarms that have the attribute(s) defined by filter patterns, time interval, and a trigger threshold or rate are grouped',
    },
    [META_ALARMS_RULE_TYPES.valuegroup]: {
      label: 'With trigger threshold or rate + custom attributes values',
      text: 'Group by group of values',
      helpText: 'All alarms that have the attribute(s) defined by filter patterns, time interval, trigger threshold or rate, and the value path are grouped',
    },
    [META_ALARMS_RULE_TYPES.corel]: {
      label: 'Custom defined parent',
      text: 'Group under custom defined parent',
      helpText: 'All alarms that have the attribute(s) defined by filter patterns, time interval, threshold count, and the correlation identifiers are grouped',
    },
  },
  patternsTabLabel: 'Pattern and params',
  errors: {
    noValuePaths: 'You have to add at least 1 value path',
  },
  field: {
    title: 'Meta alarm rule',
    noData: 'No meta alarm rules is found according to the patterns defined',
  },
};
