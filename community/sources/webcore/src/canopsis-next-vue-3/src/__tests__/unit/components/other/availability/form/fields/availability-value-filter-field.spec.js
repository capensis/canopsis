import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import { AVAILABILITY_SHOW_TYPE, AVAILABILITY_VALUE_FILTER_METHODS } from '@/constants';

import AvailabilityValueFilterField from '@/components/other/availability/form/fields/availability-value-filter-field.vue';

const stubs = {
  'availability-value-filter-method-field': {
    template: '<button class="method-field" />',
    props: ['method'],
  },
  'c-number-field': {
    template: '<button class="percent-field" />',
    props: ['value'],
  },
  'c-splitted-duration-field': {
    template: '<button class="duration-field" />',
    props: ['value'],
  },
  'c-action-btn': true,
};

const selectAvailabilityValueFilterMethodField = wrapper =>
  wrapper.find('availability-value-filter-method-field-stub');
const selectPercentField = wrapper => wrapper.find('c-number-field-stub');
const selectSplittedDurationField = wrapper => wrapper.find('c-splitted-duration-field-stub');

describe('availability-value-filter-field', () => {
  const factory = generateShallowRenderer(AvailabilityValueFilterField, { stubs });

  test('Value filter changed after trigger availability value filter method field without value', async () => {
    const wrapper = factory();

    wrapper.vm.method = AVAILABILITY_VALUE_FILTER_METHODS.greater;
    expect(wrapper.emitted('update:modelValue')).toEqual([[
      { method: AVAILABILITY_VALUE_FILTER_METHODS.greater, value: 0 },
    ]]);
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

    wrapper.vm.method = AVAILABILITY_VALUE_FILTER_METHODS.less;

    expect(wrapper.emitted('update:modelValue')).toEqual([[{
      ...valueFilter,
      method: AVAILABILITY_VALUE_FILTER_METHODS.less,
    }]]);
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

    wrapper.vm.updateField('value', newValue);

    expect(wrapper.emitted('update:modelValue')).toEqual([[{
      ...valueFilter,
      value: newValue,
    }]]);
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

    wrapper.vm.updateField('value', newValue);

    expect(wrapper.emitted('update:modelValue')).toEqual([[{
      ...valueFilter,
      value: newValue,
    }]]);
  });
});
