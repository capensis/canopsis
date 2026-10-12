import { LLM_MODEL_TYPES } from '@/constants';

/**
 * @typedef {Object} LlmConfig
 * @property {boolean} enabled
 * @property {string} name
 * @property {string} type
 * @property {string} api_key
 * @property {string} model
 * @property {string | null} thinking_level
 * @property {boolean} default
 */

/**
 * Convert LLM entity to form object
 *
 * @param {LlmConfig} [llm={}]
 * @param {boolean} [defaultValue=false]
 * @returns {LlmConfig}
 */
export const llmToForm = (llm = {}, defaultValue = false) => ({
  enabled: llm.enabled ?? true,
  type: LLM_MODEL_TYPES.gemini,
  name: llm.name ?? '',
  api_key: llm.api_key ?? '',
  model: llm.model ?? '',
  thinking_level: llm.thinking_level ?? '',
  default: llm.default ?? defaultValue,
});

/**
 * Convert form object to LLM entity
 *
 * @param {LlmConfig} form
 * @returns {LlmConfig}
 */
export const formToLlm = form => ({
  ...form,

  thinking_level: form.thinking_level || null,
});
