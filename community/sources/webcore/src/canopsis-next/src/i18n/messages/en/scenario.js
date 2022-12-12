import { ACTION_TYPES } from '@/constants';

export default {
  triggers: 'Triggers',
  emitTrigger: 'Emit trigger',
  withAuth: 'Do you need auth fields?',
  emptyResponse: 'Empty response',
  isRegexp: 'The value can be a RegExp',
  headerKey: 'Header key',
  headerValue: 'Header value',
  key: 'Key',
  skipVerify: 'Ignore HTTPS certificate verification',
  headers: 'Headers',
  declareTicket: 'Declare ticket',
  workflow: 'Workflow if this action didn’t match:',
  remainingAction: 'Continue with remaining actions',
  addAction: 'Add action',
  emptyActions: 'No actions added yet',
  output: 'Output Action Format',
  forwardAuthor: 'Forward author to the next step',
  urlHelp: '<p>The accessible variables are: <strong>.Alarm</strong>, <strong>.Entity</strong> and <strong>.Children</strong></p>'
    + '<i>For example:</i>'
    + '<pre>"https://exampleurl.com?resource={{ .Alarm.Value.Resource }}"</pre>'
    + '<pre>"https://exampleurl.com?entity_id={{ .Entity.ID }}"</pre>'
    + '<pre>"https://exampleurl.com?children_count={{ len .Children }}"</pre>'
    + '<pre>"https://exampleurl.com?children={{ range .Children }}{{ .ID }}{{ end }}"</pre>',
  outputHelp: '<p>The accessible variables are: <strong>.Alarm</strong> and <strong>.Entity</strong></p>'
    + '<i>For example:</i>'
    + '<pre>Resource - {{ .Alarm.Value.Resource }}. Entity - {{ .Entity.ID }}.</pre>',
  payloadHelp: '<p>The accessible variables are: <strong>.Alarm</strong>, <strong>.Entity</strong> and <strong>.Children</strong></p>'
    + '<i>For example:</i>'
    + '<pre>{\n'
    + '  resource: "{{ .Alarm.Value.Resource }}",\n'
    + '  entity: "{{ .Entity.ID }}",\n'
    + '  children_count: "{{ len .Children }}",\n'
    + '  children: {{ range .Children }}{{ .ID }}{{ end }}\n'
    + '}</pre>',
  actions: {
    [ACTION_TYPES.snooze]: 'Snooze',
    [ACTION_TYPES.pbehavior]: 'Pbehavior',
    [ACTION_TYPES.changeState]: 'Change state (Change and lock severity)',
    [ACTION_TYPES.ack]: 'Acknowledge',
    [ACTION_TYPES.ackremove]: 'Acknowledge remove',
    [ACTION_TYPES.assocticket]: 'Associate ticket',
    [ACTION_TYPES.cancel]: 'Cancel',
    [ACTION_TYPES.webhook]: 'Webhook',
  },
  tabs: {
    pattern: 'Pattern',
  },
  errors: {
    actionRequired: 'Please add at least one action',
    priorityExist: 'The priority of current scenario is already in use. Do you want to change the current scenario priority to {priority}?',
  },
};
