import { get } from 'lodash';
import tinycolor from 'tinycolor2';

import { COLORS } from '@/config';
import { CAT_ENGINES, COLOR_INDICATOR_TYPES, ENTITIES_STATES_STYLES } from '@/constants';

/**
 * Get most readable text color ('white' or 'black')
 *
 * @param {string} color
 * @param {{ level: 'AA' | 'AAA', size: 'small' | 'large' }} [options = {}]
 */
export const getMostReadableTextColor = (color, options = {}) => {
  if (!color) {
    return 'black';
  }

  const isWhiteReadable = tinycolor.isReadable(color, 'white', options);

  return isWhiteReadable ? 'white' : 'black';
};

/**
 * Get color for a entity by colorIndicator and isGrey parameters
 *
 * @param {Service | Entity | {}} [entity = {}]
 * @param {string} [colorIndicator = COLOR_INDICATOR_TYPES.state]
 * @returns {string|*}
 */
export const getEntityColor = (entity = {}, colorIndicator = COLOR_INDICATOR_TYPES.state) => {
  if (entity.is_grey) {
    return COLORS.state.pause;
  }

  if (colorIndicator === COLOR_INDICATOR_TYPES.state) {
    const state = get(entity, 'state.val');

    return get(ENTITIES_STATES_STYLES, [state, 'color']);
  }

  return COLORS.impactState[entity.impact_state];
};

/**
 * Get color for a node
 *
 * @param {HealthcheckNode} node
 * @returns {string}
 */
export const getHealthcheckNodeColor = (node) => {
  if (node.is_unknown) {
    return COLORS.healthcheck.unknown;
  }

  if (!node.is_running || node.is_queue_overflown) {
    return COLORS.healthcheck.error;
  }

  if (node.is_too_few_instances || node.is_diff_instances_config) {
    return COLORS.healthcheck.warning;
  }

  return CAT_ENGINES.includes(node.name) ? COLORS.secondary : COLORS.primary;
};
