/**
 * @typedef {Object} TicketStatusJob
 * @property {string} [_id]
 * @property {string} [rule_type]
 * @property {string} [rule_name]
 * @property {string} [ticket_system_name]
 * @property {string} [ticket_id]
 * @property {string} [ticket]
 */

/**
 * @typedef {Object} TicketStatusJobForm
 * @property {string} rule_type
 * @property {string} rule_name
 * @property {string} ticket_system_name
 * @property {string} ticket_id
 */

/**
 * Convert ticket status job entity to form object
 *
 * @param {TicketStatusJob} job
 * @returns {TicketStatusJobForm}
 */
export const ticketStatusJobToForm = (job = {}) => ({
  rule_type: job.rule_type ?? '',
  rule_name: job.rule_name ?? '',
  ticket_system_name: job.ticket_system_name ?? '',
  ticket_id: job.ticket_id ?? job.ticket ?? '',
});

/**
 * Convert ticket status job form to entity
 *
 * @param {TicketStatusJobForm} form
 * @returns {TicketStatusJob}
 */
export const formToTicketStatusJob = form => ({
  rule_type: form.rule_type,
  rule_name: form.rule_name,
  ticket_system_name: form.ticket_system_name,
  ticket_id: form.ticket_id,
});
