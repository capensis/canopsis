import { META_ALARMS_RULE_TYPES } from '@/constants';

export default {
  outputTemplate: 'Output template',
  thresholdType: 'Threshold type',
  thresholdRate: 'Threshold rate',
  thresholdRateHelpText: 'After achieving this threshold rate, alarms that fit the patterns and created during the defined time interval are grouped',
  thresholdCount: 'Threshold count',
  thresholdCountHelpText: 'After achieving this threshold, alarms that fit the patterns and created during the defined time interval are grouped',
  timeInterval: 'Time interval',
  timeIntervalHelpText: 'Alarms created during this time interval are grouped',
  childInactiveDelay: 'Child inactive delay',
  childInactiveDelayHelpText: 'The alarm matched with this rule is activated only after the inactivity delay',
  valuePath: 'Value path | Value paths',
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
  valuePathHelpText: 'Custom attribute for grouping alarms defined by a value path',
  steps: {
    basics: 'Basics',
    defineType: 'Define type',
    addParameters: 'Add parameters',
  },
  types: {
    [META_ALARMS_RULE_TYPES.relation]: {
      text: 'Parent-child relationship',
      helpText: 'All alarms raised on dependent entities are grouped',
    },
    [META_ALARMS_RULE_TYPES.timebased]: {
      text: 'Grouping by time interval',
      helpText: 'All alarms raised during a defined time interval are grouped',
    },
    [META_ALARMS_RULE_TYPES.attribute]: {
      text: 'Grouping by attribute',
      helpText: 'All alarms that fit a pattern with defined attributes are grouped',
    },
    [META_ALARMS_RULE_TYPES.complex]: {
      text: 'Complex grouping with trigger threshold or rate',
      helpText: 'All alarms that fit a pattern with defined attributes during the defined time interval are grouped',
    },
    [META_ALARMS_RULE_TYPES.valuegroup]: {
      text: 'Grouping by group of values',
      helpText: 'It is a complex grouping with value paths as additional parameters for grouping',
    },
    [META_ALARMS_RULE_TYPES.corel]: {
      text: 'Grouping by correlation identifiers',
      helpText: 'Grouping for already correlated alarms: all alarms of the same correlation identifier are grouped',
    },
  },
  parametersTitle: {
    [META_ALARMS_RULE_TYPES.relation]: 'Parent-child relationship',
    [META_ALARMS_RULE_TYPES.timebased]: 'Time based relationship',
    [META_ALARMS_RULE_TYPES.attribute]: 'Attribute based relationship',
    [META_ALARMS_RULE_TYPES.complex]: 'Complex grouping with a trigger threshold or rate',
    [META_ALARMS_RULE_TYPES.valuegroup]: 'Grouping by group of values',
    [META_ALARMS_RULE_TYPES.corel]: 'Grouping by correlation identifiers',
  },
  parametersDescription: {
    [META_ALARMS_RULE_TYPES.relation]: 'Define the filter patterns for entities that should have all alarms raised on its dependencies grouped',
    [META_ALARMS_RULE_TYPES.timebased]: 'All alarms that match the patterns and raised during a defined time interval are grouped',
    [META_ALARMS_RULE_TYPES.attribute]: 'All alarms that have the attribute(s) defined by filter patterns are grouped',
    [META_ALARMS_RULE_TYPES.complex]: 'All alarms that have the attribute(s) defined by filter patterns, time interval, and a trigger threshold or rate are grouped',
    [META_ALARMS_RULE_TYPES.valuegroup]: 'All alarms that have the attribute(s) defined by filter patterns, time interval, threshold or rate, and the value path are grouped',
    [META_ALARMS_RULE_TYPES.corel]: 'All alarms that have the attribute(s) defined by filter patterns, time interval, threshold count, and the correlation identifiers are grouped',
  },
  errors: {
    noValuePaths: 'You have to add at least 1 value path',
  },
  field: {
    title: 'Meta alarm rule',
    noData: 'No meta alarm rules is found according to the patterns defined',
  },
};
