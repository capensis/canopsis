import { get, camelCase, isNumber } from 'lodash';
import tinycolor from 'tinycolor2';

import { COLORS } from '@/config';
import { PRO_ENGINES, COLOR_INDICATOR_TYPES, ENTITIES_STATES_STYLES, EVENT_ENTITY_COLORS_BY_TYPE } from '@/constants';

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
 * Get color by entity impact state
 *
 * @param {number} value
 * @returns {string}
 */
export const getImpactStateColor = value => COLORS.impactState[value];

/**
 * Get color by entity impact state
 *
 * @param {number} value
 * @returns {string}
 */
export const getEntityStateColor = value => get(ENTITIES_STATES_STYLES, [value, 'color']);

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
    const state = isNumber(entity.state) ? entity.state : entity.state?.val;

    return getEntityStateColor(state);
  }

  return getImpactStateColor(entity.impact_state);
};

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

/**
 * Convert color to rgb
 *
 * @param {string|Object} color
 * @return {string}
 */
export const colorToRgb = color => tinycolor(color).toRgbString();

/**
 * Convert color to rgba with alpha
 *
 * @param {string|Object} color
 * @param {number} alpha
 * @return {string}
 */
export const colorToRgba = (color, alpha = 1.0) => tinycolor(color)
  .setAlpha(alpha)
  .toRgbString();

/**
 * Convert color to hex
 *
 * @param {string|Object} color
 * @return {string}
 */
export const colorToHex = color => tinycolor(color).toHexString();

/**
 * Check color is valid
 *
 * @param {string|Object} color
 * @return {boolean}
 */
export const isValidColor = color => tinycolor(color).isValid();

/**
 * Get color for metric
 *
 * @param {string} metric
 */
export const getMetricColor = metric => COLORS.metrics[camelCase(metric)] || COLORS.secondary;

/**
 * Get color for entity event
 *
 * @param {string} type
 */
export const getEntityEventColor = type => EVENT_ENTITY_COLORS_BY_TYPE[type];

/**
 * Get darken color
 *
 * @param {string} color
 * @param {number} amount
 */
export const getDarkenColor = (color, amount) => tinycolor(color)
  .darken(amount)
  .toString();
