import Faker from 'faker';

import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';

import { TIME_UNITS } from '@/constants';

import AlarmStatusRuleGeneralForm from '@/components/other/alarm-status-rule/form/alarm-status-rule-general-form.vue';

const stubs = {
  'c-form-block': true,
  'c-form-block-row': true,
  'c-name-field': true,
  'c-priority-field': true,
  'c-duration-field': true,
  'c-number-field': true,
  'c-description-field': true,
};

const selectNameField = wrapper => wrapper.find('c-name-field-stub');
const selectPriorityField = wrapper => wrapper.find('c-priority-field-stub');
const selectDurationField = wrapper => wrapper.find('c-duration-field-stub');
const selectFreqLimitField = wrapper => wrapper.find('c-number-field-stub');
const selectDescriptionField = wrapper => wrapper.find('c-description-field-stub');

describe('alarm-status-rule-general-form', () => {
  const form = {
    name: Faker.datatype.string(),
    duration: {
      value: Faker.datatype.number(),
      unit: TIME_UNITS.hour,
    },
    priority: Faker.datatype.number(),
    freq_limit: Faker.datatype.number(),
    description: Faker.datatype.string(),
  };

  const factory = generateShallowRenderer(AlarmStatusRuleGeneralForm, { stubs });
  const snapshotFactory = generateRenderer(AlarmStatusRuleGeneralForm, { stubs });

  test('Name changed after trigger name field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newName = Faker.datatype.string();

    selectNameField(wrapper).triggerCustomEvent('input', newName);

    expect(wrapper).toEmitInput({ ...form, name: newName });
  });

  test('Duration changed after trigger duration field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newDuration = {
      value: Faker.datatype.number(),
      unit: TIME_UNITS.day,
    };

    selectDurationField(wrapper).triggerCustomEvent('input', newDuration);

    expect(wrapper).toEmitInput({ ...form, duration: newDuration });
  });

  test('Priority changed after trigger priority field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newPriority = Faker.datatype.number();

    selectPriorityField(wrapper).triggerCustomEvent('input', newPriority);

    expect(wrapper).toEmitInput({ ...form, priority: newPriority });
  });

  test('Freq limit changed after trigger freq limit field when flapping is true', () => {
    const wrapper = factory({
      propsData: {
        form,
        flapping: true,
      },
    });

    const newFreqLimit = Faker.datatype.number();

    selectFreqLimitField(wrapper).triggerCustomEvent('input', newFreqLimit);

    expect(wrapper).toEmitInput({ ...form, freq_limit: newFreqLimit });
  });

  test('Description changed after trigger description field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newDescription = Faker.datatype.string();

    selectDescriptionField(wrapper).triggerCustomEvent('input', newDescription);

    expect(wrapper).toEmitInput({ ...form, description: newDescription });
  });

  test('Renders `alarm-status-rule-general-form` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {},
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `alarm-status-rule-general-form` with flapping prop', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          name: 'alarm-status-rule-name',
          duration: {
            value: 1,
            unit: TIME_UNITS.hour,
          },
          priority: 2,
          freq_limit: 3,
          description: 'alarm-status-rule-description',
        },
        flapping: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
