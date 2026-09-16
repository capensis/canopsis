import { createNumberInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';

import CSplittedDurationField from '@/components/forms/fields/duration/c-splitted-duration-field.vue';

const stubs = {
  'c-number-field': createNumberInputStub('c-number-field'),
};

const selectNumberFields = wrapper => wrapper.findAll('input.c-number-field');
const selectDaysNumberField = wrapper => selectNumberFields(wrapper).at(0);
const selectHoursNumberField = wrapper => selectNumberFields(wrapper).at(1);
const selectMinutesNumberField = wrapper => selectNumberFields(wrapper).at(2);

describe('c-splitted-duration-field', () => {
  const factory = generateShallowRenderer(CSplittedDurationField, { stubs });

  test('Minutes changed after trigger minutes field with value less then 1 hour', async () => {
    const wrapper = factory({
      propsData: {
        modelValue: 0,
      },
    });
    const newValue = 1;

    await selectMinutesNumberField(wrapper).triggerCustomEvent('update:value', newValue);

    expect(wrapper).toEmitInput(60);

    await wrapper.setProps({
      modelValue: 60,
    });

    expect(Number(selectMinutesNumberField(wrapper).element.value)).toBe(newValue);
  });

  test('Minutes changed after trigger minutes field with value greater then 1 hour', async () => {
    const wrapper = factory({
      propsData: {
        modelValue: 0,
      },
    });
    const newValue = 61;

    await selectMinutesNumberField(wrapper).triggerCustomEvent('update:value', newValue);

    expect(wrapper).toEmitInput(3660);

    await wrapper.setProps({
      modelValue: 3660,
    });

    expect(Number(selectHoursNumberField(wrapper).element.value)).toBe(1);
    expect(Number(selectMinutesNumberField(wrapper).element.value)).toBe(1);
  });

  test('Hours changed after trigger minutes field with value less then 1 hour', async () => {
    const wrapper = factory({
      propsData: {
        modelValue: 0,
      },
    });
    const newValue = 23;

    await selectHoursNumberField(wrapper).triggerCustomEvent('update:value', newValue);

    expect(wrapper).toEmitInput(82800);

    await wrapper.setProps({
      modelValue: 82800,
    });

    expect(Number(selectDaysNumberField(wrapper).element.value)).toBe(0);
    expect(Number(selectHoursNumberField(wrapper).element.value)).toBe(newValue);
    expect(Number(selectMinutesNumberField(wrapper).element.value)).toBe(0);
  });

  test('Hours changed after trigger hours field with value greater then 1 day', async () => {
    const wrapper = factory({
      propsData: {
        modelValue: 0,
      },
    });
    const newValue = 25;

    await selectHoursNumberField(wrapper).triggerCustomEvent('update:value', newValue);

    expect(wrapper).toEmitInput(90000);

    await wrapper.setProps({
      modelValue: 90000,
    });

    expect(Number(selectDaysNumberField(wrapper).element.value)).toBe(1);
    expect(Number(selectHoursNumberField(wrapper).element.value)).toBe(1);
    expect(Number(selectMinutesNumberField(wrapper).element.value)).toBe(0);
  });

  test('Days changed after trigger days field', async () => {
    const wrapper = factory({
      propsData: {
        modelValue: 0,
      },
    });
    const newValue = 13;

    await selectDaysNumberField(wrapper).triggerCustomEvent('update:value', newValue);

    expect(wrapper).toEmitInput(1123200);

    await wrapper.setProps({
      modelValue: 1123200,
    });

    expect(Number(selectDaysNumberField(wrapper).element.value)).toBe(newValue);
    expect(Number(selectHoursNumberField(wrapper).element.value)).toBe(0);
    expect(Number(selectMinutesNumberField(wrapper).element.value)).toBe(0);
  });

  test('Value changed to maxValue after trigger with big value', async () => {
    const maxValue = 1234567;
    const wrapper = factory({
      propsData: {
        modelValue: 0,
        maxValue,
      },
    });

    await selectDaysNumberField(wrapper).triggerCustomEvent('update:value', Infinity);

    expect(wrapper).toEmitInput(maxValue);
  });
});
