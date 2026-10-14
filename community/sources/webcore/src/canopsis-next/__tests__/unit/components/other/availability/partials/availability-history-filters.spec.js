import { randomArrayItem } from '@unit/utils/array';
import { generateShallowRenderer } from '@unit/utils/vue';

import { AVAILABILITY_SHOW_TYPE, SAMPLINGS } from '@/constants';

import AvailabilityHistoryFilters from '@/components/other/availability/partials/availability-history-filters.vue';

const stubs = {
  'c-sampling-field': true,
  'availability-show-type-field': true,
};

const selectSamplingField = wrapper => wrapper.find('c-sampling-field-stub');
const selectAvailabilityShowTypeField = wrapper =>
  wrapper.find('availability-show-type-field-stub');

describe('availability-history-filters', () => {
  const factory = generateShallowRenderer(AvailabilityHistoryFilters, { stubs });

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
});
