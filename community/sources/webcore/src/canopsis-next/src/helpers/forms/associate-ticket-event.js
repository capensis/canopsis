import { objectToTextPairs, textPairsToObject } from '@/helpers/text-pairs';

/**
 * @typedef {Object} AssociateTicketEvent
 * @property {string} ticket
 * @property {string} ticket_url
 * @property {string} ticket_system_name
 * @property {Object} ticket_data
 * @property {string} output
 */

/**
 * @typedef {AssociateTicketEvent} AssociateTicketEventForm
 * @property {string} ticket_id
 * @property {string} ticket_url
 * @property {string} system_name
 * @property {TextPairObject[]} mapping
 * @property {string} output
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
  system_name: event.ticket_system_name ?? '',
  output: event.output ?? '',
  mapping: objectToTextPairs(event.ticket_data),
});

/**
 * Convert form object to associate ticket API compatible object
 *
 * @param {AssociateTicketEventForm} form
 * @return {AssociateTicketEvent}
 */
export const formToAssociateTicketEvent = form => ({
  ticket: form.ticket_id,
  ticket_url: form.ticket_url,
  ticket_system_name: form.system_name,
  output: form.output,
  ticket_data: textPairsToObject(form.mapping),
});
