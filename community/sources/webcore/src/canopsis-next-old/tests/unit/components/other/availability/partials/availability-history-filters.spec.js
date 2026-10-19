import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { randomArrayItem } from '@unit/utils/array';

import { AVAILABILITY_SHOW_TYPE, QUICK_RANGES, SAMPLINGS } from '@/constants';

import AvailabilityHistoryFilters from '@/components/other/availability/partials/availability-history-filters.vue';

const stubs = {
  'c-sampling-field': true,
  'availability-show-type-field': true,
};

const selectSamplingField = wrapper => wrapper.find('c-sampling-field-stub');
const selectAvailabilityShowTypeField = wrapper => wrapper.find('availability-show-type-field-stub');

describe('availability-history-filters', () => {
  const factory = generateShallowRenderer(AvailabilityHistoryFilters, { stubs });
  const snapshotFactory = generateRenderer(AvailabilityHistoryFilters, { stubs });

  test('Sampling changed after trigger sampling field', async () => {
    const wrapper = factory();

    const newValue = randomArrayItem(Object.values(SAMPLINGS));

    await selectSamplingField(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmit('update:sampling', newValue);
  });

  test('Show type changed after trigger show type field', async () => {
    const wrapper = factory();

    const newValue = randomArrayItem(Object.values(AVAILABILITY_SHOW_TYPE));

    await selectAvailabilityShowTypeField(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmit('update:showType', newValue);
  });

  test('Renders `availability-history-filters` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `availability-history-filters` with custom props', () => {
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
