import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { metaAlarmRuleToForm } from '@/helpers/entities/meta-alarm/rule/form';

import MetaAlarmRuleGeneralForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-general-form.vue';

const stubs = {
  'c-id-field': true,
  'c-name-field': true,
  'c-description-field': true,
  'c-enabled-field': true,
};

const selectIdField = wrapper => wrapper.find('c-id-field-stub');
const selectDescriptionField = wrapper => wrapper.find('c-description-field-stub');
const selectNameField = wrapper => wrapper.find('c-name-field-stub');
const selectEnabledField = wrapper => wrapper.find('c-enabled-field-stub');

describe('meta-alarm-rule-general-form', () => {
  const form = metaAlarmRuleToForm();

  const factory = generateShallowRenderer(MetaAlarmRuleGeneralForm, { stubs });
  const snapshotFactory = generateRenderer(MetaAlarmRuleGeneralForm, { stubs });

  test('ID changed after trigger id field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const idField = selectIdField(wrapper);

    const newId = Faker.datatype.string();

    idField.triggerCustomEvent('input', newId);

    expect(wrapper).toEmitInput({
      ...form,
      _id: newId,
    });
  });

  test('Output template changed after trigger description field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const descriptionField = selectDescriptionField(wrapper);

    const outputTemplate = Faker.datatype.string();

    descriptionField.triggerCustomEvent('input', outputTemplate);

    expect(wrapper).toEmitInput({
      ...form,
      output_template: outputTemplate,
    });
  });

  test('Name changed after trigger text field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const nameField = selectNameField(wrapper);

    const name = Faker.datatype.string();

    nameField.triggerCustomEvent('input', name);

    expect(wrapper).toEmitInput({
      ...form,
      name,
    });
  });

  test('Enabled changed after trigger enabled field', () => {
    const autoResolve = Faker.datatype.boolean();
    const wrapper = factory({
      propsData: {
        form: {
          ...form,
          auto_resolve: autoResolve,
        },
      },
    });

    const enabledField = selectEnabledField(wrapper);

    const newAutoResolve = !autoResolve;

    enabledField.triggerCustomEvent('input', newAutoResolve);

    expect(wrapper).toEmitInput({
      ...form,
      auto_resolve: newAutoResolve,
    });
  });

  test('Renders `meta-alarm-rule-general-form` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `meta-alarm-rule-general-form` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
