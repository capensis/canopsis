import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';
import { mockDateNow } from '@unit/utils/mock-hooks';
import { fakeTimestamp } from '@unit/data/date';
import { ALARM_INTERVAL_FIELDS, QUICK_RANGES, TIME_UNITS } from '@/constants';

import DateIntervalSelector from '@/components/forms/date-interval-selector.vue';

const localVue = createVueInstance();

const stubs = {
  'date-time-picker-text-field': true,
  'c-quick-date-interval-type-field': true,
};

const factory = (options = {}) => shallowMount(DateIntervalSelector, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(DateIntervalSelector, {
  localVue,
  stubs,

  ...options,
});

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
    time_field: ALARM_INTERVAL_FIELDS.creationDate,
  };

  test('Start time updated after trigger tstart field', () => {
    const wrapper = factory({
      propsData: {
        value,
      },
    });

    const tstartField = selectTstartField(wrapper);

    tstartField.vm.$emit('input', QUICK_RANGES.last6Hour.start);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual({
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

    const tstopField = selectTstopField(wrapper);

    tstopField.vm.$emit('input', QUICK_RANGES.last6Hour.stop);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual({
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

    const tstartField = selectTstartField(wrapper);

    const timestamp = fakeTimestamp();

    tstartField.vm.$emit('update:objectValue', timestamp);

    const updateStopObjectValueEvents = wrapper.emitted('update:startObjectValue');

    expect(updateStopObjectValueEvents).toHaveLength(1);

    const [eventData] = updateStopObjectValueEvents[0];
    expect(eventData).toBe(timestamp);
  });

  test('Stop object value updated after trigger update object value', () => {
    const wrapper = factory({
      propsData: {
        value,
      },
    });

    const tstopField = selectTstopField(wrapper);

    const timestamp = fakeTimestamp();

    tstopField.vm.$emit('update:objectValue', timestamp);

    const updateStopObjectValueEvents = wrapper.emitted('update:stopObjectValue');

    expect(updateStopObjectValueEvents).toHaveLength(1);

    const [eventData] = updateStopObjectValueEvents[0];
    expect(eventData).toBe(timestamp);
  });

  test('Range updated after trigger quick range field', () => {
    const wrapper = factory({
      propsData: {
        value,
      },
    });

    const quickDateIntervalTypeField = selectQuickDateIntervalTypeField(wrapper);

    quickDateIntervalTypeField.vm.$emit('input', QUICK_RANGES.previousMonth);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual({
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

    const quickDateIntervalTypeField = selectQuickDateIntervalTypeField(wrapper);

    quickDateIntervalTypeField.vm.$emit('input', QUICK_RANGES.custom);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual({
      ...value,
      periodUnit: TIME_UNITS.hour,
      periodValue: 1,
      tstart: '07/12/2013 23:00',
      tstop: '08/12/2013 00:00',
    });
  });

  test('Range didn\'t updated after trigger quick range field with previous value', () => {
    const wrapper = factory({
      propsData: {
        value: QUICK_RANGES.last3Hour,
      },
    });

    const quickDateIntervalTypeField = selectQuickDateIntervalTypeField(wrapper);

    quickDateIntervalTypeField.vm.$emit('input', QUICK_RANGES.last3Hour);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toBeFalsy();
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
    expect(stopDate.getTime()).toBe(1386435660000);
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
    expect(stopDate.getTime()).toBe(1386435600000);
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

    expect(startDate.getTime()).toBe(1385251200000);
    expect(stopDate.getTime()).toBe(1385855940000);
  });

  test('Renders `date-interval-selector` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: {
          tstart: QUICK_RANGES.last3Hour.start,
          tstop: QUICK_RANGES.last3Hour.stop,
          time_field: ALARM_INTERVAL_FIELDS.creationDate,
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `date-interval-selector` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: {
          tstart: QUICK_RANGES.last2Days.start,
          tstop: QUICK_RANGES.last2Days.stop,
          time_field: ALARM_INTERVAL_FIELDS.creationDate,
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

    expect(wrapper.element).toMatchSnapshot();
  });
});
