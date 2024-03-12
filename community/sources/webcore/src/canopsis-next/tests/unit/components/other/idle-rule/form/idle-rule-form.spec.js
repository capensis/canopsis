import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { IDLE_RULE_TYPES } from '@/constants';

import IdleRuleForm from '@/components/other/idle-rule/form/idle-rule-form.vue';

const stubs = {
  'c-enabled-field': true,
  'idle-rule-general-form': true,
  'idle-rule-patterns-form': true,
};

const selectEnabledField = wrapper => wrapper.find('c-enabled-field-stub');
const selectGeneralForm = wrapper => wrapper.find('idle-rule-general-form-stub');
const selectPatternsForm = wrapper => wrapper.find('idle-rule-patterns-form-stub');

describe('idle-rule-form', () => {
  const factory = generateShallowRenderer(IdleRuleForm, { stubs });
  const snapshotFactory = generateRenderer(IdleRuleForm, { stubs });

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

  test('IDLE Rule fields changed after trigger general form', () => {
    const wrapper = factory({
      propsData: {
        form: {},
      },
    });

    const generalForm = selectGeneralForm(wrapper);

    const newValue = {
      name: Faker.datatype.string(),
      description: Faker.datatype.string(),
    };

    generalForm.triggerCustomEvent('input', newValue);

    expect(wrapper).toEmitInput(newValue);
  });

  test('IDLE Rule patterns changed after trigger patterns form', () => {
    const wrapper = factory({
      propsData: {
        form: {
          enabled: true,
          name: 'Name',
          patterns: {},
        },
      },
    });

    const patternsForm = selectPatternsForm(wrapper);

    const newPatterns = {
      alarm_pattern: {},
    };

    patternsForm.triggerCustomEvent('input', newPatterns);

    expect(wrapper).toEmitInput({
      enabled: true,
      name: 'Name',
      patterns: newPatterns,
    });
  });

  test('Renders `idle-rule-form` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `idle-rule-form` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          type: IDLE_RULE_TYPES.entity,
          enabled: true,
          patterns: {},
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `idle-rule-form` with errors', async () => {
    const wrapper = snapshotFactory();

    await wrapper.setData({
      hasGeneralError: true,
      hasPatternsError: true,
    });

    expect(wrapper).toMatchSnapshot();
  });
});
