/**
 * @typedef {Object} DeclareTicketEvent
 * @property {string} [_id]
 * @property {string[]} [alarms]
 * @property {string} comment
 */

/**
 * @typedef {DeclareTicketEvent} DeclareTicketEventForm
 * @property {tickets_by_alarms} Object.<string, string[]>
 */

/**
 * Convert declare ticket event object to form compatible object
 *
 * @param {DeclareTicketEvent} [declareTicketEvent = {}]
 * @param {Alarm[]} [alarms = []]
 * @return {DeclareTicketEventForm}
 */
export const declareTicketEventToForm = (declareTicketEvent = {}, alarms = []) => ({
  tickets_by_alarms: alarms.reduce((acc, { _id: id }) => {
    acc[id] = [];

    return acc;
  }, {}),
  comment: declareTicketEvent.comment ?? '',
});

/**
 * Convert form object to declare ticket API compatible object
 *
 * @param {DeclareTicketEventForm} form
 * @return {DeclareTicketEvent[]}
 */
export const formToDeclareTicketEvents = (form) => {
  const eventsByTickets = Object.entries(form.tickets_by_alarms)
    .reduce((acc, [alarmId, tickets]) => {
      tickets.forEach((ticket) => {
        if (acc[ticket]) {
          acc[ticket].push(alarmId);
        } else {
          acc[ticket] = [alarmId];
        }
      });

      return acc;
    }, {});

  return Object.entries(eventsByTickets).map(([ticketId, alarms]) => ({
    _id: ticketId,
    alarms,
    comment: form.comment,
  }));
};
