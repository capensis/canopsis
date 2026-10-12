import { generateShallowRenderer } from '@unit/utils/vue';

import { QUICK_RANGES } from '@/constants';

import CDateIntervalField from '@/components/forms/fields/date-picker/c-date-interval-field.vue';

const stubs = {
  'c-date-picker-field': {
    name: 'CDatePickerFieldStub',
    props: ['value', 'allowedDates'],
    emits: ['update:value'],
    template: '<input class="c-date-picker-field" />',
  },
};

const selectDateIntervalFields = wrapper =>
  wrapper.findAllComponents({ name: 'CDatePickerFieldStub' });
const selectFromDateIntervalField = wrapper => selectDateIntervalFields(wrapper).at(0);
const selectToDateIntervalField = wrapper => selectDateIntervalFields(wrapper).at(1);

describe('c-date-interval-field', () => {
  const factory = generateShallowRenderer(CDateIntervalField, { stubs });

  test('From value changed after trigger from date picker', () => {
    const wrapper = factory({
      propsData: {
        value: {
          from: QUICK_RANGES.last12Hour.start,
          to: QUICK_RANGES.last12Hour.stop,
        },
      },
    });

    selectFromDateIntervalField(wrapper).triggerCustomEvent(
      'update:value',
      QUICK_RANGES.last3Hour.start,
    );

    expect(wrapper).toEmitInput({
      from: QUICK_RANGES.last3Hour.start,
      to: QUICK_RANGES.last12Hour.stop,
    });
  });

  test('To value changed after trigger from date picker', () => {
    const wrapper = factory({
      propsData: {
        value: {
          from: QUICK_RANGES.last12Hour.start,
          to: QUICK_RANGES.last12Hour.stop,
        },
      },
    });

    selectToDateIntervalField(wrapper).triggerCustomEvent(
      'update:value',
      QUICK_RANGES.last3Hour.stop,
    );

    expect(wrapper).toEmitInput({
      from: QUICK_RANGES.last12Hour.start,
      to: QUICK_RANGES.last3Hour.stop,
    });
  });

  test('Allowed dates on the "from" picker works', () => {
    const isAllowedFromDate = vi.fn();

    const wrapper = factory({
      propsData: {
        value: {
          from: QUICK_RANGES.last12Hour.start,
          to: QUICK_RANGES.last12Hour.stop,
        },
        isAllowedFromDate,
      },
    });

    selectFromDateIntervalField(wrapper).vm.allowedDates();

    expect(isAllowedFromDate).toBeCalled();
  });

  test('Allowed dates on the "to" picker works', () => {
    const isAllowedToDate = vi.fn();

    const wrapper = factory({
      propsData: {
        value: {
          from: QUICK_RANGES.last12Hour.start,
          to: QUICK_RANGES.last12Hour.stop,
        },
        isAllowedToDate,
      },
    });

    selectToDateIntervalField(wrapper).vm.allowedDates();

    expect(isAllowedToDate).toBeCalled();
  });
});
