import Faker from 'faker';

import { generateRenderer } from '@unit/utils/vue';
import { getFormGeneralPatternsTabsStub } from '@unit/stubs/form';

import { TIME_UNITS } from '@/constants';

import AlarmStatusRuleForm from '@/components/other/alarm-status-rule/form/alarm-status-rule-form.vue';

const stubs = {
  'c-enabled-field': true,
  'c-form-general-patterns-tabs': getFormGeneralPatternsTabsStub(),
  'alarm-status-rule-general-form': true,
  'alarm-status-rule-patterns-form': true,
};

const selectGeneralForm = wrapper => wrapper.find('alarm-status-rule-general-form-stub');
const selectAlarmStatusRulePatternsForm = wrapper => wrapper.find('alarm-status-rule-patterns-form-stub');

describe('alarm-status-rule-form', () => {
  const factory = generateRenderer(AlarmStatusRuleForm, { stubs });
  const snapshotFactory = generateRenderer(AlarmStatusRuleForm, { stubs });

  test('General form is rendered in general tab', () => {
    const wrapper = factory({
      propsData: {
        form: {},
      },
    });

    expect(selectGeneralForm(wrapper).exists()).toBe(true);
  });

  test('Name changed after trigger general form', () => {
    const wrapper = factory({
      propsData: {
        form: {},
      },
    });

    const newValue = Faker.datatype.string();

    selectGeneralForm(wrapper).triggerCustomEvent('input', { name: newValue });

    expect(wrapper).toEmitInput({ name: newValue });
  });

  test('Duration changed after trigger general form', () => {
    const wrapper = factory({
      propsData: {
        form: {},
      },
    });

    const newDuration = {
      value: Faker.datatype.number(),
      unit: TIME_UNITS.hour,
    };

    selectGeneralForm(wrapper).triggerCustomEvent('input', { duration: newDuration });

    expect(wrapper).toEmitInput({ duration: newDuration });
  });

  test('Priority changed after trigger general form', () => {
    const wrapper = factory({
      propsData: {
        form: {},
      },
    });

    const newPriority = Faker.datatype.number();

    selectGeneralForm(wrapper).triggerCustomEvent('input', { priority: newPriority });

    expect(wrapper).toEmitInput({ priority: newPriority });
  });

  test('Freq limit changed after trigger general form', () => {
    const wrapper = factory({
      propsData: {
        form: {},
        flapping: true,
      },
    });

    const newFreqLimit = Faker.datatype.number();

    selectGeneralForm(wrapper).triggerCustomEvent('input', { freq_limit: newFreqLimit });

    expect(wrapper).toEmitInput({ freq_limit: newFreqLimit });
  });

  test('Description changed after trigger general form', () => {
    const wrapper = factory({
      propsData: {
        form: {},
      },
    });

    const newDescription = Faker.datatype.string();

    selectGeneralForm(wrapper).triggerCustomEvent('input', { description: newDescription });

    expect(wrapper).toEmitInput({ description: newDescription });
  });

  test('Patterns changed after trigger patterns field', () => {
    const wrapper = factory({
      propsData: {
        form: {},
      },
    });

    const newPatterns = {
      alarm_pattern: {},
    };

    selectAlarmStatusRulePatternsForm(wrapper).triggerCustomEvent('input', newPatterns);

    expect(wrapper).toEmitInput({ patterns: newPatterns });
  });

  test('Renders `alarm-status-rule-form` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `alarm-status-rule-form` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          name: 'name',
          duration: {
            value: 1,
            unit: TIME_UNITS.year,
          },
          priority: 2,
          freq_limit: 3,
          description: 'description',
          patterns: {},
        },
        flapping: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
