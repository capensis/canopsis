import {
  formToTicketStatusJob,
  ticketStatusJobToForm,
} from '@/helpers/entities/ticket-status-job/form';

describe('helpers/entities/ticket-status-job/form', () => {
  const job = {
    rule_name: 'Rule',
    rule_type: 1,
    ticket_system_name: 'Jira',
    ticket_id: '123',
    check_ticket_status: {
      request: {},
      status_mapping: { Open: 1 },
      ticket_status_tpl: '{{ .ticket.status }}',
    },
  };

  it('disables Go templates in the edit form', () => {
    const form = ticketStatusJobToForm(job);

    expect(form.check_ticket_status.ticket_status).toEqual({
      template: false,
      value: '{{ .ticket.status }}',
    });
  });

  it('always serializes the ticket status source as a plain field', () => {
    const form = ticketStatusJobToForm(job);
    const result = formToTicketStatusJob(form);

    expect(result.check_ticket_status.ticket_status).toBe('{{ .ticket.status }}');
    expect(result.check_ticket_status.ticket_status_tpl).toBe('');
  });
});
