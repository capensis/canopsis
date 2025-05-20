import { isNumber } from 'lodash';

import { CSS_COLORS_VARS } from '@/config';
import {
  COLOR_INDICATOR_TYPES,
  EVENT_ENTITY_COLORS_BY_TYPE,
  PBEHAVIOR_CANONICAL_TYPES,
  ALARM_STATES_CLASSES,
} from '@/constants';

import { getAlarmImpactStateColor, getAlarmStateColor, getAlarmImpactStateGroupedColorIndex } from '../alarm/color';

/**
 * Determines if the entity should be considered in a pause (grey) state.
 *
 * @param {Object} pbehaviorInfo - The pbehavior_info object from the entity.
 * @param {boolean} isGrey - The is_grey property from the entity.
 * @returns {boolean} True if the entity is in a pause state, otherwise false.
 */
const isEntityPauseState = (pbehaviorInfo, isGrey) => (
  isGrey
  || (
    pbehaviorInfo?.canonical_type
    && pbehaviorInfo?.canonical_type !== PBEHAVIOR_CANONICAL_TYPES.active
  )
);

/**
 * Get color class for a entity by colorIndicator and isGrey parameters
 *
 * @param {Service | Entity | {}} [entity = {}]
 * @param {string} [colorIndicator = COLOR_INDICATOR_TYPES.state]
 * @returns {string|*}
 */
export const getEntityColorClass = (
  {
    state,
    is_grey: isGrey,
    pbehavior_info: pbehaviorInfo,
    impact_state: impactState,
  } = {},
  colorIndicator = COLOR_INDICATOR_TYPES.state,
) => {
  if (isEntityPauseState(pbehaviorInfo, isGrey)) {
    return 'state-pause';
  }

  if (colorIndicator === COLOR_INDICATOR_TYPES.state) {
    return ALARM_STATES_CLASSES[state?.val];
  }

  return `impact-state-${getAlarmImpactStateGroupedColorIndex(impactState)}`;
};

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
  if (isEntityPauseState(pbehaviorInfo, isGrey)) {
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
