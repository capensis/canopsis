import {
  declareTicketRuleCheckTicketStatusToForm,
  formToDeclareTicketRuleCheckTicketStatus,
} from '@/helpers/entities/declare-ticket/rule/form';

/**
 * @typedef {Object} TicketStatusJob
 * @property {string} [ticket_id]
 * @property {string} [ticket_system_name]
 * @property {DeclareTicketRuleCheckTicketStatus} [check_ticket_status]
 */

/**
 * @typedef {Object} TicketStatusJobForm
 * @property {string} [ticket_system_name]
 * @property {string} [ticket_id]
 * @property {DeclareTicketRuleCheckTicketStatusForm} [check_ticket_status]
 */

/**
 * Convert ticket status job entity to form object
 *
 * @param {TicketStatusJob} job
 * @returns {TicketStatusJobForm}
 */
export const ticketStatusJobToForm = (job = {}) => ({
  ticket_system_name: job.ticket_system_name ?? '',
  ticket_id: job.ticket_id ?? job.ticket ?? '',
  check_ticket_status: declareTicketRuleCheckTicketStatusToForm(job.check_ticket_status),
});

/**
 * Convert ticket status job form to entity
 *
 * @param {TicketStatusJobForm} form
 * @returns {TicketStatusJob}
 */
export const formToTicketStatusJob = form => ({
  ticket_id: form.ticket_id,
  ticket_system_name: form.ticket_system_name,
  check_ticket_status: formToDeclareTicketRuleCheckTicketStatus(form.check_ticket_status),
});
