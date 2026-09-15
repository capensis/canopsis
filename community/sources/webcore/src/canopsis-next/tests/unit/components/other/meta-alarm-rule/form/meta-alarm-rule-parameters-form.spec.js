import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { META_ALARMS_RULE_TYPES } from '@/constants';

import MetaAlarmRuleParametersForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-parameters-form.vue';

const stubs = {
  'c-form-block': true,
  'c-form-block-row': true,
  'meta-alarm-rule-corel-form': true,
  'meta-alarm-rule-threshold-field': true,
  'meta-alarm-rule-time-based-field': true,
  'meta-alarm-rule-child-inactive-delay-field': true,
  'meta-alarm-rule-value-paths-field': true,
  'meta-alarm-rule-patterns-form': true,
};

const selectMetaAlarmRuleCorelForm = wrapper => wrapper.find('meta-alarm-rule-corel-form-stub');
const selectMetaAlarmRuleThresholdField = wrapper => wrapper.find('meta-alarm-rule-threshold-field-stub');
const selectMetaAlarmRuleTimeBasedField = wrapper => wrapper.find('meta-alarm-rule-time-based-field-stub');
const selectMetaAlarmRuleValuePathsField = wrapper => wrapper.find('meta-alarm-rule-value-paths-field-stub');
const selectMetaAlarmRulePatternsForm = wrapper => wrapper.find('meta-alarm-rule-patterns-form-stub');

describe('meta-alarm-rule-parameters-form', () => {
  const form = {
    _id: 'meta-alarm-rule-id',
    name: 'event-filter-name',
    output_template: 'event-filter-output-template',
    auto_resolve: true,
    type: META_ALARMS_RULE_TYPES.complex,
    config: {},
  };

  const factory = generateShallowRenderer(MetaAlarmRuleParametersForm, { stubs });
  const snapshotFactory = generateRenderer(MetaAlarmRuleParametersForm, { stubs });

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

    expect(wrapper).toEmitInput({
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

    const metaAlarmRuleThresholdField = selectMetaAlarmRuleThresholdField(wrapper);

    const newThresholdConfig = {
      threshold_count: 0.2,
    };

    metaAlarmRuleThresholdField.triggerCustomEvent('input', newThresholdConfig);

    expect(wrapper).toEmitInput({
      ...form,
      config: newThresholdConfig,
    });
  });

  test('Time based changed after trigger time based field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const metaAlarmRuleTimeBasedField = selectMetaAlarmRuleTimeBasedField(wrapper);

    const newTimeInterval = {};

    metaAlarmRuleTimeBasedField.triggerCustomEvent('input', newTimeInterval);

    expect(wrapper).toEmitInput({
      ...form,
      config: {
        ...form.config,
        time_interval: newTimeInterval,
      },
    });
  });

  test('Value paths changed after trigger value paths field', () => {
    const valuegroupForm = {
      ...form,
      type: META_ALARMS_RULE_TYPES.valuegroup,
    };
    const wrapper = factory({
      propsData: {
        form: valuegroupForm,
      },
    });

    const metaAlarmRuleValuePathsField = selectMetaAlarmRuleValuePathsField(wrapper);

    const newValuePaths = [Faker.datatype.string()];

    metaAlarmRuleValuePathsField.triggerCustomEvent('input', newValuePaths);

    expect(wrapper).toEmitInput({
      ...valuegroupForm,
      config: {
        ...valuegroupForm.config,
        value_paths: newValuePaths,
      },
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

    expect(wrapper).toEmitInput({
      ...form,
      patterns: newPatterns,
    });
  });

  test('Renders `meta-alarm-rule-parameters-form` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          type: META_ALARMS_RULE_TYPES.complex,
          config: {},
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `meta-alarm-rule-parameters-form` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test.each(
    Object.values(META_ALARMS_RULE_TYPES).filter(type => type !== META_ALARMS_RULE_TYPES.manualgroup),
  )('Renders `meta-alarm-rule-parameters-form` with `%s` type', (type) => {
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
