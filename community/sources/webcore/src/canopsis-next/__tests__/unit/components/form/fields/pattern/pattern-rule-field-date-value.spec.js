import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import {
  PATTERN_OPERATORS,
  PATTERN_QUICK_RANGES_WITHOUT_CUSTOM,
  QUICK_RANGES,
  TIME_UNITS,
} from '@/constants';

import { getDefaultDateFormRange } from '@/helpers/entities/pattern/form';

import PatternRuleFieldDateValue from '@/components/forms/fields/pattern/pattern-rule-field-date-value.vue';

const stubs = {
  'c-date-time-interval-field': {
    name: 'CDateTimeIntervalField',
    props: ['value', 'name', 'disabled'],
    emits: ['update:value'],
    template: '<c-date-time-interval-field-stub />',
  },
  'c-quick-date-interval-type-field': {
    name: 'CQuickDateIntervalTypeField',
    props: ['value', 'name', 'disabled', 'ranges'],
    emits: ['update:value'],
    template: '<c-quick-date-interval-type-field-stub />',
  },
  'c-quick-date-interval-type-range-field': {
    name: 'CQuickDateIntervalTypeRangeField',
    props: ['value', 'name', 'disabled', 'fromRanges', 'toRanges'],
    emits: ['update:value'],
    template: '<c-quick-date-interval-type-range-field-stub />',
  },
};

const selectDateTimeIntervalField = wrapper =>
  wrapper.findComponent({ name: 'CDateTimeIntervalField' });
const selectQuickDateIntervalTypeField = wrapper =>
  wrapper.findComponent({ name: 'CQuickDateIntervalTypeField' });
const selectQuickDateIntervalTypeRangeField = wrapper =>
  wrapper.findComponent({ name: 'CQuickDateIntervalTypeRangeField' });

