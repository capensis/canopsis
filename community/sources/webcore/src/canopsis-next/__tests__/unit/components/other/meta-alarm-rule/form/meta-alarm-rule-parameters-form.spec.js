import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import { META_ALARMS_RULE_TYPES } from '@/constants';

import MetaAlarmRuleParametersForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-parameters-form.vue';

const stubs = {
  'c-form-block': true,
  'meta-alarm-rule-corel-form': {
    name: 'MetaAlarmRuleCorelForm',
    emits: ['update:config'],
    template: '<div />',
  },
  'meta-alarm-rule-threshold-field': {
    name: 'MetaAlarmRuleThresholdField',
    emits: ['update:config'],
    template: '<div />',
  },
  'meta-alarm-rule-time-based-field': {
    name: 'MetaAlarmRuleTimeBasedField',
    emits: ['update:modelValue'],
    template: '<div />',
  },
  'meta-alarm-rule-child-inactive-delay-field': true,
  'meta-alarm-value-paths-field': {
    name: 'MetaAlarmValuePathsField',
    emits: ['update:items'],
    template: '<div />',
  },
  'meta-alarm-rule-patterns-form': {
    name: 'MetaAlarmRulePatternsForm',
    emits: ['update:form'],
    template: '<div />',
  },
};

const selectMetaAlarmRuleCorelForm = wrapper =>
  wrapper.findComponent({ name: 'MetaAlarmRuleCorelForm' });
const selectMetaAlarmRuleThresholdField = wrapper =>
  wrapper.findComponent({ name: 'MetaAlarmRuleThresholdField' });
const selectMetaAlarmRuleTimeBasedField = wrapper =>
  wrapper.findComponent({ name: 'MetaAlarmRuleTimeBasedField' });
const selectMetaAlarmRuleValuePathsField = wrapper =>
  wrapper.findComponent({ name: 'MetaAlarmValuePathsField' });
const selectMetaAlarmRulePatternsForm = wrapper =>
  wrapper.findComponent({ name: 'MetaAlarmRulePatternsForm' });

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

  test('Corel changed after trigger corel form', async () => {
    const corelForm = {
      ...form,
      type: META_ALARMS_RULE_TYPES.corel,
    };
    const wrapper = factory({
      props: {
        form: corelForm,
      },
    });

    const metaAlarmRuleCorelForm = selectMetaAlarmRuleCorelForm(wrapper);

    const newCorelConfig = {
      corel_id: Faker.datatype.string(),
    };

    metaAlarmRuleCorelForm.vm.$emit('update:config', newCorelConfig);
    await wrapper.vm.$nextTick();

    expect(globalThis.__TEST_FORM__.getValue('config')).toEqual(newCorelConfig);
  });

  test('Threshold changed after trigger threshold form', async () => {
    const wrapper = factory({
      props: {
        form,
      },
    });

    const metaAlarmRuleThresholdField = selectMetaAlarmRuleThresholdField(wrapper);

    const newThresholdConfig = {
      threshold_count: 0.2,
    };

    metaAlarmRuleThresholdField.vm.$emit('update:config', newThresholdConfig);
    await wrapper.vm.$nextTick();

    expect(globalThis.__TEST_FORM__.getValue('config')).toEqual(newThresholdConfig);
  });

  test('Time based changed after trigger time based field', async () => {
    const wrapper = factory({
      props: {
        form,
      },
    });

    const metaAlarmRuleTimeBasedField = selectMetaAlarmRuleTimeBasedField(wrapper);

    const newTimeInterval = {};

    metaAlarmRuleTimeBasedField.vm.$emit('update:modelValue', newTimeInterval);
    await wrapper.vm.$nextTick();

    expect(globalThis.__TEST_FORM__.getValue('config.time_interval')).toEqual(newTimeInterval);
  });

  test('Value paths changed after trigger value paths field', async () => {
    const valuegroupForm = {
      ...form,
      type: META_ALARMS_RULE_TYPES.valuegroup,
    };
    const wrapper = factory({
      props: {
        form: valuegroupForm,
      },
    });

    const metaAlarmRuleValuePathsField = selectMetaAlarmRuleValuePathsField(wrapper);

    const newValuePaths = [Faker.datatype.string()];

    metaAlarmRuleValuePathsField.vm.$emit('update:items', newValuePaths);
    await wrapper.vm.$nextTick();

    expect(globalThis.__TEST_FORM__.getValue('config.value_paths')).toEqual(newValuePaths);
  });

  test('Patterns changed after trigger patterns field', async () => {
    const wrapper = factory({
      props: {
        form,
      },
    });

    const patternsField = selectMetaAlarmRulePatternsForm(wrapper);

    const newPatterns = {
      alarm_pattern: {
        id: Faker.datatype.string(),
      },
    };

    patternsField.vm.$emit('update:form', newPatterns);
    await wrapper.vm.$nextTick();

    expect(globalThis.__TEST_FORM__.getValue('patterns')).toEqual(newPatterns);
  });
});
