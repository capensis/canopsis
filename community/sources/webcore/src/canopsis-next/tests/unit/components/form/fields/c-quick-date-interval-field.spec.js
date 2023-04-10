import flushPromises from 'flush-promises';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { mockDateNow } from '@unit/utils/mock-hooks';
import { QUICK_RANGES } from '@/constants';

import CQuickDateIntervalField from '@/components/forms/fields/c-quick-date-interval-field.vue';

const stubs = {
  'c-date-interval-field': true,
  'c-information-block': true,
  'c-quick-date-interval-type-field': true,
};

const selectDateIntervalField = wrapper => wrapper.find('c-date-interval-field-stub');
const selectQuickDateIntervalTypeField = wrapper => wrapper.find('c-quick-date-interval-type-field-stub');

describe('c-quick-date-interval-field', () => {
  const nowTimestamp = 1386435500000;
  mockDateNow(nowTimestamp);

  const factory = generateShallowRenderer(CQuickDateIntervalField, { stubs });
  const snapshotFactory = generateRenderer(CQuickDateIntervalField, { stubs });

  it('Value changed after trigger date interval field', () => {
    const wrapper = factory();

    const dateIntervalField = selectDateIntervalField(wrapper);

    const interval = {
      from: 1384435500000,
      to: 1386435500000,
    };

    dateIntervalField.vm.$emit('input', interval);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe(interval);
  });

  it('Value changed after trigger quick date interval type field', () => {
    const wrapper = factory();

    selectQuickDateIntervalTypeField(wrapper).vm.$emit('input', QUICK_RANGES.last12Hour);

    expect(wrapper).toEmit('input', {
      from: QUICK_RANGES.last12Hour.start,
      to: QUICK_RANGES.last12Hour.stop,
    });
  });

  it('Value changed after trigger quick date interval type field with short format', () => {
    const wrapper = factory({
      propsData: {
        short: true,
      },
    });

    selectQuickDateIntervalTypeField(wrapper).vm.$emit('input', QUICK_RANGES.last12Hour);

    expect(wrapper).toEmit('input', {
      from: QUICK_RANGES.last12Hour.start,
      to: QUICK_RANGES.last12Hour.stop,
    });
  });

  it('Value not changed after trigger date interval field with custom quick range', () => {
    const wrapper = factory();

    selectQuickDateIntervalTypeField(wrapper).vm.$emit('input', QUICK_RANGES.custom);

    expect(wrapper).not.toEmit('input');
  });

  it('Ranges filtered with accumulatedBefore and min', () => {
    const wrapper = factory({
      propsData: {
        accumulatedBefore: 1386335500,
        min: 1386235500,
      },
    });

    const quickDateIntervalTypeField = selectQuickDateIntervalTypeField(wrapper);

    expect(quickDateIntervalTypeField.vm.$attrs.ranges).toEqual([
      QUICK_RANGES.last15Minutes,
      QUICK_RANGES.last30Minutes,
      QUICK_RANGES.last1Hour,
      QUICK_RANGES.last3Hour,
      QUICK_RANGES.last6Hour,
      QUICK_RANGES.last12Hour,
      QUICK_RANGES.last24Hour,
      QUICK_RANGES.todaySoFar,
      QUICK_RANGES.custom,
    ]);
  });

  it('Allowed from date function work correctly with accumulated before', () => {
    const accumulatedBefore = 1386335500;

    const wrapper = factory({
      propsData: {
        interval: {
          to: 1386435500,
          from: 0,
        },
        accumulatedBefore,
      },
    });

    const dateIntervalField = selectDateIntervalField(wrapper);

    const isAllowedFromDate = dateIntervalField.vm.$attrs['is-allowed-from-date'];

    expect(isAllowedFromDate(accumulatedBefore - 1)).toBeFalsy();
    expect(isAllowedFromDate(accumulatedBefore)).toBeTruthy();
    expect(isAllowedFromDate(accumulatedBefore + 1)).toBeTruthy();
    /** Monday */
    expect(isAllowedFromDate(1384125400)).toBeTruthy();
  });

  it('Allowed from date function work correctly with min', () => {
    const min = 1386335500;

    const wrapper = factory({
      propsData: {
        interval: {
          to: 1386435500,
          from: 0,
        },
        min,
      },
    });

    const dateIntervalField = selectDateIntervalField(wrapper);

    const isAllowedFromDate = dateIntervalField.vm.$attrs['is-allowed-from-date'];

    expect(isAllowedFromDate(min - 1)).toBeFalsy();
    expect(isAllowedFromDate(min)).toBeTruthy();
    expect(isAllowedFromDate(min + 1)).toBeTruthy();
  });

  it('Allowed from date function work correctly with to date', () => {
    const to = 1386435500;

    const wrapper = factory({
      propsData: {
        interval: {
          to,
          from: 0,
        },
      },
    });

    const dateIntervalField = selectDateIntervalField(wrapper);

    const isAllowedFromDate = dateIntervalField.vm.$attrs['is-allowed-from-date'];

    expect(isAllowedFromDate(to - 1)).toBeTruthy();
    expect(isAllowedFromDate(to)).toBeFalsy();
    expect(isAllowedFromDate(to + 1)).toBeFalsy();
  });

  it('Allowed to date function work correctly with accumulated before', () => {
    const accumulatedBefore = 1386335500;

    const wrapper = factory({
      propsData: {
        interval: {
          to: 1386435500,
          from: 0,
        },
        accumulatedBefore,
      },
    });

    const dateIntervalField = selectDateIntervalField(wrapper);

    const isAllowedToDate = dateIntervalField.vm.$attrs['is-allowed-to-date'];

    expect(isAllowedToDate(accumulatedBefore - 1)).toBeFalsy();
    expect(isAllowedToDate(accumulatedBefore)).toBeTruthy();
    expect(isAllowedToDate(accumulatedBefore + 1)).toBeTruthy();
    /** Sunday */
    expect(isAllowedToDate(1385299900)).toBeTruthy();
  });

  it('Allowed to date function work correctly with now', () => {
    const wrapper = factory({
      propsData: {
        interval: {
          to: 1386435500,
          from: 0,
        },
      },
    });

    const dateIntervalField = selectDateIntervalField(wrapper);

    const isAllowedToDate = dateIntervalField.vm.$attrs['is-allowed-to-date'];

    const now = nowTimestamp / 1000;

    expect(isAllowedToDate(now - 1)).toBeTruthy();
    expect(isAllowedToDate(now)).toBeTruthy();
    expect(isAllowedToDate(now + 1)).toBeFalsy();
  });

  it('Allowed to date function work correctly with from date', () => {
    const from = 1386435500;

    const wrapper = factory({
      propsData: {
        interval: {
          to: 1386435500,
          from,
        },
      },
    });

    const dateIntervalField = selectDateIntervalField(wrapper);

    const isAllowedToDate = dateIntervalField.vm.$attrs['is-allowed-to-date'];

    expect(isAllowedToDate(from - 1)).toBeFalsy();
    expect(isAllowedToDate(from)).toBeFalsy();
    expect(isAllowedToDate(from + 1)).toBeFalsy();
  });

  it('Renders `c-quick-date-interval-field` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-quick-date-interval-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        interval: {
          from: QUICK_RANGES.last2Days.start,
          to: QUICK_RANGES.last2Days.stop,
        },
        accumulatedBefore: 1385435500,
        min: 1384435500,
        disabled: true,
        reverse: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-quick-date-interval-field` with short format', async () => {
    snapshotFactory({
      propsData: {
        interval: {
          from: QUICK_RANGES.last2Days.start,
          to: QUICK_RANGES.last2Days.stop,
        },
        accumulatedBefore: 1385435500,
        min: 1384435500,
        disabled: true,
        reverse: true,
        short: true,
      },
    });

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
  });

  it('Renders `c-quick-date-interval-field` with short format and custom interval', async () => {
    snapshotFactory({
      propsData: {
        interval: {
          from: 1385435500,
          to: 1385435500,
        },
        accumulatedBefore: 1385435500,
        min: 1384435500,
        disabled: true,
        reverse: true,
        short: true,
      },
    });

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
  });
});
