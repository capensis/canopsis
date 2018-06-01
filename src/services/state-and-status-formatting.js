import { STATES_CHIPS_AND_FLAGS_STYLE, STATUS_CHIPS_AND_FLAGS_STYLE } from '@/config';

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
  if (isStatus && STATUS_CHIPS_AND_FLAGS_STYLE[value]) {
    if (isCroppedState) {
      return ({
        ...STATUS_CHIPS_AND_FLAGS_STYLE[value],
        icon: 'vertical_align_center',
      });
    }

    return STATUS_CHIPS_AND_FLAGS_STYLE[value];
  }

  if (!isStatus && STATES_CHIPS_AND_FLAGS_STYLE[value]) {
    return STATES_CHIPS_AND_FLAGS_STYLE[value];
  }

  return ({
    color: 'black',
    text: 'Invalid val',
    icon: 'clear',
  });
}
