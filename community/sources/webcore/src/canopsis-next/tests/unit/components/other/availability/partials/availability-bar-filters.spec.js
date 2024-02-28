import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { randomArrayItem } from '@unit/utils/array';

import { AVAILABILITY_SHOW_TYPE, QUICK_RANGES } from '@/constants';

import AvailabilityBarFilters from '@/components/other/availability/partials/availability-bar-filters.vue';

const stubs = {
  'c-quick-date-interval-field': true,
  'c-availability-show-type-field': true,
};

const selectQuickDateIntervalField = wrapper => wrapper.find('c-quick-date-interval-field-stub');
const selectAvailabilityShowTypeField = wrapper => wrapper.find('c-availability-show-type-field-stub');

describe('availability-bar-filters', () => {
  const factory = generateShallowRenderer(AvailabilityBarFilters, { stubs });
  const snapshotFactory = generateRenderer(AvailabilityBarFilters, { stubs });

  test('Interval changed after trigger quick interval field', async () => {
    const wrapper = factory();

    const newValue = randomArrayItem(Object.values(QUICK_RANGES));

    await selectQuickDateIntervalField(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmit('update:interval', newValue);
  });

  test('Show type changed after trigger show type field', async () => {
    const wrapper = factory();

    const newValue = randomArrayItem(Object.values(AVAILABILITY_SHOW_TYPE));

    await selectAvailabilityShowTypeField(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmit('update:showType', newValue);
  });

  test('Renders `availability-bar-filters` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `availability-bar-filters` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        interval: QUICK_RANGES.yesterday,
        showType: AVAILABILITY_SHOW_TYPE.duration,
        minIntervalDate: 1000000000,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
