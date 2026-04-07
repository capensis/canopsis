import { LLM_AI_CHAT_ERROR_CODES, LLM_AI_CHAT_SUGGESTION_TYPES, LLM_THINKING_LEVELS } from '@/constants';

export default {
  expandTabs: {
    promptsHistory: 'Prompts history',
  },
  promptsHistory: {
    tabs: {
      allPrompts: 'All prompts',
    },
    expandUserHistory: {
      empty: 'No prompts for this user.',
      lastPromptDate: 'Last prompt date',
    },
    filterByUser: 'User filter',
    groupByChat: 'Group by chat',
    onlyOffTopic: 'Only off-topic',
    allUsers: 'All users',
    search: 'Search by user, modal or name',
    searchByUserName: 'Search by user name',
    searchByModalOrName: 'Search by modal or name',
    userHistoryTitle: 'User history',
    userHistoryEmpty: 'No users yet for this configuration.',
    columns: {
      userName: 'User name',
      datetime: 'Datetime',
      tokensUsed: 'Tokens used',
      context: 'Context',
      name: 'Name',
      offTopic: 'Off topic',
      prompt: 'Prompt',
      lastUsed: 'Last used',
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
  currentDefaultModelLine: 'Current default model: {name}',
  lastUsedDate: 'Last used date',
  importantNotesMessage:
    '<div class="font-weight-regular">Important notes:</div>'
    + '<ul>'
    + '<li>Gemini models can be enabled or disabled through this module; however, '
    + '<strong>token limits and quota management</strong>'
    + ' are configured directly within the '
    + '<strong><a href="{geminiConsoleUrl}" target="_blank">Gemini AI Console</a></strong>.</li>'
    + '<li><strong>Regular monitoring</strong> of Google Cloud limits is required, as <strong>API policies</strong> and usage conditions are subject to change.</li>'
    + '</ul>',
  thinkingLevels: {
    [LLM_THINKING_LEVELS.minimal]: 'Minimal',
    [LLM_THINKING_LEVELS.low]: 'Low',
    [LLM_THINKING_LEVELS.medium]: 'Medium',
    [LLM_THINKING_LEVELS.high]: 'High',
  },
  chat: {
    howCanIHelp: 'How can I help?',
    promptPlaceholder: 'Describe what you need…',
    ask: 'Ask',
    chatHistoryTitle: 'Chat history',
    thinking: 'Thinking…',
    tryLabel: 'Try:',
    suggestions: {
      createPattern: 'Create pattern',
      editPattern: 'Edit pattern',
      validatePattern: 'Validate pattern',
    },
    patternsMessage: '{patterns} pattern | {patterns} patterns',
    patternsEditedMessage: 'You edited {patterns} pattern | You edited {patterns} patterns',
    patternCreatedMessage: 'Patterns created',
    patternUpdatedMessage: 'Patterns updated',
    emptyPatternsMessage: 'Empty patterns',
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
    warningAlert: '<div>Please be aware that:</div>'
      + '<ul>'
      + '<li>AI-generated output <strong>may contain errors</strong> — always review and validate patterns before use.</li>'
      + '<li>This chat session is <strong>temporary</strong> and will be <strong>cleared</strong> when the modal is closed.</li>'
      + '</ul>',
    errors: {
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
    },
  },
};
