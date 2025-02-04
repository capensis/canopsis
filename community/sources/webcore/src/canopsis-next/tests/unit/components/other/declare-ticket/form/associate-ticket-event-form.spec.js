import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import AssociateTicketEventForm from '@/components/other/declare-ticket/form/associate-ticket-event-form.vue';

const stubs = {
  'c-information-block': true,
  'c-name-field': true,
  'declare-ticket-rule-ticket-id-text-field': true,
  'declare-ticket-rule-ticket-url-text-field': true,
  'declare-ticket-rule-ticket-url-title-field': true,
  'declare-ticket-rule-ticket-custom-fields-field': true,
};

const selectSystemNameField = wrapper => wrapper.find('c-name-field-stub');
const selectDeclareTicketRuleTicketIdField = wrapper => wrapper.find('declare-ticket-rule-ticket-id-text-field-stub');
const selectDeclareTicketRuleTicketUrlTextField = wrapper => wrapper.find('declare-ticket-rule-ticket-url-text-field-stub');
const selectDeclareTicketRuleTicketUrlTitleField = wrapper => wrapper.find('declare-ticket-rule-ticket-url-title-field-stub');
const selectDeclareTicketRuleTicketCustomFieldsField = wrapper => wrapper.find('declare-ticket-rule-ticket-custom-fields-field-stub');

describe('associate-ticket-event-form', () => {
  const form = {
    system_name: 'System name',
    ticket: 'Ticket ID',
    ticket_url: 'Ticket URL',
    ticket_url_title: 'Ticket URL',
    output: 'Output',
    mapping: [
      {
        value: 'value',
        text: 'text',
      },
    ],
  };

  const factory = generateShallowRenderer(AssociateTicketEventForm, { stubs });
  const snapshotFactory = generateRenderer(AssociateTicketEventForm, { stubs });

  test('System name changed after trigger name field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newName = Faker.datatype.string();

    selectSystemNameField(wrapper).triggerCustomEvent('input', newName);

    expect(wrapper).toEmitInput({
      ...form,
      ticket_system_name: newName,
    });
  });

  test('Ticket changed after trigger ticket id field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newTicket = Faker.datatype.string();

    selectDeclareTicketRuleTicketIdField(wrapper).triggerCustomEvent('input', newTicket);

    expect(wrapper).toEmitInput({
      ...form,
      ticket: newTicket,
    });
  });

  test('Ticket url changed after trigger ticket url field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newTicketUrl = Faker.datatype.string();

    selectDeclareTicketRuleTicketUrlTextField(wrapper).triggerCustomEvent('input', newTicketUrl);

    expect(wrapper).toEmitInput({
      ...form,
      ticket_url: newTicketUrl,
    });
  });

  test('Ticket url title changed after trigger ticket url field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newTicketUrlTitle = Faker.datatype.string();

    selectDeclareTicketRuleTicketUrlTitleField(wrapper).triggerCustomEvent('input', newTicketUrlTitle);

    expect(wrapper).toEmitInput({
      ...form,
      ticket_url_title: newTicketUrlTitle,
    });
  });

  test('Mapping changed after trigger custom fields field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newMapping = [
      ...form.mapping,
      {
        value: Faker.datatype.string(),
        text: Faker.datatype.string(),
      },
    ];

    selectDeclareTicketRuleTicketCustomFieldsField(wrapper).triggerCustomEvent('input', newMapping);

    expect(wrapper).toEmitInput({
      ...form,
      mapping: newMapping,
    });
  });

  test('Renders `associate-ticket-event-form` with form', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
