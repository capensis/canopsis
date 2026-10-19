import { getFormGeneralPatternsTabsStub } from '@unit/stubs/form';
import { generateRenderer } from '@unit/utils/vue';

import AlarmStatusRuleForm from '@/components/other/alarm-status-rule/form/alarm-status-rule-form.vue';

const stubs = {
  'c-enabled-field': true,
  'c-form-general-patterns-tabs': getFormGeneralPatternsTabsStub(),
  'alarm-status-rule-general-form': true,
  'alarm-status-rule-patterns-form': true,
};

describe('alarm-status-rule-form', () => {
  const factory = generateRenderer(AlarmStatusRuleForm, { stubs });

  test('renders general and patterns form sections', () => {
    const wrapper = factory({ propsData: { form: {} } });

    expect(wrapper.find('alarm-status-rule-general-form-stub').exists()).toBe(true);
    expect(wrapper.find('alarm-status-rule-patterns-form-stub').exists()).toBe(true);
  });

  test('hides patterns for a default rule', () => {
    const wrapper = factory({ propsData: { defaultRule: true, form: {} } });

    expect(wrapper.find('alarm-status-rule-general-form-stub').exists()).toBe(true);
    expect(wrapper.find('alarm-status-rule-patterns-form-stub').exists()).toBe(false);
  });
});
