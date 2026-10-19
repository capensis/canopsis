import { createFormStub } from '@unit/stubs/form';
import { createModalWrapperStub } from '@unit/stubs/modal';
import { mockModals } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';

import CreateRemoveAssociatedTicketEvent
  from '@/components/modals/alarm/create-remove-associated-ticket-event.vue';

const stubs = {
  'c-modal-wrapper': createModalWrapperStub('c-modal-wrapper'),
  'alarm-general-table': true,
  'c-collapse-panel': true,
  'remove-associated-ticket-event-form': true,
  'v-form': createFormStub('v-form'),
};

describe('create-remove-associated-ticket-event', () => {
  const $modals = mockModals();
  const factory = generateShallowRenderer(CreateRemoveAssociatedTicketEvent, {
    stubs,
    mocks: { $modals },
  });

  test('preselects and locks a ticket provided by a declared-ticket row', async () => {
    const action = vi.fn();
    const ticket = 'TICKET-123';
    const wrapper = factory({
      props: {
        modal: {
          id: 'remove-ticket-modal',
          config: {
            action,
            items: [{ _id: 'alarm-id', v: { tickets: [{ ticket }] } }],
            ticket,
            ticketDisabled: true,
          },
        },
      },
    });

    const form = wrapper.find('remove-associated-ticket-event-form-stub');

    expect(wrapper.vm.values).toEqual({ ticket, comment: '' });
    expect(form.attributes('ticket-disabled')).toBe('true');

    await wrapper.find('form').trigger('submit');

    expect(action).toHaveBeenCalledWith({ ticket, comment: '' });
    expect($modals.hide).toHaveBeenCalledWith({ id: 'remove-ticket-modal' });
  });
});
