import { generateRenderer } from '@unit/utils/vue';

import DeclareTicketRuleTicketIdTextField
  from '@/components/other/declare-ticket/form/fields/declare-ticket-rule-ticket-id-text-field.vue';

const stubs = {
  'c-text-field': {
    template: '<div><slot name="append-inner" /></div>',
  },
  'c-help-icon': {
    template: '<div class="help-tooltip"><slot /></div>',
  },
};

describe('declare-ticket-rule-ticket-id-text-field', () => {
  const factory = generateRenderer(DeclareTicketRuleTicketIdTextField, { stubs });

  test('shows every configured ticket format on a separate line in the help tooltip', () => {
    const ticketFormats = ['^INC-[0-9]+$', '^REQ-[A-Z]+$'];
    const wrapper = factory({
      propsData: { ticketFormats },
    });
    const formats = wrapper.findAll('.help-tooltip code');

    expect(wrapper.find('.help-tooltip').text()).toContain('declareTicket.responseFieldHelpText');
    expect(wrapper.find('.help-tooltip').text()).toContain('declareTicket.allowedTicketFormats');
    expect(formats).toHaveLength(ticketFormats.length);
    expect(formats.map(format => format.text())).toEqual(ticketFormats);
    expect(formats.every(format => format.classes('d-block'))).toBe(true);
  });
});
