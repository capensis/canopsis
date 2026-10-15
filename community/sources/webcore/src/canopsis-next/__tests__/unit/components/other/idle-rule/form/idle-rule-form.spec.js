import { getFormGeneralPatternsTabsStub } from '@unit/stubs/form';
import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import IdleRuleForm from '@/components/other/idle-rule/form/idle-rule-form.vue';

const stubs = {
  'c-enabled-field': true,
  'c-form-general-patterns-tabs': getFormGeneralPatternsTabsStub(),
  'idle-rule-general-form': true,
  'idle-rule-patterns-form': true,
};

const selectEnabledField = wrapper => wrapper.find('c-enabled-field-stub');
const selectGeneralForm = wrapper => wrapper.find('idle-rule-general-form-stub');

describe('idle-rule-form', () => {
  const factory = generateShallowRenderer(IdleRuleForm, { stubs });

  test('General form is rendered in general tab', () => {
    const wrapper = factory({
      propsData: {
        form: {
          enabled: true,
        },
      },
    });

    expect(selectGeneralForm(wrapper).exists()).toBe(true);
  });

  test('IDLE Rule enabled after trigger enabled field', () => {
    const enabled = Faker.datatype.boolean();
    const wrapper = factory({
      propsData: {
        form: {
          enabled,
        },
      },
    });

    const enabledField = selectEnabledField(wrapper);

    const newValue = !enabled;

    enabledField.triggerCustomEvent('input', newValue);

    expect(wrapper).toEmitInput({ enabled: newValue });
  });
});
