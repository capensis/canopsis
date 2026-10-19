import { generateShallowRenderer } from '@unit/utils/vue';

import { TEMPLATE_TESTING_TEST_TYPES } from '@/constants';

import { metaAlarmRuleToForm } from '@/helpers/entities/meta-alarm/rule/form';

import MetaAlarmRuleForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-form.vue';

const stubs = {
  'c-enabled-field': {
    name: 'CEnabledField',
    emits: ['update:modelValue'],
    template: '<div />',
  },
  'c-form-general-patterns-tabs': {
    name: 'CFormGeneralPatternsTabs',
    props: ['form', 'ruleId', 'type'],
    template: `
      <div>
        <slot name="general" :set-ref="() => {}" />
        <slot name="patterns" :set-ref="() => {}" />
      </div>
    `,
  },
  'meta-alarm-rule-general-form': {
    name: 'MetaAlarmRuleGeneralForm',
    props: ['disabledIdField'],
    template: '<div />',
  },
  'meta-alarm-rule-parameters-form': {
    name: 'MetaAlarmRuleParametersForm',
    template: '<div />',
  },
};

const selectEnabledField = wrapper => wrapper.findComponent({ name: 'CEnabledField' });
const selectTabs = wrapper => wrapper.findComponent({ name: 'CFormGeneralPatternsTabs' });
const selectMetaAlarmRuleGeneralForm = wrapper =>
  wrapper.findComponent({ name: 'MetaAlarmRuleGeneralForm' });
const selectMetaAlarmRuleParametersForm = wrapper =>
  wrapper.findComponent({ name: 'MetaAlarmRuleParametersForm' });

describe('meta-alarm-rule-form', () => {
  const form = metaAlarmRuleToForm();

  const factory = generateShallowRenderer(MetaAlarmRuleForm, { stubs });

  test('General form is rendered in general tab', () => {
    const wrapper = factory({
      props: {
        form,
      },
    });

    expect(selectMetaAlarmRuleGeneralForm(wrapper).exists()).toBe(true);
  });

  test('Enabled changed after trigger enabled field', async () => {
    const wrapper = factory({
      props: {
        form,
      },
    });

    selectEnabledField(wrapper).vm.$emit('update:modelValue', false);
    await wrapper.vm.$nextTick();

    expect(globalThis.__TEST_FORM__.getValue('enabled')).toBe(false);
  });

  test('Disabled ID state forwarded to general form', () => {
    const wrapper = factory({
      props: {
        form,
        disabledIdField: true,
      },
    });

    expect(selectMetaAlarmRuleGeneralForm(wrapper).props('disabledIdField')).toBe(true);
  });

  test('Tabs receive form context, rule ID and test type', () => {
    const ruleId = 'rule-id';
    const wrapper = factory({
      props: {
        form,
        ruleId,
      },
    });

    expect(selectTabs(wrapper).props()).toMatchObject({
      form,
      ruleId,
      type: TEMPLATE_TESTING_TEST_TYPES.metaAlarmRule,
    });
    expect(selectMetaAlarmRuleParametersForm(wrapper).exists()).toBe(true);
  });
});
