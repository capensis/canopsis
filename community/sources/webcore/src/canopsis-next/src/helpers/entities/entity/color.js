import { isNumber } from 'lodash';

import { CSS_COLORS_VARS } from '@/config';
import { COLOR_INDICATOR_TYPES, EVENT_ENTITY_COLORS_BY_TYPE, PBEHAVIOR_CANONICAL_TYPES } from '@/constants';

import { getAlarmImpactStateColor, getAlarmStateColor } from '../alarm/color';

/**
 * Get color for a entity by colorIndicator and isGrey parameters
 *
 * @param {Service | Entity | {}} [entity = {}]
 * @param {string} [colorIndicator = COLOR_INDICATOR_TYPES.state]
 * @returns {string|*}
 */
export const getEntityColor = (
  {
    state,
    is_grey: isGrey,
    pbehavior_info: pbehaviorInfo,
    impact_state: impactState,
  } = {},
  colorIndicator = COLOR_INDICATOR_TYPES.state,
) => {
  if (
    isGrey
    || (pbehaviorInfo?.canonical_type && pbehaviorInfo?.canonical_type !== PBEHAVIOR_CANONICAL_TYPES.active)
  ) {
    return CSS_COLORS_VARS.state.pause;
  }

  if (colorIndicator === COLOR_INDICATOR_TYPES.state) {
    return getAlarmStateColor(isNumber(state) ? state : state?.val);
  }

  return getAlarmImpactStateColor(impactState);
};

/**
 * Get color for entity event
 *
 * @param {string} type
 */
export const getEntityEventColor = type => EVENT_ENTITY_COLORS_BY_TYPE[type];
