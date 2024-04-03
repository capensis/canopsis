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
  steps: {
    basics: 'Basics',
    defineType: 'Define type',
    addParameters: 'Add parameters',
  },
  errors: {
    noValuePaths: 'You have to add at least 1 value path',
  },
};
