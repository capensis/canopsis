import { ACTION_TYPES } from '@/constants';

export default {
  withAuth: 'Do you need auth fields?',
  key: 'Key',
  declareTicket: 'Declare ticket',
  workflow: 'Workflow if this action didn\'t match:',
  workflowInCaseOfSuccess: 'Workflow in case of success',
  workflowInCaseOfFailure: 'Workflow in case of failure',
  caseFailure: 'In case of failure',
  caseSuccess: 'In case of success',
  actionDidNotMatch: 'Action didn\'t match',
  remainingAction: 'Continue with remaining actions',
  remainingStep: 'Continue with remaining steps',
  addAction: 'Add action',
  emptyActions: 'No actions added yet',
  output: 'Output Action Format',
  forwardAuthor: 'Forward author to the next step',
  skipForChild: 'Skip for meta alarm children',
  skipForInstruction: 'Skip if event triggered an auto instruction',
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
  allowMultipleUrls: 'Allow multiple URLs',
  allowMultipleUrlsTooltip: 'If turned on, URL field can contain multiple URLs, it will be parsed and request will be sent to all of them.\nComma should be used as delimiter.\nFor example, URL could contain: {{ range $link := .Alarm.Value.Infos }}{{ if (eq $link.type "webhook") }}{{ $link.webhookURL }},{{ end }}{{ end }}',
  massRemove: 'Remove scenarios',
  massEnable: 'Enable scenarios',
  massDisable: 'Disable scenarios',
  actions: {
    [ACTION_TYPES.snooze]: 'Snooze',
    [ACTION_TYPES.unsnooze]: 'Unsnooze',
    [ACTION_TYPES.pbehavior]: 'Pbehavior',
    [ACTION_TYPES.pbehaviorRemove]: 'Pbehavior remove',
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
    deprecatedTriggerExist: 'This scenario is not supported anymore due to its old format and thus disabled. \n'
      + 'Please update the scenario triggers or create a new ticket declaration rule.',
    testQueryRequireSteps: 'Test query is unavailable: no webhooks were added to the scenario',
  },
  tooltips: {
    pbehaviorActionsNamePrefix: 'Name is going to be `{{prefix}} {{entity_id}} {{start}}-{{stop}}`',
  },
};
