import { COLORS } from '@/config';
import { PRO_ENGINES } from '@/constants';

/**
 * Get color for a node
 *
 * @param {HealthcheckNode} node
 * @returns {string}
 */
export const getHealthcheckNodeColor = (node = {}) => {
  if (node.is_unknown) {
    return COLORS.healthcheck.unknown;
  }

  if (!node.is_running || node.is_queue_overflown) {
    return COLORS.healthcheck.error;
  }

  if (node.is_too_few_instances || node.is_diff_instances_config) {
    return COLORS.healthcheck.warning;
  }

  return PRO_ENGINES.includes(node.name) ? COLORS.secondary : COLORS.primary;
};
