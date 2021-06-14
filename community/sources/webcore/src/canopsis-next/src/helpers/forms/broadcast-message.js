import moment from 'moment';

import { DEFAULT_BROADCAST_MESSAGE_COLOR } from '@/constants';

import { convertTimestampToMoment } from '@/helpers/date/date';

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
  start: broadcastMessage.start
    ? convertTimestampToMoment(broadcastMessage.start).toDate()
    : new Date(),
  end: broadcastMessage.end
    ? convertTimestampToMoment(broadcastMessage.end).toDate()
    : new Date(),
});

/**
 * Convert broadcast form to broadcast object
 *
 * @param {BroadcastForm} form
 * @return {Broadcast}
 */
export const formToMessage = (form = {}) => ({
  ...form,
  start: moment(form.start).unix(),
  end: moment(form.end).unix(),
});
