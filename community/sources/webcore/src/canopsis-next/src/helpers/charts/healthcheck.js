import { HEALTHCHECK_NETWORK_GRAPH_OPTIONS } from '@/constants';

/**
 * @typedef {Object} HealthcheckNode
 * @property {string} name
 * @property {boolean} is_running
 * @property {boolean} is_queue_overflown
 * @property {boolean} is_too_few_instances
 * @property {boolean} is_diff_instances_config
 */

/**
 * @typedef {HealthcheckNode} HealthcheckEnginesNode
 * @property {number} instances
 * @property {number} min_instances
 * @property {number} optimal_instances
 * @property {number} queue_length
 * @property {number} time
 */

/**
 * @typedef {Object} HealthcheckEnginesEdge
 * @property {string} from
 * @property {string} to
 */

/**
 * Get node rendered position diff by factor and constants
 *
 * @param {number} factor
 * @returns {number}
 */
export const getHealthcheckNodeRenderedPositionDiff = (factor = 1) =>
  HEALTHCHECK_NETWORK_GRAPH_OPTIONS.nodeSpace * HEALTHCHECK_NETWORK_GRAPH_OPTIONS.spacingFactor * factor;

/**
 * Get translate keys for status message
 *
 * @param {HealthcheckNode} data
 * @return {string[]}
 */
export const getHealthcheckStatusMessageKeys = (data) => {
  const statusKeys = [];

  if (!data.is_running) {
    statusKeys.push('healthcheck.statuses.notRunning');
  }

  if (data.is_queue_overflown) {
    statusKeys.push('healthcheck.statuses.queueOverflow');
  }

  if (data.is_too_few_instances) {
    statusKeys.push('healthcheck.statuses.tooFewInstances');
  }

  if (data.is_diff_instances_config) {
    statusKeys.push('healthcheck.statuses.diffInstancesConfig');
  }

  return statusKeys;
};
