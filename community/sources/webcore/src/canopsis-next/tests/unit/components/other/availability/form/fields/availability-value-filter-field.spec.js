import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { AVAILABILITY_SHOW_TYPE, AVAILABILITY_VALUE_FILTER_METHODS } from '@/constants';

import AvailabilityValueFilterField from '@/components/other/availability/form/fields/availability-value-filter-field.vue';

const stubs = {
  'availability-value-filter-method-field': true,
  'c-percents-field': true,
  'c-splitted-duration-field': true,
  'c-action-btn': true,
};

const selectAvailabilityValueFilterMethodField = wrapper => wrapper.find('availability-value-filter-method-field-stub');
const selectPercentField = wrapper => wrapper.find('c-percents-field-stub');
const selectSplittedDurationField = wrapper => wrapper.find('c-splitted-duration-field-stub');

describe('availability-value-filter-field', () => {
  const factory = generateShallowRenderer(AvailabilityValueFilterField, { stubs });
  const snapshotFactory = generateRenderer(AvailabilityValueFilterField, { stubs });

  test('Value filter changed after trigger availability value filter method field without value', async () => {
    const wrapper = factory();

    await selectAvailabilityValueFilterMethodField(wrapper).triggerCustomEvent('input', AVAILABILITY_VALUE_FILTER_METHODS.greater);

    expect(wrapper).toEmitInput({
      method: AVAILABILITY_VALUE_FILTER_METHODS.greater,
      value: 0,
    });
  });

  test('Value filter changed after trigger availability value filter method field with value', async () => {
    const valueFilter = {
      method: AVAILABILITY_VALUE_FILTER_METHODS.greater,
      value: 10,
    };
    const wrapper = factory({
      propsData: {
        value: valueFilter,
      },
    });

    await selectAvailabilityValueFilterMethodField(wrapper).triggerCustomEvent('input', AVAILABILITY_VALUE_FILTER_METHODS.less);

    expect(wrapper).toEmitInput({
      ...valueFilter,
      method: AVAILABILITY_VALUE_FILTER_METHODS.less,
    });
  });

  test('Value changed after trigger percent field', async () => {
    const valueFilter = {
      method: AVAILABILITY_VALUE_FILTER_METHODS.greater,
      value: 0,
    };
    const wrapper = factory({
      propsData: {
        value: valueFilter,
      },
    });

    const newValue = Faker.datatype.number({ min: 1, max: 100 });

    await selectPercentField(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmitInput({
      ...valueFilter,
      value: newValue,
    });
  });

  test('Value changed after trigger duration field', async () => {
    const valueFilter = {
      method: AVAILABILITY_VALUE_FILTER_METHODS.greater,
      value: 0,
    };
    const wrapper = factory({
      propsData: {
        value: valueFilter,
        showType: AVAILABILITY_SHOW_TYPE.duration,
      },
    });

    const newValue = Faker.datatype.number({ min: 1, max: 10000 });

    await selectSplittedDurationField(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmitInput({
      ...valueFilter,
      value: newValue,
    });
  });

  test('Renders `availability-value-filter-field` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `availability-value-filter-field` with percent show type props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: {
          method: AVAILABILITY_VALUE_FILTER_METHODS.greater,
          value: 50,
        },
        showType: AVAILABILITY_SHOW_TYPE.percent,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `availability-value-filter-field` with duration show type props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: {
          method: AVAILABILITY_VALUE_FILTER_METHODS.greater,
          value: 500,
        },
        showType: AVAILABILITY_SHOW_TYPE.duration,
        maxSeconds: 499,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
