import { generateShallowRenderer } from '@unit/utils/vue';

import TicketFormatsForm from '@/components/forms/ticket-formats-form.vue';

const stubs = {
  'c-card-iterator-form': true,
  'c-card-iterator-item': true,
  'c-text-field': true,
};

describe('ticket-formats-form', () => {
  const factory = generateShallowRenderer(TicketFormatsForm, { stubs });

  test('uses stable form keys as draggable item identities', () => {
    const wrapper = factory({
      propsData: {
        modelValue: [
          { key: 'first-format', value: '^INC-[0-9]+$' },
          { key: 'second-format', value: '^REQ-[0-9]+$' },
        ],
      },
    });

    expect(wrapper.findComponent({ name: 'CCardIteratorForm' }).attributes('item-key')).toBe('key');
  });
});
