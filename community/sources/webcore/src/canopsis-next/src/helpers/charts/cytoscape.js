import { COLORS } from '@/config';
import { CAT_ENGINES, HEALTHCHECK_NETWORK_GRAPH_OPTIONS, HEALTHCHECK_STATUSES } from '@/constants';

/**
 * @typedef {Object} HealthcheckNode
 * @property {string} name
 * @property {number} status
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
 * Get color for a node
 *
 * @param {HealthcheckNode} node
 * @returns {string}
 */
export const getNodeColor = node => ({
  [HEALTHCHECK_STATUSES.ok]: CAT_ENGINES[node.name] ? COLORS.secondary : COLORS.primary,
  [HEALTHCHECK_STATUSES.notRunning]: COLORS.healthcheck.error,
  [HEALTHCHECK_STATUSES.unknown]: COLORS.healthcheck.unknown,
  [HEALTHCHECK_STATUSES.queueOverflow]: COLORS.healthcheck.error,
  [HEALTHCHECK_STATUSES.tooFewInstances]: COLORS.healthcheck.warning,
  [HEALTHCHECK_STATUSES.diffInstancesConfig]: COLORS.healthcheck.warning,
}[node.status] || COLORS.healthcheck.unknown);

/**
 * Get node rendered position diff by factor and constants
 *
 * @param {number} factor
 * @returns {number}
 */
export const getNodeRenderedPositionDiff = (factor = 1) =>
  HEALTHCHECK_NETWORK_GRAPH_OPTIONS.nodeSpace * HEALTHCHECK_NETWORK_GRAPH_OPTIONS.spacingFactor * factor;
