import { objectToTextPairs, textPairsToObject } from '@/helpers/text-pairs';

/**
 * @typedef {Object} AssociateTicketEvent
 * @property {string} ticket
 * @property {string} url
 * @property {string} url_title
 * @property {string} system_name
 * @property {string} comment
 * @property {boolean} ticket_resources
 * @property {Object} data
 */

/**
 * @typedef {Object} AssociateTicketEventForm
 * @property {string} ticket
 * @property {string} ticket_url
 * @property {string} ticket_url_title
 * @property {string} ticket_system_name
 * @property {string} comment
 * @property {boolean} ticket_resources
 * @property {TextPairObject[]} mapping
 */

/**
 * Convert associate ticket event object to form compatible object
 *
 * @param {AssociateTicketEvent} [event = {}]
 * @return {AssociateTicketEventForm}
 */
export const eventToAssociateTicketForm = (event = {}) => ({
  ticket: event.ticket ?? '',
  ticket_url: event.url ?? '',
  ticket_url_title: event.url_title ?? '',
  ticket_resources: event.ticket_resources ?? false,
  ticket_system_name: event.system_name ?? '',
  comment: event.comment ?? '',
  mapping: objectToTextPairs(event.data),
});

/**
 * Convert form object to associate ticket API compatible object
 *
 * @param {AssociateTicketEventForm} form
 * @return {AssociateTicketEvent}
 */
export const formToAssociateTicketEvent = form => ({
  comment: form.comment,
  data: textPairsToObject(form.mapping),
  system_name: form.ticket_system_name,
  ticket: form.ticket,
  ticket_resources: form.ticket_resources,
  url: form.ticket_url,
  url_title: form.ticket_url_title,
});
