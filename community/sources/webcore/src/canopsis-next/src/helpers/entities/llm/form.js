import { LLM_MODEL_TYPES } from '@/constants';

import { durationToForm } from '@/helpers/date/duration';

/**
 * @typedef {Object} LlmConfig
 * @property {boolean} enabled
 * @property {string} name
 * @property {string} type
 * @property {string} api_key
 * @property {Duration} request_timeout
 * @property {string} model
 * @property {string} thinking_level
 * @property {boolean} default
 */

/**
 * Convert LLM entity to form object
 *
 * @param {LlmConfig} [llm={}]
 * @returns {LlmConfig}
 */
export const llmToForm = (llm = {}) => ({
  enabled: llm.enabled ?? true,
  type: LLM_MODEL_TYPES.gemini,
  name: llm.name ?? '',
  api_key: llm.api_key ?? '',
  request_timeout: durationToForm(llm.request_timeout ?? {}),
  model: llm.model ?? '',
  thinking_level: llm.thinking_level ?? '',
  default: llm.default ?? false,
});

/**
 * Convert form object to LLM entity
 *
 * @param {LlmConfig} form
 * @returns {LlmConfig}
 */
export const formToLlm = form => ({
  enabled: form.enabled,
  name: form.name,
  type: form.type,
  api_key: form.api_key,
  request_timeout: form.request_timeout,
  model: form.model,
  thinking_level: form.thinking_level || null,
  default: form.default,
});
