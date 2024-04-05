import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { mockDateNow } from '@unit/utils/mock-hooks';
import { fakeTimestamp } from '@unit/data/date';

import { ALARM_INTERVAL_FIELDS, QUICK_RANGES, TIME_UNITS } from '@/constants';

import DateIntervalSelector from '@/components/forms/date-interval-selector.vue';

const stubs = {
  'date-time-picker-text-field': true,
  'c-quick-date-interval-type-field': true,
};

const selectDateTimePickerTextFields = wrapper => wrapper.findAll('date-time-picker-text-field-stub');
const selectTstartField = wrapper => selectDateTimePickerTextFields(wrapper).at(0);
const selectTstopField = wrapper => selectDateTimePickerTextFields(wrapper).at(1);
const selectQuickDateIntervalTypeField = wrapper => wrapper.find('c-quick-date-interval-type-field-stub');

describe('date-interval-selector', () => {
  const nowTimestamp = 1386435600000;
  mockDateNow(nowTimestamp);

  const value = {
    tstart: QUICK_RANGES.last3Hour.start,
    tstop: QUICK_RANGES.last3Hour.stop,
    time_field: ALARM_INTERVAL_FIELDS.timestamp,
  };

  const factory = generateShallowRenderer(DateIntervalSelector, { stubs });
  const snapshotFactory = generateRenderer(DateIntervalSelector, { stubs });

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

    selectTstartField(wrapper).triggerCustomEvent('update:objectValue', timestamp);

    expect(wrapper).toEmit('update:startObjectValue', timestamp);
  });

  test('Stop object value updated after trigger update object value', () => {
    const wrapper = factory({
      propsData: {
        value,
      },
    });

    const timestamp = fakeTimestamp();

    selectTstopField(wrapper).triggerCustomEvent('update:objectValue', timestamp);

    expect(wrapper).toEmit('update:stopObjectValue', timestamp);
  });

  test('Range updated after trigger quick range field', () => {
    const wrapper = factory({
      propsData: {
        value,
      },
    });

    selectQuickDateIntervalTypeField(wrapper).triggerCustomEvent('input', QUICK_RANGES.previousMonth);

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

    selectQuickDateIntervalTypeField(wrapper).triggerCustomEvent('input', QUICK_RANGES.custom);

    expect(wrapper).toEmitInput({
      ...value,
      periodUnit: TIME_UNITS.hour,
      periodValue: 1,
      tstart: '07/12/2013 17:00',
      tstop: '07/12/2013 18:00',
    });
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

    const tstartField = selectTstartField(wrapper);
    const tstopField = selectTstopField(wrapper);

    const unix = 1386435712;
    const startDate = tstartField.vm.dateObjectPreparer(unix);
    const stopDate = tstopField.vm.dateObjectPreparer(unix);

    expect(startDate.getTime()).toBe(1386435660000);
    expect(stopDate.getTime()).toBe(1386435719999);
  });

  test('Dates prepared after trigger prepare callback with rounded hours', () => {
    const wrapper = factory({
      propsData: {
        value,
        roundHours: true,
      },
    });

    const tstartField = selectTstartField(wrapper);
    const tstopField = selectTstopField(wrapper);

    const unix = 1386435712;
    const startDate = tstartField.vm.dateObjectPreparer(unix);
    const stopDate = tstopField.vm.dateObjectPreparer(unix);

    expect(startDate.getTime()).toBe(1386435600000);
    expect(stopDate.getTime()).toBe(1386439199999);
  });

  test('Dates prepared after trigger prepare callback with string value', () => {
    const wrapper = factory({
      propsData: {
        value,
      },
    });

    const tstartField = selectTstartField(wrapper);
    const tstopField = selectTstopField(wrapper);

    const startDate = tstartField.vm.dateObjectPreparer(QUICK_RANGES.previousWeek.start);
    const stopDate = tstopField.vm.dateObjectPreparer(QUICK_RANGES.previousWeek.stop);

    expect(startDate.getTime()).toBe(1385334000000);
    expect(stopDate.getTime()).toBe(1385938799999);
  });

  test('Renders `date-interval-selector` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: {
          tstart: QUICK_RANGES.last3Hour.start,
          tstop: QUICK_RANGES.last3Hour.stop,
          time_field: ALARM_INTERVAL_FIELDS.timestamp,
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `date-interval-selector` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: {
          tstart: QUICK_RANGES.last2Days.start,
          tstop: QUICK_RANGES.last2Days.stop,
          time_field: ALARM_INTERVAL_FIELDS.timestamp,
        },
        roundHours: true,
        tstartRules: {
          required: true,
        },
        tstopRules: {
          required: true,
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
