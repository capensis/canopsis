import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import AssociateTicketEventForm from '@/components/other/declare-ticket/form/associate-ticket-event-form.vue';

const stubs = {
  'c-information-block': true,
  'c-name-field': true,
  'declare-ticket-rule-ticket-id-field': true,
  'declare-ticket-rule-ticket-url-text-field': true,
  'declare-ticket-rule-ticket-custom-fields-field': true,
};

const selectSystemNameField = wrapper => wrapper.find('c-name-field-stub');
const selectDeclareTicketRuleTicketIdField = wrapper => wrapper.find('declare-ticket-rule-ticket-id-field-stub');
const selectDeclareTicketRuleTicketUrlTextField = wrapper => wrapper.find('declare-ticket-rule-ticket-url-text-field-stub');
const selectDeclareTicketRuleTicketCustomFieldsField = wrapper => wrapper.find('declare-ticket-rule-ticket-custom-fields-field-stub');

describe('associate-ticket-event-form', () => {
  const form = {
    system_name: 'System name',
    ticket_id: 'Ticket ID',
    ticket_url: 'Ticket URL',
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

    selectSystemNameField(wrapper).vm.$emit('input', newName);

    expect(wrapper).toEmit('input', {
      ...form,
      system_name: newName,
    });
  });

  test('Ticket id changed after trigger ticket id field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newTicketId = Faker.datatype.string();

    selectDeclareTicketRuleTicketIdField(wrapper).vm.$emit('input', newTicketId);

    expect(wrapper).toEmit('input', {
      ...form,
      ticket_id: newTicketId,
    });
  });

  test('Ticket url changed after trigger ticket url field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newTicketUrl = Faker.datatype.string();

    selectDeclareTicketRuleTicketUrlTextField(wrapper).vm.$emit('input', newTicketUrl);

    expect(wrapper).toEmit('input', {
      ...form,
      ticket_url: newTicketUrl,
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

    selectDeclareTicketRuleTicketCustomFieldsField(wrapper).vm.$emit('input', newMapping);

    expect(wrapper).toEmit('input', {
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

    expect(wrapper.element).toMatchSnapshot();
  });
});
