import { generateShallowRenderer } from '@unit/utils/vue';

import { CUSTOM_WIDGET_TEMPLATE } from '@/constants';

import AssociateTicket from '@/components/sidebars/alarm/form/fields/associate-ticket.vue';

const stubs = {
  'widget-settings-group': {
    template: '<div><slot /></div>',
  },
  'v-container': {
    template: '<div><slot /></div>',
  },
  'c-widget-template-field': true,
  'ticket-formats-form': true,
};

const selectTicketFormatsForm = wrapper => wrapper.findComponent({ name: 'TicketFormatsForm' });

describe('associate-ticket', () => {
  const factory = generateShallowRenderer(AssociateTicket, { stubs });

  test('shows the formats editor for an empty selected template', () => {
    const wrapper = factory({
      propsData: {
        modelValue: [],
        template: 'default-ticket-formats-template-id',
      },
    });

    expect(selectTicketFormatsForm(wrapper).exists()).toBe(true);
    expect(wrapper.text()).toContain('This template is empty');
  });

  test('switches to custom when selected template formats are changed', () => {
    const wrapper = factory({
      propsData: {
        modelValue: [],
        template: 'default-ticket-formats-template-id',
      },
    });
    const formats = [{ key: 'format-key', value: '^INC-[0-9]+$' }];

    selectTicketFormatsForm(wrapper).vm.$emit('update:modelValue', formats);

    expect(wrapper.emitted('update:template')).toEqual([[CUSTOM_WIDGET_TEMPLATE, formats]]);
  });
});
