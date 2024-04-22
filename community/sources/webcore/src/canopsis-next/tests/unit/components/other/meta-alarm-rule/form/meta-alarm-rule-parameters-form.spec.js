import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { META_ALARMS_RULE_TYPES } from '@/constants';

import MetaAlarmRuleParametersForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-parameters-form.vue';

const stubs = {
  'meta-alarm-rule-corel-form': true,
  'meta-alarm-rule-threshold-form': true,
  'meta-alarm-rule-time-based-form': true,
  'meta-alarm-rule-value-paths-form': true,
  'meta-alarm-rule-patterns-form': true,
};

const selectMetaAlarmRuleCorelForm = wrapper => wrapper.find('meta-alarm-rule-corel-form-stub');
const selectMetaAlarmRuleThresholdForm = wrapper => wrapper.find('meta-alarm-rule-threshold-form-stub');
const selectMetaAlarmRuleTimeBasedForm = wrapper => wrapper.find('meta-alarm-rule-time-based-form-stub');
const selectMetaAlarmRuleValuePathsForm = wrapper => wrapper.find('meta-alarm-rule-value-paths-form-stub');
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

    const metaAlarmRuleThresholdForm = selectMetaAlarmRuleThresholdForm(wrapper);

    const newThresholdConfig = {
      threshold_count: 0.2,
    };

    metaAlarmRuleThresholdForm.triggerCustomEvent('input', newThresholdConfig);

    expect(wrapper).toEmitInput({
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

    expect(wrapper).toEmitInput({
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

    expect(wrapper).toEmitInput({
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

    expect(wrapper).toEmitInput({
      ...form,
      patterns: newPatterns,
    });
  });

  test('Renders `meta-alarm-rule-parameters-form` with default props', () => {
    const wrapper = snapshotFactory();

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
    Object.values(META_ALARMS_RULE_TYPES),
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
