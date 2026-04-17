import { HEALTHCHECK_NETWORK_GRAPH_OPTIONS } from '@/constants';

/**
 * @typedef {Object} HealthcheckNode
 * @property {string} name
 * @property {boolean} is_running
 * @property {boolean} [is_queue_overflown]
 * @property {boolean} [is_too_few_instances]
 * @property {boolean} [is_diff_instances_config]
 * @property {boolean} [is_unknown]
 * @property {number} [nodes]
 * @property {number} [running_nodes]
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
export const getHealthcheckNodeRenderedPositionDiff = (factor = 1) => HEALTHCHECK_NETWORK_GRAPH_OPTIONS.nodeSpace
  * HEALTHCHECK_NETWORK_GRAPH_OPTIONS.spacingFactor
  * factor;

/**
 * Get HTML div element for healthcheck network graph info label
 *
 * @param {string} message
 * @param {boolean} [hasError=false]
 * @returns {string}
 */
export const getHealthcheckInfoDiv = (message, hasError = false) => `<div class="body-1 grey--text darken-3 pre-wrap ${hasError ? 'error--text' : ''}">${message}</div>`;