describe('pattern-rule-field-date-value', () => {
  const defaultValue = {
    type: QUICK_RANGES.last1Hour.value,
    from: 0,
    to: 0,
  };

  const factory = generateShallowRenderer(PatternRuleFieldDateValue, { stubs });

  test('Mount c-date-time-interval-field when operator is inRangeDates', () => {
    const wrapper = factory({
      propsData: {
        value: { from: new Date(), to: new Date() },
        operator: PATTERN_OPERATORS.inRangeDates,
      },
    });

    expect(selectDateTimeIntervalField(wrapper).exists()).toBe(true);
    expect(selectQuickDateIntervalTypeField(wrapper).exists()).toBe(false);
    expect(selectQuickDateIntervalTypeRangeField(wrapper).exists()).toBe(false);
  });

  test('Mount c-quick-date-interval-type-field when operator is within', () => {
    const wrapper = factory({
      propsData: {
        value: defaultValue,
        operator: PATTERN_OPERATORS.within,
      },
    });

    expect(selectDateTimeIntervalField(wrapper).exists()).toBe(false);
    expect(selectQuickDateIntervalTypeField(wrapper).exists()).toBe(true);
    expect(selectQuickDateIntervalTypeRangeField(wrapper).exists()).toBe(false);
  });

  test('Mount c-quick-date-interval-type-field when operator is olderThan', () => {
    const wrapper = factory({
      propsData: {
        value: defaultValue,
        operator: PATTERN_OPERATORS.olderThan,
      },
    });

    expect(selectDateTimeIntervalField(wrapper).exists()).toBe(false);
    expect(selectQuickDateIntervalTypeField(wrapper).exists()).toBe(true);
    expect(selectQuickDateIntervalTypeRangeField(wrapper).exists()).toBe(false);
  });

  test('Mount c-quick-date-interval-type-range-field when operator is inRangePeriod', () => {
    const wrapper = factory({
      propsData: {
        value: defaultValue,
        operator: PATTERN_OPERATORS.inRangePeriod,
      },
    });

    expect(selectDateTimeIntervalField(wrapper).exists()).toBe(false);
    expect(selectQuickDateIntervalTypeField(wrapper).exists()).toBe(false);
    expect(selectQuickDateIntervalTypeRangeField(wrapper).exists()).toBe(true);
  });

  test('Mount nothing when operator is not date-related', () => {
    const wrapper = factory({
      propsData: {
        value: defaultValue,
        operator: PATTERN_OPERATORS.equal,
      },
    });

    expect(selectDateTimeIntervalField(wrapper).exists()).toBe(false);
    expect(selectQuickDateIntervalTypeField(wrapper).exists()).toBe(false);
    expect(selectQuickDateIntervalTypeRangeField(wrapper).exists()).toBe(false);
  });

  test('Passes correct props to c-date-time-interval-field', () => {
    const value = { from: new Date(), to: new Date() };
    const name = Faker.datatype.string();
    const disabled = Faker.datatype.boolean();

    const wrapper = factory({
      propsData: {
        value,
        operator: PATTERN_OPERATORS.inRangeDates,
        name,
        disabled,
      },
    });

    const dateTimeIntervalField = selectDateTimeIntervalField(wrapper);

    expect(dateTimeIntervalField.props('name')).toBe(name);
    expect(dateTimeIntervalField.props('disabled')).toBe(disabled);
    expect(dateTimeIntervalField.props('value')).toEqual(value);
  });

  test('Passes correct props to c-quick-date-interval-type-field', () => {
    const value = { type: QUICK_RANGES.last2Days.value };
    const name = Faker.datatype.string();
    const disabled = Faker.datatype.boolean();
    const intervalRanges = [{ value: QUICK_RANGES.custom.value, text: 'Custom' }];

    const wrapper = factory({
      propsData: {
        value,
        operator: PATTERN_OPERATORS.within,
        name,
        disabled,
        intervalRanges,
      },
    });

    const quickDateIntervalTypeField = selectQuickDateIntervalTypeField(wrapper);

    expect(quickDateIntervalTypeField.props('name')).toBe(name);
    expect(quickDateIntervalTypeField.props('disabled')).toBe(disabled);
    expect(quickDateIntervalTypeField.props('value')).toEqual(value.type);
  });

  test('Passes correct props to c-quick-date-interval-type-range-field', () => {
    const value = { from: QUICK_RANGES.last1Hour.value, to: QUICK_RANGES.last3Hour.value };
    const name = Faker.datatype.string();
    const disabled = Faker.datatype.boolean();
    const intervalRanges = [{ value: QUICK_RANGES.custom.value, text: 'Custom' }];

    const wrapper = factory({
      propsData: {
        value,
        operator: PATTERN_OPERATORS.inRangePeriod,
        name,
        disabled,
        intervalRanges,
      },
    });

    const quickDateIntervalTypeRangeField = selectQuickDateIntervalTypeRangeField(wrapper);

    expect(quickDateIntervalTypeRangeField.props('name')).toBe(name);
    expect(quickDateIntervalTypeRangeField.props('disabled')).toBe(disabled);
    expect(quickDateIntervalTypeRangeField.props('value')).toEqual(value);
  });

  test('Returns true for isDateRange when operator is inRangeDates', () => {
    const wrapper = factory({
      propsData: {
        value: defaultValue,
        operator: PATTERN_OPERATORS.inRangeDates,
      },
    });

    expect(wrapper.vm.isDateRange).toBe(true);
  });

  test('Returns false for isDateRange when operator is not inRangeDates', () => {
    const wrapper = factory({
      propsData: {
        value: defaultValue,
        operator: PATTERN_OPERATORS.within,
      },
    });

    expect(wrapper.vm.isDateRange).toBe(false);
  });

  test('Returns true for isIntervalRange when operator is inRangePeriod', () => {
    const wrapper = factory({
      propsData: {
        value: defaultValue,
        operator: PATTERN_OPERATORS.inRangePeriod,
      },
    });

    expect(wrapper.vm.isIntervalRange).toBe(true);
  });

  test('Returns false for isIntervalRange when operator is not inRangePeriod', () => {
    const wrapper = factory({
      propsData: {
        value: defaultValue,
        operator: PATTERN_OPERATORS.within,
      },
    });

    expect(wrapper.vm.isIntervalRange).toBe(false);
  });

  test('Returns true for isInterval when operator is within', () => {
    const wrapper = factory({
      propsData: {
        value: defaultValue,
        operator: PATTERN_OPERATORS.within,
      },
    });

    expect(wrapper.vm.isInterval).toBe(true);
  });

  test('Returns true for isInterval when operator is olderThan', () => {
    const wrapper = factory({
      propsData: {
        value: defaultValue,
        operator: PATTERN_OPERATORS.olderThan,
      },
    });

    expect(wrapper.vm.isInterval).toBe(true);
  });

  test('Returns false for isInterval when operator is not within or olderThan', () => {
    const wrapper = factory({
      propsData: {
        value: defaultValue,
        operator: PATTERN_OPERATORS.inRangeDates,
      },
    });

    expect(wrapper.vm.isInterval).toBe(false);
  });

  test('Returns true for isAllowedFromDate when no toTimestamp is set', () => {
    const wrapper = factory({
      propsData: {
        value: { from: null, to: null },
        operator: PATTERN_OPERATORS.inRangeDates,
      },
    });

    const testDate = new Date();
    expect(wrapper.vm.isAllowedFromDate(testDate)).toBe(true);
  });

  test('Returns true for isAllowedFromDate when from date is before to date', () => {
    const fromDate = new Date('2023-01-01');
    const toDate = new Date('2023-01-31');

    const wrapper = factory({
      propsData: {
        value: { from: fromDate, to: toDate },
        operator: PATTERN_OPERATORS.inRangeDates,
      },
    });

    const testDate = new Date('2023-01-15');
    expect(wrapper.vm.isAllowedFromDate(testDate)).toBe(true);
  });

  test('Returns false for isAllowedFromDate when from date is after to date', () => {
    const fromDate = new Date('2023-01-01');
    const toDate = new Date('2023-01-15');

    const wrapper = factory({
      propsData: {
        value: { from: fromDate, to: toDate },
        operator: PATTERN_OPERATORS.inRangeDates,
      },
    });

    const testDate = new Date('2023-01-31');
    expect(wrapper.vm.isAllowedFromDate(testDate)).toBe(false);
  });

  test('Returns true for isAllowedToDate when no fromTimestamp is set', () => {
    const wrapper = factory({
      propsData: {
        value: { from: null, to: null },
        operator: PATTERN_OPERATORS.inRangeDates,
      },
    });

    const testDate = new Date();
    expect(wrapper.vm.isAllowedToDate(testDate)).toBe(true);
  });

  test('Returns true for isAllowedToDate when to date is after from date', () => {
    const fromDate = new Date('2023-01-01');
    const toDate = new Date('2023-01-31');

    const wrapper = factory({
      propsData: {
        value: { from: fromDate, to: toDate },
        operator: PATTERN_OPERATORS.inRangeDates,
      },
    });

    const testDate = new Date('2023-01-15');
    expect(wrapper.vm.isAllowedToDate(testDate)).toBe(true);
  });

  test('Returns false for isAllowedToDate when to date is before from date', () => {
    const fromDate = new Date('2023-01-15');
    const toDate = new Date('2023-01-31');

    const wrapper = factory({
      propsData: {
        value: { from: fromDate, to: toDate },
        operator: PATTERN_OPERATORS.inRangeDates,
      },
    });

    const testDate = new Date('2023-01-01');
    expect(wrapper.vm.isAllowedToDate(testDate)).toBe(false);
  });

  test('Updates model with default date form range when operator changes', async () => {
    const wrapper = factory({
      propsData: {
        value: defaultValue,
        operator: PATTERN_OPERATORS.within,
      },
    });

    await wrapper.setProps({ operator: PATTERN_OPERATORS.inRangeDates });

    expect(wrapper).toEmitInput(getDefaultDateFormRange());
  });

  test('Does not update model when operator remains the same', async () => {
    const wrapper = factory({
      propsData: {
        value: defaultValue,
        operator: PATTERN_OPERATORS.within,
      },
    });

    await wrapper.setProps({ operator: PATTERN_OPERATORS.within });

    expect(wrapper.emitted('update:value')).toBeUndefined();
  });

  test('Uses default intervalRanges when not provided', () => {
    const wrapper = factory({
      propsData: {
        value: defaultValue,
        operator: PATTERN_OPERATORS.within,
      },
    });

    expect(wrapper.vm.intervalRanges).toEqual(PATTERN_QUICK_RANGES_WITHOUT_CUSTOM);
  });

  test('Uses default disabled value when not provided', () => {
    const wrapper = factory({
      propsData: {
        value: defaultValue,
        operator: PATTERN_OPERATORS.within,
      },
    });

    expect(wrapper.vm.disabled).toBe(false);
  });

  test('Does not require name prop', () => {
    const wrapper = factory({
      propsData: {
        value: defaultValue,
        operator: PATTERN_OPERATORS.within,
      },
    });

    expect(wrapper.vm.name).toBeUndefined();
  });

  test('Emits input event when c-date-time-interval-field value changes', () => {
    const wrapper = factory({
      propsData: {
        value: { from: new Date(), to: new Date() },
        operator: PATTERN_OPERATORS.inRangeDates,
      },
    });

    const newValue = { from: new Date('2023-01-01'), to: new Date('2023-01-31') };

    selectDateTimeIntervalField(wrapper).triggerCustomEvent('update:value', newValue);

    expect(wrapper).toEmitInput(newValue);
  });

  test('Emits input event when c-quick-date-interval-type-field value changes', () => {
    const value = { type: QUICK_RANGES.last1Hour.value };
    const wrapper = factory({
      propsData: {
        value,
        operator: PATTERN_OPERATORS.within,
      },
    });

    const newType = QUICK_RANGES.last2Days.value;

    selectQuickDateIntervalTypeField(wrapper).triggerCustomEvent('update:value', newType);

    expect(wrapper).toEmitInput({
      ...value,
      type: newType,
    });
  });

  test('Emits input event when c-quick-date-interval-type-range-field value changes', () => {
    const wrapper = factory({
      propsData: {
        value: { from: QUICK_RANGES.last1Hour.value, to: QUICK_RANGES.last3Hour.value },
        operator: PATTERN_OPERATORS.inRangePeriod,
      },
    });

    const newValue = { from: QUICK_RANGES.last2Days.value, to: QUICK_RANGES.last7Days.value };

    selectQuickDateIntervalTypeRangeField(wrapper).triggerCustomEvent('update:value', newValue);

    expect(wrapper).toEmitInput(newValue);
  });

  test('Returns intervalRanges when value.type is not custom', () => {
    const customIntervalRanges = [
      { value: QUICK_RANGES.last15Minutes.value, text: 'Last 5 minutes' },
      { value: QUICK_RANGES.last30Minutes.value, text: 'Last 10 minutes' },
    ];

    const wrapper = factory({
      propsData: {
        value: { type: QUICK_RANGES.last1Hour.value },
        operator: PATTERN_OPERATORS.within,
        intervalRanges: customIntervalRanges,
      },
    });

    expect(wrapper.vm.typeQuickRanges).toEqual(customIntervalRanges);
  });

  test('Returns intervalRanges with custom range item when value.type is custom', () => {
    const customIntervalRanges = [
      { value: QUICK_RANGES.last1Hour.value, text: 'Last 1 hour' },
      { value: QUICK_RANGES.last3Hour.value, text: 'Last 3 hours' },
    ];
    const customDuration = { value: 5, unit: TIME_UNITS.minute };

    const wrapper = factory({
      propsData: {
        value: {
          type: QUICK_RANGES.custom.value,
          typeCustom: customDuration,
        },
        operator: PATTERN_OPERATORS.within,
        intervalRanges: customIntervalRanges,
      },
    });

    const result = wrapper.vm.typeQuickRanges;
    expect(result).toHaveLength(3);
    expect(result.slice(0, 2)).toEqual(customIntervalRanges);
    expect(result[2]).toEqual({
      value: QUICK_RANGES.custom.value,
      text: 'Last 5 minutes',
      start: 'now-5m',
      stop: 'now',
    });
  });

  test('Returns intervalRanges when value.from is not custom', () => {
    const customIntervalRanges = [
      { value: QUICK_RANGES.last15Minutes.value, text: 'Last 5 minutes' },
      { value: QUICK_RANGES.last30Minutes.value, text: 'Last 10 minutes' },
    ];

    const wrapper = factory({
      propsData: {
        value: { from: QUICK_RANGES.last1Hour.value, to: QUICK_RANGES.last3Hour.value },
        operator: PATTERN_OPERATORS.inRangePeriod,
        intervalRanges: customIntervalRanges,
      },
    });

    expect(wrapper.vm.fromQuickRanges).toEqual(customIntervalRanges);
  });

  test('Returns intervalRanges with custom range item when value.from is custom', () => {
    const customIntervalRanges = [
      { value: QUICK_RANGES.last1Hour.value, text: 'Last 1 hour' },
      { value: QUICK_RANGES.last3Hour.value, text: 'Last 3 hours' },
    ];
    const customDuration = { value: 10, unit: TIME_UNITS.hour };

    const wrapper = factory({
      propsData: {
        value: {
          from: QUICK_RANGES.custom.value,
          fromCustom: customDuration,
          to: QUICK_RANGES.last1Hour.value,
        },
        operator: PATTERN_OPERATORS.inRangePeriod,
        intervalRanges: customIntervalRanges,
      },
    });

    const result = wrapper.vm.fromQuickRanges;
    expect(result).toHaveLength(3);
    expect(result.slice(0, 2)).toEqual(customIntervalRanges);
    expect(result[2]).toEqual({
      value: QUICK_RANGES.custom.value,
      text: 'Last 10 hours',
      start: 'now-10h',
      stop: 'now',
    });
  });

  test('Returns intervalRanges when value.to is not custom', () => {
    const customIntervalRanges = [
      { value: QUICK_RANGES.last15Minutes.value, text: 'Last 5 minutes' },
      { value: QUICK_RANGES.last30Minutes.value, text: 'Last 10 minutes' },
    ];

    const wrapper = factory({
      propsData: {
        value: { from: QUICK_RANGES.last1Hour.value, to: QUICK_RANGES.last3Hour.value },
        operator: PATTERN_OPERATORS.inRangePeriod,
        intervalRanges: customIntervalRanges,
      },
    });

    expect(wrapper.vm.toQuickRanges).toEqual(customIntervalRanges);
  });

  test('Returns intervalRanges with custom range item when value.to is custom', () => {
    const customIntervalRanges = [
      { value: QUICK_RANGES.last1Hour.value, text: 'Last 1 hour' },
      { value: QUICK_RANGES.last3Hour.value, text: 'Last 3 hours' },
    ];
    const customDuration = { value: 2, unit: TIME_UNITS.day };

    const wrapper = factory({
      propsData: {
        value: {
          from: QUICK_RANGES.last1Hour.value,
          to: QUICK_RANGES.custom.value,
          toCustom: customDuration,
        },
        operator: PATTERN_OPERATORS.inRangePeriod,
        intervalRanges: customIntervalRanges,
      },
    });

    const result = wrapper.vm.toQuickRanges;
    expect(result).toHaveLength(3);
    expect(result.slice(0, 2)).toEqual(customIntervalRanges);
    expect(result[2]).toEqual({
      value: QUICK_RANGES.custom.value,
      text: 'Last 2 days',
      start: 'now-2d',
      stop: 'now',
    });
  });

  test('Creates correct custom range item for minutes', () => {
    const customDuration = { value: 30, unit: TIME_UNITS.minute };

    const wrapper = factory({
      propsData: {
        value: {
          type: QUICK_RANGES.custom.value,
          typeCustom: customDuration,
        },
        operator: PATTERN_OPERATORS.within,
      },
    });

    const result = wrapper.vm.typeQuickRanges;
    const customItem = result.find(item => item.value === QUICK_RANGES.custom.value);

    expect(customItem).toEqual({
      value: QUICK_RANGES.custom.value,
      text: 'Last 30 minutes',
      start: 'now-30m',
      stop: 'now',
    });
  });

  test('Creates correct custom range item for hours with plural form', () => {
    const customDuration = { value: 5, unit: TIME_UNITS.hour };

    const wrapper = factory({
      propsData: {
        value: {
          type: QUICK_RANGES.custom.value,
          typeCustom: customDuration,
        },
        operator: PATTERN_OPERATORS.within,
      },
    });

    const result = wrapper.vm.typeQuickRanges;
    const customItem = result.find(item => item.value === QUICK_RANGES.custom.value);

    expect(customItem).toEqual({
      value: QUICK_RANGES.custom.value,
      text: 'Last 5 hours',
      start: 'now-5h',
      stop: 'now',
    });
  });

  test('Creates correct custom range item for single day', () => {
    const customDuration = { value: 1, unit: TIME_UNITS.day };

    const wrapper = factory({
      propsData: {
        value: {
          type: QUICK_RANGES.custom.value,
          typeCustom: customDuration,
        },
        operator: PATTERN_OPERATORS.within,
      },
    });

    const result = wrapper.vm.typeQuickRanges;
    const customItem = result.find(item => item.value === QUICK_RANGES.custom.value);

    expect(customItem).toEqual({
      value: QUICK_RANGES.custom.value,
      text: 'Last 1 day',
      start: 'now-1d',
      stop: 'now',
    });
  });

  test('Handles missing duration gracefully', () => {
    const wrapper = factory({
      propsData: {
        value: {
          type: QUICK_RANGES.custom.value,
          typeCustom: { value: 1, unit: TIME_UNITS.day },
        },
        operator: PATTERN_OPERATORS.within,
      },
    });

    const result = wrapper.vm.typeQuickRanges;
    const customItem = result.find(item => item.value === QUICK_RANGES.custom.value);

    expect(customItem).toEqual({
      value: QUICK_RANGES.custom.value,
      text: 'Last 1 day',
      start: 'now-1d',
      stop: 'now',
    });
  });
});
