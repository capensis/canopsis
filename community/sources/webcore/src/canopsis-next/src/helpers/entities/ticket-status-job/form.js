import {
  declareTicketRuleCheckTicketStatusToForm,
  formToDeclareTicketRuleCheckTicketStatus,
} from '@/helpers/entities/declare-ticket/rule/form';

/**
 * Convert ticket status job object to form compatible object.
 *
 * @param {Object} job
 * @returns {Object}
 */
export const ticketStatusJobToForm = (job = {}) => ({
  rule_name: job.rule_name ?? '',
  rule_type: job.rule_type ?? '',
  ticket_system_name: job.ticket_system_name ?? '',
  ticket_id: job.ticket_id ?? job.ticket ?? '',
  check_ticket_status: declareTicketRuleCheckTicketStatusToForm(job.check_ticket_status),
});

/**
 * Convert ticket status job form to API compatible object.
 *
 * @param {Object} form
 * @returns {Object}
 */
export const formToTicketStatusJob = form => ({
  rule_name: form.rule_name,
  rule_type: form.rule_type,
  ticket_system_name: form.ticket_system_name,
  ticket_id: form.ticket_id,
  check_ticket_status: formToDeclareTicketRuleCheckTicketStatus(form.check_ticket_status),
});
