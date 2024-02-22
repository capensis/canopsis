import Faker from 'faker';

import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';

import { REQUEST_METHODS } from '@/constants';

import PbehaviorGeneralForm from '@/components/other/declare-ticket/form/declare-ticket-rule-general-form.vue';

const stubs = {
  'c-enabled-field': true,
  'c-name-field': true,
  'declare-ticket-rule-webhooks-field': true,
};

const selectEnabledFields = wrapper => wrapper.findAll('c-enabled-field-stub');
const selectEnabledField = wrapper => selectEnabledFields(wrapper).at(0);
const selectEmitTriggerField = wrapper => selectEnabledFields(wrapper).at(1);
const selectNameFields = wrapper => wrapper.findAll('c-name-field-stub');
const selectNameField = wrapper => selectNameFields(wrapper).at(0);
const selectSystemNameField = wrapper => selectNameFields(wrapper).at(1);
const selectDeclareTicketRuleWebhooksField = wrapper => wrapper.find('declare-ticket-rule-webhooks-field-stub');

describe('declare-ticket-rule-general-form', () => {
  const form = {
    name: Faker.datatype.string(),
    system_name: Faker.datatype.string(),
    enabled: Faker.datatype.boolean(),
    emit_trigger: Faker.datatype.boolean(),
    webhooks: [],
  };
  const factory = generateShallowRenderer(PbehaviorGeneralForm, { stubs });
  const snapshotFactory = generateRenderer(PbehaviorGeneralForm, { stubs });

  test('Name changed after trigger name field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newName = Faker.datatype.string();

    selectNameField(wrapper).triggerCustomEvent('input', newName);

    expect(wrapper).toEmitInput({ ...form, name: newName });
  });

  test('System name changed after trigger name field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newName = Faker.datatype.string();

    selectSystemNameField(wrapper).triggerCustomEvent('input', newName);

    expect(wrapper).toEmitInput({ ...form, system_name: newName });
  });

  test('Enabled changed after trigger enabled field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newEnabled = !form.enabled;

    selectEnabledField(wrapper).triggerCustomEvent('input', newEnabled);

    expect(wrapper).toEmitInput({ ...form, enabled: newEnabled });
  });

  test('Emit trigger field changed after trigger enabled field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newEnabled = !form.emit_trigger;

    selectEmitTriggerField(wrapper).triggerCustomEvent('input', newEnabled);

    expect(wrapper).toEmitInput({ ...form, emit_trigger: newEnabled });
  });

  test('Webhooks field changed after trigger webhooks field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newWebhooks = [{
      method: REQUEST_METHODS.post,
      url: 'url',
    }];

    selectDeclareTicketRuleWebhooksField(wrapper).triggerCustomEvent('input', newWebhooks);

    expect(wrapper).toEmitInput({ ...form, webhooks: newWebhooks });
  });

  test('Renders `declare-ticket-rule-general-form` with empty form', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          name: '',
          system_name: '',
          enabled: true,
          emit_trigger: false,
          webhooks: [],
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `declare-ticket-rule-general-form` with filled form', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          name: 'name',
          system_name: 'system-name',
          enabled: true,
          emit_trigger: false,
          webhooks: [{
            method: REQUEST_METHODS.post,
            url: 'url',
          }],
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
