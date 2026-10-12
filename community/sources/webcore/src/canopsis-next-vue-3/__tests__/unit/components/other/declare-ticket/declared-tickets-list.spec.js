import { generateShallowRenderer } from '@unit/utils/vue';

import DeclaredTicketsList from '@/components/other/declare-ticket/declared-tickets-list.vue';

const stubs = {
  'c-advanced-data-table': {
    props: ['headers', 'items'],
    template: `
      <div>
        <slot
          v-for="item in items"
          name="actions"
          :item="item"
        />
      </div>
    `,
  },
  'c-action-btn': {
    emits: ['click'],
    template: '<button class="remove-ticket" @click="$emit(\'click\', $event)" />',
  },
  'declare-ticket-rule-ticket-status-chip': true,
};

describe('declared-tickets-list', () => {
  const tickets = [{ ticket: 'TICKET-123' }];
  const factory = generateShallowRenderer(DeclaredTicketsList, { stubs });

  test('does not show ticket removal actions without access', () => {
    const wrapper = factory({ props: { tickets } });

    expect(wrapper.find('.remove-ticket').exists()).toBe(false);
    expect(wrapper.vm.headers).not.toEqual(
      expect.arrayContaining([expect.objectContaining({ value: 'actions' })]),
    );
  });

  test('emits the selected ticket from the remove action', async () => {
    const wrapper = factory({ props: { tickets, removable: true } });

    await wrapper.find('.remove-ticket').trigger('click');

    expect(wrapper.emitted('remove')).toEqual([[tickets[0]]]);
    expect(wrapper.vm.headers).toEqual(
      expect.arrayContaining([expect.objectContaining({ value: 'actions', sortable: false })]),
    );
  });
});
