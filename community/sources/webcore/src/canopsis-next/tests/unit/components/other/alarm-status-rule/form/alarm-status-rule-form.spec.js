import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createInputStub } from '@unit/stubs/input';

import { TIME_UNITS } from '@/constants';

import AlarmStatusRuleForm from '@/components/other/alarm-status-rule/form/alarm-status-rule-form.vue';

const stubs = {
  'c-name-field': true,
  'c-duration-field': true,
  'c-priority-field': true,
  'c-number-field': true,
  'c-description-field': true,
  'alarm-status-rule-patterns-form': true,
  'v-text-field': createInputStub('v-text-field'),
};

const snapshotStubs = {
  'c-name-field': true,
  'c-duration-field': true,
  'c-priority-field': true,
  'c-number-field': true,
  'c-description-field': true,
  'alarm-status-rule-patterns-form': true,
};

const selectNameField = wrapper => wrapper.find('c-name-field-stub');
const selectDurationField = wrapper => wrapper.find('c-duration-field-stub');
const selectPriorityField = wrapper => wrapper.find('c-priority-field-stub');
const selectNumberField = wrapper => wrapper.find('c-number-field-stub');
const selectDescriptionField = wrapper => wrapper.find('c-description-field-stub');
const selectAlarmStatusRulePatternsForm = wrapper => wrapper.find('alarm-status-rule-patterns-form-stub');

describe('alarm-status-rule-form', () => {
  const factory = generateShallowRenderer(AlarmStatusRuleForm, { stubs });
  const snapshotFactory = generateRenderer(AlarmStatusRuleForm, { stubs: snapshotStubs });

  test('Name changed after trigger text field', () => {
    const wrapper = factory({
      propsData: {
        form: {},
      },
    });

    const nameField = selectNameField(wrapper);

    const newValue = Faker.datatype.string();

    nameField.triggerCustomEvent('input', newValue);

    expect(wrapper).toEmitInput({ name: newValue });
  });

  test('Duration changed after trigger duration field', () => {
    const wrapper = factory({
      propsData: {
        form: {},
      },
    });

    const durationField = selectDurationField(wrapper);

    const newDuration = {
      value: Faker.datatype.number(),
      unit: TIME_UNITS.hour,
    };

    durationField.triggerCustomEvent('input', newDuration);

    expect(wrapper).toEmitInput({ duration: newDuration });
  });

  test('Priority changed after trigger priority field', () => {
    const wrapper = factory({
      propsData: {
        form: {},
      },
    });

    const priorityField = selectPriorityField(wrapper);

    const newPriority = Faker.datatype.number();

    priorityField.triggerCustomEvent('input', newPriority);

    expect(wrapper).toEmitInput({ priority: newPriority });
  });

  test('Freq limit changed after trigger number field', () => {
    const wrapper = factory({
      propsData: {
        form: {},
        flapping: true,
      },
    });

    const numberField = selectNumberField(wrapper);

    const newFreqLimit = Faker.datatype.number();

    numberField.triggerCustomEvent('input', newFreqLimit);

    expect(wrapper).toEmitInput({ freq_limit: newFreqLimit });
  });

  test('Description changed after trigger description field', () => {
    const wrapper = factory({
      propsData: {
        form: {},
      },
    });

    const descriptionField = selectDescriptionField(wrapper);

    const newDescription = Faker.datatype.string();

    descriptionField.triggerCustomEvent('input', newDescription);

    expect(wrapper).toEmitInput({ description: newDescription });
  });

  test('Patterns changed after trigger patterns field', () => {
    const wrapper = factory({
      propsData: {
        form: {},
      },
    });

    const alarmStatusRulePatternsForm = selectAlarmStatusRulePatternsForm(wrapper);

    const newPatterns = {
      alarm_pattern: {},
    };

    alarmStatusRulePatternsForm.triggerCustomEvent('input', newPatterns);

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
