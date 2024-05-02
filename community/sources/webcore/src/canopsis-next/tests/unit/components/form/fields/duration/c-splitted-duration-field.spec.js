import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createNumberInputStub } from '@unit/stubs/input';

import CSplittedDurationField from '@/components/forms/fields/duration/c-splitted-duration-field.vue';

const stubs = {
  'c-number-field': createNumberInputStub('c-number-field'),
};

const snapshotStubs = {
  'c-number-field': true,
};

const selectNumberFieldByLabel = (wrapper, label) => wrapper.find(`input.c-number-field[label="${label}"]`);
const selectMinutesNumberField = wrapper => selectNumberFieldByLabel(wrapper, 'minutes');
const selectHoursNumberField = wrapper => selectNumberFieldByLabel(wrapper, 'hours');
const selectDaysNumberField = wrapper => selectNumberFieldByLabel(wrapper, 'days');

describe('c-splitted-duration-field', () => {
  const factory = generateShallowRenderer(CSplittedDurationField, { stubs });
  const snapshotFactory = generateRenderer(CSplittedDurationField, {
    stubs: snapshotStubs,
  });

  test('Minutes changed after trigger minutes field with value less then 1 hour', async () => {
    const value = 0;
    const wrapper = factory({
      propsData: {
        value,
      },
    });
    const newValue = 1;

    await selectMinutesNumberField(wrapper).setValue(newValue);

    expect(wrapper).toEmitInput(60);

    await wrapper.setProps({
      value: 60,
    });

    expect(selectMinutesNumberField(wrapper).vm.value).toEqual(newValue);
  });

  test('Minutes changed after trigger minutes field with value greater then 1 hour', async () => {
    const value = 0;
    const wrapper = factory({
      propsData: {
        value,
      },
    });
    const newValue = 61;

    await selectMinutesNumberField(wrapper).setValue(newValue);

    expect(wrapper).toEmitInput(3660);

    await wrapper.setProps({
      value: 3660,
    });

    expect(selectHoursNumberField(wrapper).vm.value).toEqual(1);
    expect(selectMinutesNumberField(wrapper).vm.value).toEqual(1);
  });

  test('Hours changed after trigger minutes field with value less then 1 hour', async () => {
    const value = 0;
    const wrapper = factory({
      propsData: {
        value,
      },
    });
    const newValue = 23;

    await selectHoursNumberField(wrapper).setValue(newValue);

    expect(wrapper).toEmitInput(82800);

    await wrapper.setProps({
      value: 82800,
    });

    expect(selectDaysNumberField(wrapper).vm.value).toEqual(0);
    expect(selectHoursNumberField(wrapper).vm.value).toEqual(newValue);
    expect(selectMinutesNumberField(wrapper).vm.value).toEqual(0);
  });

  test('Hours changed after trigger hours field with value greater then 1 day', async () => {
    const value = 0;
    const wrapper = factory({
      propsData: {
        value,
      },
    });
    const newValue = 25;

    await selectHoursNumberField(wrapper).setValue(newValue);

    expect(wrapper).toEmitInput(90000);

    await wrapper.setProps({
      value: 90000,
    });

    expect(selectDaysNumberField(wrapper).vm.value).toEqual(1);
    expect(selectHoursNumberField(wrapper).vm.value).toEqual(1);
    expect(selectMinutesNumberField(wrapper).vm.value).toEqual(0);
  });

  test('Days changed after trigger days field', async () => {
    const value = 0;
    const wrapper = factory({
      propsData: {
        value,
      },
    });
    const newValue = 13;

    await selectDaysNumberField(wrapper).setValue(newValue);

    expect(wrapper).toEmitInput(1123200);

    await wrapper.setProps({
      value: 1123200,
    });

    expect(selectDaysNumberField(wrapper).vm.value).toEqual(newValue);
    expect(selectHoursNumberField(wrapper).vm.value).toEqual(0);
    expect(selectMinutesNumberField(wrapper).vm.value).toEqual(0);
  });

  test('Value changed to maxValue after trigger with big value', async () => {
    const value = 0;
    const maxValue = 1234567;
    const wrapper = factory({
      propsData: {
        value,
        maxValue,
      },
    });

    await selectDaysNumberField(wrapper).setValue(Infinity);

    expect(wrapper).toEmitInput(maxValue);
  });

  test('Renders `c-splitted-duration-field` with required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 0,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `c-splitted-duration-field` with big value and without maxValue', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 10000,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `c-splitted-duration-field` with big value and with maxValue is one day', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 10000,
        /**
         * One day
         */
        maxValue: 86400,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
