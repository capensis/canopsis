/**
 * @typedef {Object} DeclareTicketEvent
 * @property {string} [_id]
 * @property {string[]} [alarms]
 * @property {string} comment
 * @property {boolean} ticket_resources
 */

/**
 * @typedef {DeclareTicketEvent[]} DeclareTicketEvents
 */

/**
 * @typedef {DeclareTicketEvent} DeclareTicketEventForm
 * @property {Object.<string, string[]>} alarms_by_tickets
 * @property {Object.<string, string>} comments_by_tickets
 * @property {Object.<string, boolean>} ticket_resources_by_tickets
 */

/**
 * @typedef {DeclareTicketEventForm[]} DeclareTicketEventsForm
 */

/**
 * Convert declare ticket event object to form compatible object
 *
 * @param {Object.<string, { name: string, alarms: string[] }>} [alarmIdsByTickets = {}]
 * @return {DeclareTicketEventForm}
 */
export const alarmsToDeclareTicketEventForm = (alarmIdsByTickets = {}) => Object.keys(alarmIdsByTickets)
  .reduce((acc, ticketId) => {
    acc.alarms_by_tickets[ticketId] = [];
    acc.comments_by_tickets[ticketId] = '';
    acc.ticket_resources_by_tickets[ticketId] = false;

    return acc;
  }, {
    alarms_by_tickets: {},
    comments_by_tickets: {},
    ticket_resources_by_tickets: {},
  });

/**
 * Convert form object to declare ticket API compatible object
 *
 * @param {DeclareTicketEventForm} form
 * @return {DeclareTicketEvent[]}
 */
export const formToDeclareTicketEvents = form => Object.entries(form.alarms_by_tickets)
  .reduce((acc, [ticketId, alarms]) => {
    if (alarms.length) {
      acc.push({
        _id: ticketId,
        alarms,
        comment: form.comments_by_tickets[ticketId],
        ticket_resources: form.ticket_resources_by_tickets[ticketId],
      });
    }

    return acc;
  }, []);
