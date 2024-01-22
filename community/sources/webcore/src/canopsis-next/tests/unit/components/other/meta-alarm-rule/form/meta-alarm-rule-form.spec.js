import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createInputStub, createSelectInputStub } from '@unit/stubs/input';
import { META_ALARMS_RULE_TYPES } from '@/constants';

import MetaAlarmRuleForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-form.vue';

const stubs = {
  'c-id-field': true,
  'c-name-field': true,
  'c-description-field': true,
  'c-enabled-field': true,
  'meta-alarm-rule-corel-form': true,
  'meta-alarm-rule-threshold-form': true,
  'meta-alarm-rule-time-based-form': true,
  'meta-alarm-rule-value-paths-form': true,
  'meta-alarm-rule-patterns-form': true,
  'v-text-field': createInputStub('v-text-field'),
  'v-select': createSelectInputStub('v-select'),
};

const snapshotStubs = {
  'c-id-field': true,
  'c-name-field': true,
  'c-description-field': true,
  'c-enabled-field': true,
  'meta-alarm-rule-corel-form': true,
  'meta-alarm-rule-threshold-form': true,
  'meta-alarm-rule-time-based-form': true,
  'meta-alarm-rule-value-paths-form': true,
  'meta-alarm-rule-patterns-form': true,
};

const selectIdField = wrapper => wrapper.find('c-id-field-stub');
const selectDescriptionField = wrapper => wrapper.find('c-description-field-stub');
const selectNameField = wrapper => wrapper.find('c-name-field-stub');
const selectEnabledField = wrapper => wrapper.find('c-enabled-field-stub');
const selectMetaAlarmRuleCorelForm = wrapper => wrapper.find('meta-alarm-rule-corel-form-stub');
const selectMetaAlarmRuleThresholdForm = wrapper => wrapper.find('meta-alarm-rule-threshold-form-stub');
const selectMetaAlarmRuleTimeBasedForm = wrapper => wrapper.find('meta-alarm-rule-time-based-form-stub');
const selectMetaAlarmRuleValuePathsForm = wrapper => wrapper.find('meta-alarm-rule-value-paths-form-stub');
const selectMetaAlarmRulePatternsForm = wrapper => wrapper.find('meta-alarm-rule-patterns-form-stub');
const selectMetaAlarmRuleTypeField = wrapper => wrapper.find('.v-select');

describe('meta-alarm-rule-form', () => {
  const form = {
    _id: 'meta-alarm-rule-id',
    name: 'event-filter-name',
    output_template: 'event-filter-output-template',
    auto_resolve: true,
    type: META_ALARMS_RULE_TYPES.complex,
    config: {},
  };

  const factory = generateShallowRenderer(MetaAlarmRuleForm, { stubs });
  const snapshotFactory = generateRenderer(MetaAlarmRuleForm, { stubs: snapshotStubs });

  test('ID changed after trigger id field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const idField = selectIdField(wrapper);

    const newId = Faker.datatype.string();

    idField.triggerCustomEvent('input', newId);

    expect(wrapper).toEmit('input', {
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

    expect(wrapper).toEmit('input', {
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

    expect(wrapper).toEmit('input', {
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

    expect(wrapper).toEmit('input', {
      ...form,
      auto_resolve: newAutoResolve,
    });
  });

  test('Type changed after trigger type field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const metaAlarmRuleTypeField = selectMetaAlarmRuleTypeField(wrapper);

    metaAlarmRuleTypeField.triggerCustomEvent('input', META_ALARMS_RULE_TYPES.attribute);

    expect(wrapper).toEmit('input', {
      ...form,
      type: META_ALARMS_RULE_TYPES.attribute,
    });
  });

  test('Corel changed after trigger corel form', () => {
    const corelForm = {
      ...form,
      type: META_ALARMS_RULE_TYPES.corel,
    };
    const wrapper = factory({
      propsData: {
        form: corelForm,
      },
    });

    const metaAlarmRuleCorelForm = selectMetaAlarmRuleCorelForm(wrapper);

    const newCorelConfig = {
      corel_id: Faker.datatype.string(),
    };

    metaAlarmRuleCorelForm.triggerCustomEvent('input', newCorelConfig);

    expect(wrapper).toEmit('input', {
      ...corelForm,
      config: newCorelConfig,
    });
  });

  test('Threshold changed after trigger threshold form', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const metaAlarmRuleThresholdForm = selectMetaAlarmRuleThresholdForm(wrapper);

    const newThresholdConfig = {
      threshold_count: 0.2,
    };

    metaAlarmRuleThresholdForm.triggerCustomEvent('input', newThresholdConfig);

    expect(wrapper).toEmit('input', {
      ...form,
      config: newThresholdConfig,
    });
  });

  test('Time based changed after trigger time based form', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const metaAlarmRuleTimeBasedForm = selectMetaAlarmRuleTimeBasedForm(wrapper);

    const newTimeBasedConfig = {
      time_interval: {},
    };

    metaAlarmRuleTimeBasedForm.triggerCustomEvent('input', newTimeBasedConfig);

    expect(wrapper).toEmit('input', {
      ...form,
      config: newTimeBasedConfig,
    });
  });

  test('Value paths changed after trigger value paths form', () => {
    const valuegroupForm = {
      ...form,
      type: META_ALARMS_RULE_TYPES.valuegroup,
    };
    const wrapper = factory({
      propsData: {
        form: valuegroupForm,
      },
    });

    const metaAlarmRuleValuePathsForm = selectMetaAlarmRuleValuePathsForm(wrapper);

    const newValuePathsConfig = {
      value_paths: [Faker.datatype.string()],
    };

    metaAlarmRuleValuePathsForm.triggerCustomEvent('input', newValuePathsConfig);

    expect(wrapper).toEmit('input', {
      ...valuegroupForm,
      config: newValuePathsConfig,
    });
  });

  test('Patterns changed after trigger patterns field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const patternsField = selectMetaAlarmRulePatternsForm(wrapper);

    const newPatterns = {
      alarm_pattern: {
        id: Faker.datatype.string(),
      },
    };

    patternsField.triggerCustomEvent('input', newPatterns);

    expect(wrapper).toEmit('input', {
      ...form,
      patterns: newPatterns,
    });
  });

  test('Renders `meta-alarm-rule-form` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `meta-alarm-rule-form` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test.each(
    Object.values(META_ALARMS_RULE_TYPES),
  )('Renders `meta-alarm-rule-form` with `%s` type', (type) => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          ...form,
          type,
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
