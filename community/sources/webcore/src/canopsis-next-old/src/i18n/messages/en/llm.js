import {
  LLM_AI_CHAT_ERROR_CODES,
  LLM_AI_CHAT_SUGGESTION_TYPES,
  LLM_SOCKET_CONTEXTS,
  LLM_THINKING_LEVELS,
} from '@/constants';

export default {
  expandTabs: {
    promptsHistory: 'Prompts history',
  },
  massEnable: 'Enable LLMs',
  massDisable: 'Disable LLMs',
  promptsHistory: {
    tabs: {
      allPrompts: 'All prompts',
    },
    expandUserHistory: {
      lastPromptDate: 'Last prompt date',
    },
    filterByUser: 'User filter',
    groupByChat: 'Group by chat',
    notRelatedToCanopsis: 'Not related to Canopsis',
    allUsers: 'All users',
    search: 'Search by user, modal or name',
    searchByUserName: 'Search by user name',
    searchByModalOrName: 'Search by modal or name',
    userHistoryTitle: 'User history',
    userHistoryEmpty: 'No users yet for this configuration.',
    ruleNotSaved: 'Not saved',
    columns: {
      userName: 'User name',
      datetime: 'Datetime',
      tokensUsed: 'Tokens used',
      modal: 'Modal',
      usage: 'Usage',
      canopsisRelated: 'Canopsis related',
      prompt: 'Prompt',
      lastUsed: 'Last used',
    },
    contextTitles: {
      [LLM_SOCKET_CONTEXTS.idleRule]: 'Idle rule',
      [LLM_SOCKET_CONTEXTS.scenario]: 'Scenario',
      [LLM_SOCKET_CONTEXTS.flappingRule]: 'Flapping rule',
      [LLM_SOCKET_CONTEXTS.resolveRule]: 'Resolve rule',
      [LLM_SOCKET_CONTEXTS.alarmTag]: 'Alarm tag',
      [LLM_SOCKET_CONTEXTS.linkRule]: 'Link rule',
      [LLM_SOCKET_CONTEXTS.instruction]: 'Remediation instruction',
      [LLM_SOCKET_CONTEXTS.dynamicInfos]: 'Dynamic information',
      [LLM_SOCKET_CONTEXTS.metaAlarmRule]: 'Meta alarm rule',
      [LLM_SOCKET_CONTEXTS.declareTicketRule]: 'Declare ticket rule',
      [LLM_SOCKET_CONTEXTS.pbehavior]: 'Pbehavior',
      [LLM_SOCKET_CONTEXTS.entityService]: 'Service',
      [LLM_SOCKET_CONTEXTS.stateSettings]: 'State setting',
      [LLM_SOCKET_CONTEXTS.kpiFilter]: 'KPI filter',
      [LLM_SOCKET_CONTEXTS.eventFilter]: 'Event filter',
      [LLM_SOCKET_CONTEXTS.eventRecord]: 'Events record',
      [LLM_SOCKET_CONTEXTS.widgetFilter]: 'Widget filter',
      [LLM_SOCKET_CONTEXTS.corporateAlarmPattern]: 'Shared alarm filter',
      [LLM_SOCKET_CONTEXTS.corporateEntityPattern]: 'Shared entity filter',
      [LLM_SOCKET_CONTEXTS.corporatePbehaviorPattern]: 'Shared pbehavior filter',
      [LLM_SOCKET_CONTEXTS.corporateWeatherServicePattern]: 'Shared service weather filter',
    },
  },
  modelType: 'LLM model type',
  modelName: 'LLM model name',
  model: 'LLM model',
  recommendedBadge: 'recommended',
  apiKey: 'API key',
  apiKeyPlaceholder: 'New API key',
  thinkingLevel: 'Thinking level',
  isDefaultModel: 'Is default model',
  currentDefaultModelLine: 'Current default model:',
  lastUsedDate: 'Last used date',
  importantNotesMessage:
    '<div class="font-weight-regular">Important notes:</div>'
    + '<ul>'
    + '<li>Gemini models can be enabled or disabled through this module; however, '
    + '<strong>token limits and quota management</strong>'
    + ' are configured directly within the '
    + '<strong><a href="{geminiConsoleUrl}" target="_blank">Gemini AI Console</a></strong>.</li>'
    + '<li><strong>Regular monitoring</strong> of Google Cloud limits is required, as <strong>API policies and usage conditions are subject to change.</strong></li>'
    + '</ul>',
  thinkingLevels: {
    [LLM_THINKING_LEVELS.minimal]: 'Minimal',
    [LLM_THINKING_LEVELS.low]: 'Low',
    [LLM_THINKING_LEVELS.medium]: 'Medium',
    [LLM_THINKING_LEVELS.high]: 'High',
  },
  chat: {
    title: 'AI assistant',
    howCanIHelp: 'How can I help?',
    promptPlaceholder: 'Describe what you need…',
    ask: 'Ask',
    thinking: 'Thinking…',
    tryLabel: 'Try:',
    restartConfirmation: 'Restart the chat? <strong>Chat history</strong> and <strong>versions</strong> will be cleared',
    suggestions: {
      createPattern: 'Create pattern',
      editPattern: 'Edit pattern',
      validatePattern: 'Validate pattern',
    },
    patternsMessage: '{patterns} pattern | {patterns} patterns',
    patternsEditedMessage: 'You edited {patterns} pattern | You edited {patterns} patterns',
    patternCreatedMessage: 'Patterns created',
    patternUpdatedMessage: 'Patterns updated',
    emptyPatternsMessage: 'You removed all patterns.\nIf you need to clear the context, restart AI assistant.',
    patternCannotBeCreatedReasons: 'Pattern cannot be created for the following reasons:\n',
    fixPatternPrompt: 'Fix pattern:\n{jsonString}',
    suggestionPrompts: {
      [LLM_AI_CHAT_SUGGESTION_TYPES.createPattern]: 'Create pattern with following conditions:\n',
      [LLM_AI_CHAT_SUGGESTION_TYPES.editPattern]: 'The following changes are needed in pattern:\n',
      [LLM_AI_CHAT_SUGGESTION_TYPES.validatePattern]: 'Fix the following pattern if needed:\n',
    },
    pattern: {
      restoreVersion: 'Restore version',
      seePattern: 'See pattern',
      hidePattern: 'Hide pattern',
      version: 'Version {version}',
      versionRestored: 'Version {version} restored',
    },
    patternsItemsLabel: {
      [LLM_SOCKET_CONTEXTS.scenario]: 'Scenario action',
      [LLM_SOCKET_CONTEXTS.stateSettings]: 'Target pattern',
    },
    infoAlerts: {
      [LLM_SOCKET_CONTEXTS.scenario]: 'Works only with pattern for action <strong>{patternItem}</strong><br />',
      [LLM_SOCKET_CONTEXTS.stateSettings]: 'Works only with pattern for <strong>{patternItem}</strong><br />',
    },
    infoAlertEnding: '<strong>Click Restart to run AI assistant for another pattern</strong>',
    patternsItemsError: {
      [LLM_SOCKET_CONTEXTS.scenario]:
        'Selected action is deleted.<br />'
        + '<strong>You can restart the chat with another action (the history will be lost).</strong>',
    },
    warningAlert: '<div>Please be aware that:</div>'
      + '<ul>'
      + '<li>AI-generated output <strong>may contain errors</strong> — always review and validate patterns before use.</li>'
      + '<li>This chat session is <strong>temporary</strong> and will be <strong>cleared</strong> when the modal is closed.</li>'
      + '</ul>',
    errors: {
      [LLM_AI_CHAT_ERROR_CODES.badRequest]:
        'The request could not be processed.<br />'
        + '<strong>Please check your input and try again.</strong>',
      [LLM_AI_CHAT_ERROR_CODES.gone]:
        'Selected AI model is currently unavailable.<br />'
        + '<strong>Please contact your administrator. You can restart the chat with another model (the history will be lost).</strong>',
      [LLM_AI_CHAT_ERROR_CODES.tooManyRequests]:
        'You exceeded your current quota.<br />'
        + '<strong>Please try again at {retryAt} or restart the chat with another model (the history will be lost).</strong>',
      [LLM_AI_CHAT_ERROR_CODES.internalError]:
        'Something went wrong.<br />'
        + '<strong>Please contact your administrator</strong>',
      [LLM_AI_CHAT_ERROR_CODES.timeout]:
        'Request timed out.<br />'
        + '<strong>Please try again later.</strong>',
      noModels:
        'No AI models are currently available.<br />'
        + '<strong>Please contact your administrator.</strong>',
      serverError:
        'Problem with socket connection.<br />'
        + '<strong>Please contact your administrator.</strong>',
    },
  },
};
