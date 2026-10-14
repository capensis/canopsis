import { mockDateNow } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';

import { QUICK_RANGES } from '@/constants';

import CQuickDateIntervalField from '@/components/forms/fields/compat/c-quick-date-interval-field.vue';

const stubs = {
  'v-menu': {
    inheritAttrs: false,
    template: `
      <div class="v-menu">
        <slot name="activator" :props="{}" />
        <slot />
      </div>
    `,
  },
  'v-card': {
    inheritAttrs: false,
    template: '<div class="v-card"><slot /></div>',
  },
  'v-card-text': {
    inheritAttrs: false,
    template: '<div class="v-card-text"><slot /></div>',
  },
  'c-date-interval-field': {
    name: 'CDateIntervalField',
    props: [
      'value',
      'disabled',
      'isAllowedFromDate',
      'isAllowedToDate',
      'roundHours',
      'hideDetails',
      'withIcon',
    ],
    emits: ['update:value'],
    template: '<c-date-interval-field-stub />',
  },
  'c-information-block': {
    inheritAttrs: false,
    template: '<div class="c-information-block"><slot /></div>',
  },
  'c-quick-date-interval-type-field': {
    name: 'CQuickDateIntervalTypeField',
    props: ['value', 'ranges', 'disabled', 'hideDetails', 'returnObject'],
    emits: ['input', 'update:value'],
    template: '<c-quick-date-interval-type-field-stub />',
  },
};

const selectDateIntervalField = wrapper => wrapper.findComponent({ name: 'CDateIntervalField' });
const selectQuickDateIntervalTypeField = wrapper =>
  wrapper.findComponent({ name: 'CQuickDateIntervalTypeField' });

describe('c-quick-date-interval-field', () => {
  const nowTimestamp = 1386435500000;
  mockDateNow(nowTimestamp);

  const factory = generateShallowRenderer(CQuickDateIntervalField, { stubs });

  test('Value changed after trigger date interval field', async () => {
    const wrapper = factory();

    const interval = {
      from: 1384435500000,
      to: 1386435500000,
    };

    await selectDateIntervalField(wrapper).triggerCustomEvent('update:value', interval);

    expect(wrapper).toEmitInput(interval);
  });

  test('Value changed after trigger quick date interval type field', async () => {
    const wrapper = factory();

    await selectQuickDateIntervalTypeField(wrapper).triggerCustomEvent(
      'input',
      QUICK_RANGES.last12Hour,
    );

    expect(wrapper).toEmitInput({
      from: QUICK_RANGES.last12Hour.start,
      to: QUICK_RANGES.last12Hour.stop,
    });
  });

  test('Value changed after trigger quick date interval type field with short format', async () => {
    const wrapper = factory({
      propsData: {
        short: true,
      },
    });

    await selectQuickDateIntervalTypeField(wrapper).triggerCustomEvent(
      'input',
      QUICK_RANGES.last12Hour,
    );

    expect(wrapper).toEmitInput({
      from: QUICK_RANGES.last12Hour.start,
      to: QUICK_RANGES.last12Hour.stop,
    });
  });

  test('Value not changed after trigger date interval field with custom quick range', async () => {
    const wrapper = factory();

    await selectQuickDateIntervalTypeField(wrapper).triggerCustomEvent(
      'input',
      QUICK_RANGES.custom,
    );

    expect(wrapper).not.toHaveBeenEmit('input');
    expect(wrapper).not.toHaveBeenEmit('update:interval');
  });

  test('Ranges filtered with accumulatedBefore and min', () => {
    const wrapper = factory({
      propsData: {
        accumulatedBefore: 1386335500,
        min: 1386235500,
      },
    });

    expect(selectQuickDateIntervalTypeField(wrapper).props('ranges')).toEqual([
      QUICK_RANGES.last15Minutes,
      QUICK_RANGES.last30Minutes,
      QUICK_RANGES.last1Hour,
      QUICK_RANGES.last3Hour,
      QUICK_RANGES.last6Hour,
      QUICK_RANGES.last12Hour,
      QUICK_RANGES.last24Hour,
      QUICK_RANGES.today,
      QUICK_RANGES.todaySoFar,
      QUICK_RANGES.custom,
    ]);
  });

  test('Allowed from date function work correctly with accumulated before', () => {
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

    const isAllowedFromDate = selectDateIntervalField(wrapper).props('isAllowedFromDate');

    expect(isAllowedFromDate(accumulatedBefore - 1)).toBeFalsy();
    expect(isAllowedFromDate(accumulatedBefore)).toBeTruthy();
    expect(isAllowedFromDate(accumulatedBefore + 1)).toBeTruthy();
    /** Monday */
    expect(isAllowedFromDate(1384125400)).toBeTruthy();
  });

  test('Allowed from date function work correctly with min', () => {
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

    const isAllowedFromDate = selectDateIntervalField(wrapper).props('isAllowedFromDate');

    expect(isAllowedFromDate(min - 1)).toBeFalsy();
    expect(isAllowedFromDate(min)).toBeTruthy();
    expect(isAllowedFromDate(min + 1)).toBeTruthy();
  });

  test('Allowed from date function work correctly with to date', () => {
    const to = 1386435500;

    const wrapper = factory({
      propsData: {
        interval: {
          to,
          from: 0,
        },
      },
    });

    const isAllowedFromDate = selectDateIntervalField(wrapper).props('isAllowedFromDate');

    expect(isAllowedFromDate(to - 1)).toBeTruthy();
    expect(isAllowedFromDate(to)).toBeFalsy();
    expect(isAllowedFromDate(to + 1)).toBeFalsy();
  });

  test('Allowed to date function work correctly with accumulated before', () => {
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

    const isAllowedToDate = selectDateIntervalField(wrapper).props('isAllowedToDate');

    expect(isAllowedToDate(accumulatedBefore - 1)).toBeFalsy();
    expect(isAllowedToDate(accumulatedBefore)).toBeTruthy();
    expect(isAllowedToDate(accumulatedBefore + 1)).toBeTruthy();
    /** Sunday */
    expect(isAllowedToDate(1385299900)).toBeTruthy();
  });

  test('Allowed to date function work correctly with now', () => {
    const wrapper = factory({
      propsData: {
        interval: {
          to: 1386435500,
          from: 0,
        },
      },
    });

    const isAllowedToDate = selectDateIntervalField(wrapper).props('isAllowedToDate');
    const now = nowTimestamp / 1000;

    expect(isAllowedToDate(now - 1)).toBeTruthy();
    expect(isAllowedToDate(now)).toBeTruthy();
    expect(isAllowedToDate(now + 1)).toBeFalsy();
  });

  test('Allowed to date function work correctly with from date', () => {
    const from = 1386435500;

    const wrapper = factory({
      propsData: {
        interval: {
          to: 1386435500,
          from,
        },
      },
    });

    const isAllowedToDate = selectDateIntervalField(wrapper).props('isAllowedToDate');

    expect(isAllowedToDate(from - 1)).toBeFalsy();
    expect(isAllowedToDate(from)).toBeFalsy();
    expect(isAllowedToDate(from + 1)).toBeFalsy();
  });
});
