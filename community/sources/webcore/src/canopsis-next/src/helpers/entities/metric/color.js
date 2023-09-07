import { camelCase } from 'lodash';

import { COLORS } from '@/config';

/**
 * Get color for metric
 *
 * @param {string} metric
 */
export const getMetricColor = metric => COLORS.metrics[camelCase(metric)] || COLORS.secondary;
