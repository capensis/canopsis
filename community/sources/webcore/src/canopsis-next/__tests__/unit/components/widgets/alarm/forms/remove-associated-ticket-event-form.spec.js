import { generateShallowRenderer } from '@unit/utils/vue';

import RemoveAssociatedTicketEventForm
  from '@/components/widgets/alarm/forms/remove-associated-ticket-event-form.vue';

const stubs = {
  'c-lazy-search-field': true,
  'c-description-field': true,
};

describe('remove-associated-ticket-event-form', () => {
  const factory = generateShallowRenderer(RemoveAssociatedTicketEventForm, { stubs });

  test('disables the ticket selector when requested', () => {
    const wrapper = factory({
      props: {
        form: { ticket: 'TICKET-123', comment: '' },
        items: [{
          v: {
            tickets: [
              { ticket: 'TICKET-123', ticket_system_name: 'Jira' },
              { ticket: 'TICKET-456', ticket_system_name: 'Jira' },
            ],
          },
        }],
        ticketDisabled: true,
      },
    });

    expect(wrapper.find('c-lazy-search-field-stub').attributes('disabled')).toBe('true');
  });
});
