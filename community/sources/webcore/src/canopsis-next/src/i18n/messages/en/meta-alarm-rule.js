import { META_ALARMS_RULE_TYPES } from '@/constants';

export default {
  outputTemplate: 'Output template',
  thresholdType: 'Threshold type',
  thresholdRate: 'Threshold rate',
  thresholdCount: 'Threshold count',
  timeInterval: 'Time interval',
  childInactiveDelay: 'Child inactive delay',
  childInactiveDelayTooltip: 'The alarm matched with this rule is activated only after the inactivity delay',
  valuePath: 'Value path | Value paths',
  autoResolve: 'Auto resolve',
  idHelp: 'If no id is specified, a unique id will be generated automatically on rule creation',
  corelId: 'Corel ID',
  corelIdHelp: '<p>The accessible variables are: <strong>.Alarm</strong> and <strong>.Entity</strong></p>'
    + '<i>For example:</i> <span>"{{ .Alarm.Value.Connector }}", "{{ .Entity.Component }}"</span>',
  corelStatus: 'Corel status',
  corelStatusHelp: '<p>The accessible variables are: <strong>.Alarm</strong> and <strong>.Entity</strong></p>'
    + '<i>For example:</i> <span>"{{ .Alarm.Value.Connector }}", "{{ .Entity.Component }}"</span>',
  corelParent: 'Corel parent',
  corelChild: 'Corel child',
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
  errors: {
    noValuePaths: 'You have to add at least 1 value path',
  },
};
