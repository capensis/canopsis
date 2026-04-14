import { PATTERN_TYPES } from './pattern';

export const LLM_MODEL_TYPES = {
  gemini: 'gemini',
};

export const LLM_THINKING_LEVELS = {
  minimal: 'minimal',
  low: 'low',
  medium: 'medium',
  high: 'high',
};

export const LLM_AI_CHAT_WIDTH = 400;

export const LLM_AI_CHAT_PROMPT_MAX_LENGTH = 10000;

export const LLM_AI_CHAT_SUGGESTION_TYPES = {
  createPattern: 'createPattern',
  editPattern: 'editPattern',
  validatePattern: 'validatePattern',
};

export const LLM_AI_CHAT_MESSAGE_ROLES = {
  user: 'user',
  model: 'model',
};

export const LLM_AI_CHAT_MESSAGE_TYPES = {
  send: 'send',
  cancel: 'cancel',
};

export const LLM_AI_CHAT_ERROR_CODES = {
  gone: 410,
  tooManyRequests: 429,
  internalError: 500,
  timeout: 504,
};

export const LLM_AI_CHAT_TOURS = {
  sidebarMinimized: 'ai_chat_minimized',
  hiddenWarningAlert: 'ai_chat_hidden_warning_alert',
};

export const LLM_SOCKET_CONTEXTS = {
  idleRule: 'idle_rule',
  scenario: 'scenario', // TODO: here we have array of patterns for each action. Need to check.
  flappingRule: 'flapping_rule',
  resolveRule: 'resolve_rule',
  alarmTag: 'alarm_tag',
  linkRule: 'link_rule',
  instruction: 'instruction',
  dynamicInfos: 'dynamic_infos',
  metaAlarmRule: 'meta_alarm_rule',
  declareTicketRule: 'declare_ticket_rule',
  pbehavior: 'pbehavior', // TODO: here we have popups with form in modal window. Need to check.
  entityService: 'entity_service',
  stateSettings: 'state_settings',
  kpiFilter: 'kpi_filter',
  eventFilter: 'eventfilter',
  eventRecord: 'event_record',
  widgetFilter: 'widget_filter',
  corporateAlarmPattern: 'corporate_alarm_pattern',
  corporateEntityPattern: 'corporate_entity_pattern',
  corporatePbehaviorPattern: 'corporate_pbehavior_pattern',
  corporateWeatherServicePattern: 'corporate_weather_service_pattern',
};

export const PATTERN_TYPES_TO_LLM_SOCKET_CONTEXTS = {
  [PATTERN_TYPES.alarm]: LLM_SOCKET_CONTEXTS.corporateAlarmPattern,
  [PATTERN_TYPES.entity]: LLM_SOCKET_CONTEXTS.corporateEntityPattern,
  [PATTERN_TYPES.pbehavior]: LLM_SOCKET_CONTEXTS.corporatePbehaviorPattern,
  [PATTERN_TYPES.serviceWeather]: LLM_SOCKET_CONTEXTS.corporateWeatherServicePattern,
};
