import {
  declareTicketRuleCheckTicketStatusToForm,
  formToDeclareTicketRuleCheckTicketStatus,
} from '@/helpers/entities/declare-ticket/rule/form';

/**
 * @typedef {Object} TicketStatusJob
 * @property {string} [rule_name]
 * @property {string} [rule_type]
 * @property {string} [ticket_id]
 * @property {string} [ticket_system_name]
 * @property {DeclareTicketRuleCheckTicketStatus} [check_ticket_status]
 */

/**
 * @typedef {Object} TicketStatusJobForm
 * @property {string} [rule_name]
 * @property {string} [rule_type]
 * @property {string} [ticket_system_name]
 * @property {string} [ticket_id]
 * @property {DeclareTicketRuleCheckTicketStatusForm} [check_ticket_status]
 */

/**
 * Convert ticket status job object to form compatible object.
 *
 * @param {TicketStatusJob} job
 * @returns {TicketStatusJobForm}
 */
export const ticketStatusJobToForm = (job = {}) => {
  const checkTicketStatus = declareTicketRuleCheckTicketStatusToForm(job.check_ticket_status);

  checkTicketStatus.ticket_status.template = false;

  return {
    rule_name: job.rule_name ?? '',
    rule_type: job.rule_type ?? '',
    ticket_system_name: job.ticket_system_name ?? '',
    ticket_id: job.ticket_id ?? job.ticket ?? '',
    check_ticket_status: checkTicketStatus,
  };
};

/**
 * Convert ticket status job form to API compatible object.
 *
 * @param {TicketStatusJobForm} form
 * @returns {TicketStatusJob}
 */
export const formToTicketStatusJob = form => ({
  rule_name: form.rule_name,
  rule_type: form.rule_type,
  ticket_system_name: form.ticket_system_name,
  ticket_id: form.ticket_id,
  check_ticket_status: formToDeclareTicketRuleCheckTicketStatus({
    ...form.check_ticket_status,
    ticket_status: {
      ...form.check_ticket_status.ticket_status,
      template: false,
    },
  }),
});
