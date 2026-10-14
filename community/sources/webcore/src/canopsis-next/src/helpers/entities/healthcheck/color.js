import { COLORS } from '@/config';
import { HEALTHCHECK_STATUSES } from '@/constants';

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

  switch (node.status) {
    case HEALTHCHECK_STATUSES.healthy:
      return COLORS.primary;
    case HEALTHCHECK_STATUSES.degraded:
      return COLORS.healthcheck.warning;
    case HEALTHCHECK_STATUSES.unhealthy:
      return COLORS.healthcheck.error;
    default:
      return COLORS.healthcheck.unknown;
  }
};
