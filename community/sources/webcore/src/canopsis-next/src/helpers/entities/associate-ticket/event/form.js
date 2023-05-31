import { objectToTextPairs, textPairsToObject } from '@/helpers/text-pairs';

/**
 * @typedef {Object} AssociateTicketEvent
 * @property {string} ticket
 * @property {string} ticket_url
 * @property {string} ticket_system_name
 * @property {Object} ticket_data
 * @property {string} comment
 * @property {boolean} ticket_resources
 */

/**
 * @typedef {AssociateTicketEvent} AssociateTicketEventForm
 * @property {string} ticket_id
 * @property {string} system_name
 * @property {TextPairObject[]} mapping
 */

/**
 * Convert associate ticket event object to form compatible object
 *
 * @param {AssociateTicketEvent} [event = {}]
 * @return {AssociateTicketEventForm}
 */
export const eventToAssociateTicketForm = (event = {}) => ({
  ticket_id: event.ticket ?? '',
  ticket_url: event.ticket_url ?? '',
  ticket_resources: event.ticket_resources ?? false,
  system_name: event.ticket_system_name ?? '',
  comment: event.comment ?? '',
  mapping: objectToTextPairs(event.ticket_data),
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
  system_name: form.system_name,
  ticket: form.ticket_id,
  ticket_resources: form.ticket_resources,
  url: form.ticket_url,
});
