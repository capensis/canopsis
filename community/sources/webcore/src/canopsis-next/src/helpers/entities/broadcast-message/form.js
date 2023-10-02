import { DEFAULT_BROADCAST_MESSAGE_COLOR } from '@/constants';

import { convertDateToDateObject, convertDateToTimestamp } from '@/helpers/date/date';

/**
 * @typedef {Object} Broadcast
 * @property {string} message
 * @property {string} color
 * @property {number} start
 * @property {number} end
 */

/**
 * @typedef {Object} BroadcastForm
 * @property {Date} start
 * @property {Date} end
 */

/**
 * Convert broadcast object to broadcast form
 *
 * @param {Broadcast} broadcastMessage
 * @return {BroadcastForm}
 */
export const messageToForm = (broadcastMessage = {}) => ({
  message: broadcastMessage.message || '',
  color: broadcastMessage.color || DEFAULT_BROADCAST_MESSAGE_COLOR,
  start: convertDateToDateObject(broadcastMessage.start),
  end: convertDateToDateObject(broadcastMessage.end),
});

/**
 * Convert broadcast form to broadcast object
 *
 * @param {BroadcastForm} form
 * @return {Broadcast}
 */
export const formToMessage = (form = {}) => ({
  ...form,
  start: convertDateToTimestamp(form.start),
  end: convertDateToTimestamp(form.end),
});
