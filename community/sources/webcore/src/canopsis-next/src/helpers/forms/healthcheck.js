import { HEALTHCHECK_ENGINES_NAMES } from '@/constants';

/**
 * @typedef {Object} EngineParameters
 * @property {boolean} enabled
 * @property {number} minimal
 * @property {number} optimal
 */

/**
 * @typedef {Object} QueueLimit
 * @property {boolean} enabled
 * @property {number} limit
 */

/**
 * @typedef {Object} HealthcheckParameters
 * @property {QueueLimit} queue
 * @property {QueueLimit} messages
 * @property {EngineParameters} engine-webhook
 * @property {EngineParameters} engine-fifo
 * @property {EngineParameters} engine-axe
 * @property {EngineParameters} engine-che
 * @property {EngineParameters} engine-pbehavior
 * @property {EngineParameters} engine-action
 * @property {EngineParameters} engine-service
 * @property {EngineParameters} dynamic-infos
 * @property {EngineParameters} engine-correlation
 * @property {EngineParameters} engine-remediation
 */

/**
 * Convert healthcheck parameters to form
 *
 * @param {HealthcheckParameters} [healthcheckParameters = {}]
 * @return {HealthcheckParameters}
 */
export const healthcheckParametersToForm = (healthcheckParameters = {}) => ({
  queue: healthcheckParameters.queue || {
    limit: 0,
    enabled: false,
  },
  messages: healthcheckParameters.messages || {
    limit: 0,
    enabled: false,
  },
  ...Object.values(HEALTHCHECK_ENGINES_NAMES).reduce((acc, engineName) => {
    const engineParameters = healthcheckParameters[engineName] || {};

    acc[engineName] = {
      enabled: !!engineParameters.enabled,
      minimal: engineParameters.minimal,
      optimal: engineParameters.optimal,
    };

    return acc;
  }, {}),
});
