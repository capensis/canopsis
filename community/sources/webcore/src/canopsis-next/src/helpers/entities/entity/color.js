import { isNumber } from 'lodash';

import { CSS_COLORS_VARS } from '@/config';
import { COLOR_INDICATOR_TYPES, EVENT_ENTITY_COLORS_BY_TYPE } from '@/constants';

import { getAlarmImpactStateColor, getAlarmStateColor } from '../alarm/color';

/**
 * Get color for a entity by colorIndicator and isGrey parameters
 *
 * @param {Service | Entity | {}} [entity = {}]
 * @param {string} [colorIndicator = COLOR_INDICATOR_TYPES.state]
 * @returns {string|*}
 */
export const getEntityColor = (entity = {}, colorIndicator = COLOR_INDICATOR_TYPES.state) => {
  if (entity.is_grey) {
    return CSS_COLORS_VARS.state.pause;
  }

  if (colorIndicator === COLOR_INDICATOR_TYPES.state) {
    const state = isNumber(entity.state) ? entity.state : entity.state?.val;

    return getAlarmStateColor(state);
  }

  return getAlarmImpactStateColor(entity.impact_state);
};

/**
 * Get color for entity event
 *
 * @param {string} type
 */
export const getEntityEventColor = type => EVENT_ENTITY_COLORS_BY_TYPE[type];
