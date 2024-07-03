import { ALARM_LIST_STEPS } from '@/constants';

/**
 * Checks if the ticket is a successful declaration ticket or associated ticket.
 *
 * @param {Object} ticket - The ticket object to check.
 * @returns {boolean} - True if the ticket is a successful declaration ticket or associated ticket, false otherwise.
 */
export const isSuccessTicketDeclaration = ticket => [
  ALARM_LIST_STEPS.declareTicket,
  ALARM_LIST_STEPS.assocTicket,
].includes(ticket?._t);
