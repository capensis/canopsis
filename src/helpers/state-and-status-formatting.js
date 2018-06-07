import { ENTITIES_STATES_STYLES, ENTITY_STATUS_STYLES } from '@/constants';

/**
 * Format state and status on list-alarm with color, text and icon.
 *
 * @param {Number} value => Value of the state/status
 * @param {Boolean} isStatus => Is it a status, or a state
 * @param {Boolean} isCroppedState => Is this state a croppedState
 *
 * @returns {Object}
 */
export default function formatStateAndStatus(value, isStatus, isCroppedState) {
  if (isStatus && ENTITY_STATUS_STYLES[value]) {
    if (isCroppedState) {
      return {
        ...ENTITY_STATUS_STYLES[value],
        icon: 'vertical_align_center',
      };
    }
    return ENTITY_STATUS_STYLES[value];
  }

  if (!isStatus && ENTITIES_STATES_STYLES[value]) {
    return ENTITIES_STATES_STYLES[value];
  }

  return {
    color: 'black',
    text: 'Invalid val',
    icon: 'clear',
  };
}
