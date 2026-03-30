import { LLM_AI_CHAT_SUGGESTION_TYPES, LLM_THINKING_LEVELS } from '@/constants';

export default {
  expandTabs: {
    details: 'Details',
    promptsHistory: 'Prompts history',
  },
  promptsHistory: {
    tabs: {
      allPrompts: 'All prompts',
      byUser: 'By user',
    },
    searchPlaceholder: 'Search by user, modal or name',
    notRelatedToCanopsis: 'Not related to Canopsis',
    groupByChat: 'Group by chat',
    seeChat: 'See chat',
    columns: {
      userName: 'User name',
      datetime: 'Datetime',
      tokensUsed: 'Tokens used',
      modal: 'Modal',
      name: 'Name',
      usage: 'Usage',
      canopsisRelated: 'Canopsis related',
      prompt: 'Prompt',
      seeChat: 'See chat',
      promptsCount: 'Prompts',
      lastUsed: 'Last used',
    },
  },
  modelType: 'LLM model type',
  modelTypes: {
    gemini: 'Gemini',
  },
  modelName: 'LLM model name',
  model: 'LLM model',
  recommendedBadge: 'recommended',
  apiKey: 'API key',
  apiKeyPlaceholder: 'New API key',
  thinkingLevel: 'Thinking level',
  isDefaultModel: 'Is default model',
  currentDefaultModelLine: 'Current default model: {name}',
  lastUsedDate: 'Last used date',
  expandNoDetails: 'No additional details for this configuration.',
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
    modelPlaceholder: 'Model',
    ask: 'Ask',
    expandMessage: 'Show full message',
    collapseMessage: 'Show less',
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
  },
};
