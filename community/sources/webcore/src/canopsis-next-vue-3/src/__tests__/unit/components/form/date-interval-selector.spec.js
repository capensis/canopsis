import { fakeTimestamp } from '@unit/data/date';
import { mockDateNow } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';

import { ALARM_INTERVAL_FIELDS, QUICK_RANGES } from '@/constants';

import DateIntervalSelector from '@/components/forms/date-interval-selector.vue';

const stubs = {
  'date-time-picker-text-field': true,
  'c-quick-date-interval-type-field': true,
};

const selectDateTimePickerTextFields = wrapper =>
  wrapper.findAll('date-time-picker-text-field-stub');
const selectTstartField = wrapper => selectDateTimePickerTextFields(wrapper).at(0);
const selectTstopField = wrapper => selectDateTimePickerTextFields(wrapper).at(1);
const selectQuickDateIntervalTypeField = wrapper =>
  wrapper.find('c-quick-date-interval-type-field-stub');

describe('date-interval-selector', () => {
  const nowTimestamp = 1386435600000;
  mockDateNow(nowTimestamp);

  const value = {
    tstart: QUICK_RANGES.last3Hour.start,
    tstop: QUICK_RANGES.last3Hour.stop,
    time_field: ALARM_INTERVAL_FIELDS.timestamp,
  };

  const factory = generateShallowRenderer(DateIntervalSelector, { stubs });

  test('Start time updated after trigger tstart field', () => {
    const wrapper = factory({
      propsData: {
        value,
      },
    });

    selectTstartField(wrapper).triggerCustomEvent('input', QUICK_RANGES.last6Hour.start);

    expect(wrapper).toEmitInput({
      ...value,
      tstart: QUICK_RANGES.last6Hour.start,
    });
  });

  test('Stop time updated after trigger tstart field', () => {
    const wrapper = factory({
      propsData: {
        value,
      },
    });

    selectTstopField(wrapper).triggerCustomEvent('input', QUICK_RANGES.last6Hour.stop);

    expect(wrapper).toEmitInput({
      ...value,
      tstop: QUICK_RANGES.last6Hour.stop,
    });
  });

  test('Start object value updated after trigger update object value', () => {
    const wrapper = factory({
      propsData: {
        value,
      },
    });

    const timestamp = fakeTimestamp();
    const date = wrapper.vm.startDateObjectPreparer(timestamp);

    expect(date).toBeInstanceOf(Date);
  });

  test('Stop object value updated after trigger update object value', () => {
    const wrapper = factory({
      propsData: {
        value,
      },
    });

    const timestamp = fakeTimestamp();
    const date = wrapper.vm.stopDateObjectPreparer(timestamp);

    expect(date).toBeInstanceOf(Date);
  });

  test('Range updated after trigger quick range field', () => {
    const wrapper = factory({
      propsData: {
        value,
      },
    });

    selectQuickDateIntervalTypeField(wrapper).triggerCustomEvent(
      'input',
      QUICK_RANGES.previousMonth,
    );

    expect(wrapper).toEmitInput({
      ...value,
      tstart: QUICK_RANGES.previousMonth.start,
      tstop: QUICK_RANGES.previousMonth.stop,
    });
  });

  test('Range updated after trigger quick range field with custom type', () => {
    const wrapper = factory({
      propsData: {
        value,
      },
    });

    wrapper.vm.updateRange(QUICK_RANGES.custom);

    expect(wrapper).not.toHaveBeenEmit('update:value');
  });

  test('Range didn\'t updated after trigger quick range field with previous value', () => {
    const wrapper = factory({
      propsData: {
        value: { tstart: QUICK_RANGES.last3Hour.start, tstop: QUICK_RANGES.last3Hour.stop },
      },
    });

    selectQuickDateIntervalTypeField(wrapper).triggerCustomEvent('input', QUICK_RANGES.last3Hour);

    expect(wrapper).not.toHaveBeenEmit('input');
  });

  test('Dates prepared after trigger prepare callback', () => {
    const wrapper = factory({
      propsData: {
        value,
      },
    });

    const unix = 1386435712;
    const startDate = wrapper.vm.startDateObjectPreparer(unix);
    const stopDate = wrapper.vm.stopDateObjectPreparer(unix);

    expect(startDate.getTime()).toBe(1386435600000);
    expect(stopDate.getTime()).toBe(1386439199999);
  });

  test('Dates prepared after trigger prepare callback with rounded hours', () => {
    const wrapper = factory({
      propsData: {
        value,
        roundHours: true,
      },
    });

    const unix = 1386435712;
    const startDate = wrapper.vm.startDateObjectPreparer(unix);
    const stopDate = wrapper.vm.stopDateObjectPreparer(unix);

    expect(startDate.getTime()).toBe(1386435600000);
    expect(stopDate.getTime()).toBe(1386439199999);
  });

  test('Dates prepared after trigger prepare callback with string value', () => {
    const wrapper = factory({
      propsData: {
        value,
      },
    });

    const startDate = wrapper.vm.startDateObjectPreparer(QUICK_RANGES.previousWeek.start);
    const stopDate = wrapper.vm.stopDateObjectPreparer(QUICK_RANGES.previousWeek.stop);

    expect(startDate.getTime()).toBe(1385316000000);
    expect(stopDate.getTime()).toBe(1385920799999);
  });
});
