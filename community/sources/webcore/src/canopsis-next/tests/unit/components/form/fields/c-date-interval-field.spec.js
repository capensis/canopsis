import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { QUICK_RANGES } from '@/constants';

import CDateIntervalField from '@/components/forms/fields/date-picker/c-date-interval-field.vue';

const stubs = {
  'c-date-picker-field': true,
};

const selectDateIntervalFields = wrapper => wrapper.findAll('c-date-picker-field-stub');
const selectFromDateIntervalField = wrapper => selectDateIntervalFields(wrapper).at(0);
const selectToDateIntervalField = wrapper => selectDateIntervalFields(wrapper).at(1);

describe('c-date-interval-field', () => {
  const factory = generateShallowRenderer(CDateIntervalField, { stubs });
  const snapshotFactory = generateRenderer(CDateIntervalField, { stubs });

  test('From value changed after trigger from date picker', () => {
    const wrapper = factory({
      propsData: {
        value: {
          from: QUICK_RANGES.last12Hour.start,
          to: QUICK_RANGES.last12Hour.stop,
        },
      },
    });

    const fromDatePicker = selectFromDateIntervalField(wrapper);

    fromDatePicker.triggerCustomEvent('input', QUICK_RANGES.last3Hour.start);

    expect(wrapper).toEmit('input', {
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

    const toDatePicker = selectToDateIntervalField(wrapper);

    toDatePicker.triggerCustomEvent('input', QUICK_RANGES.last3Hour.stop);

    expect(wrapper).toEmit('input', {
      from: QUICK_RANGES.last12Hour.start,
      to: QUICK_RANGES.last3Hour.stop,
    });
  });

  test('Allowed dates on the "from" picker works', () => {
    const isAllowedFromDate = jest.fn();

    const wrapper = factory({
      propsData: {
        value: {
          from: QUICK_RANGES.last12Hour.start,
          to: QUICK_RANGES.last12Hour.stop,
        },
        isAllowedFromDate,
      },
    });

    const fromDatePicker = selectFromDateIntervalField(wrapper);

    const allowedDates = fromDatePicker.vm.$attrs['allowed-dates'];

    allowedDates();

    expect(isAllowedFromDate).toBeCalled();
  });

  test('Allowed dates on the "to" picker works', () => {
    const isAllowedToDate = jest.fn();

    const wrapper = factory({
      propsData: {
        value: {
          from: QUICK_RANGES.last12Hour.start,
          to: QUICK_RANGES.last12Hour.stop,
        },
        isAllowedToDate,
      },
    });

    const toDatePicker = selectToDateIntervalField(wrapper);

    const allowedDates = toDatePicker.vm.$attrs['allowed-dates'];

    allowedDates();

    expect(isAllowedToDate).toBeCalled();
  });

  test('Renders `c-date-interval-field` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: {
          from: QUICK_RANGES.last12Hour.start,
          to: QUICK_RANGES.last12Hour.stop,
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `c-date-interval-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: {
          from: QUICK_RANGES.last2Days.start,
          to: QUICK_RANGES.last2Days.stop,
        },
        isAllowedFromDate: () => {},
        isAllowedToDate: () => {},
        disabled: true,
        column: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
