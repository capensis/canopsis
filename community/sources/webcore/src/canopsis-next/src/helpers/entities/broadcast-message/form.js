import { DEFAULT_BROADCAST_MESSAGE_COLOR, BROADCAST_MESSAGE_VIEWS } from '@/constants';

import { convertDateToDateObject, convertDateToTimestamp } from '@/helpers/date/date';

/**
 * @typedef {Object} Broadcast
 * @property {string} message
 * @property {string} color
 * @property {number} start
 * @property {number} end
 * @property {string[]} views
 */

/**
 * @typedef {Broadcast} BroadcastForm
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
  views: [...(broadcastMessage.views || Object.values(BROADCAST_MESSAGE_VIEWS))],
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
