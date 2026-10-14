import { randomArrayItem } from '@unit/utils/array';
import { generateShallowRenderer } from '@unit/utils/vue';

import { AVAILABILITY_SHOW_TYPE, QUICK_RANGES } from '@/constants';

import AvailabilityFilters from '@/components/other/availability/partials/availability-filters.vue';

const stubs = {
  'c-quick-date-interval-field': true,
  'availability-show-type-field': true,
};

describe('availability-filters', () => {
  const factory = generateShallowRenderer(AvailabilityFilters, { stubs });

  test('Interval update is exposed by the component contract', () => {
    const wrapper = factory();

    const newValue = randomArrayItem(Object.values(QUICK_RANGES));

    wrapper.vm.$emit('update:interval', newValue);

    expect(wrapper.emitted('update:interval')).toEqual([[newValue]]);
  });

  test('Show type update is exposed by the component contract', () => {
    const wrapper = factory();

    const newValue = randomArrayItem(Object.values(AVAILABILITY_SHOW_TYPE));

    wrapper.vm.$emit('update:showType', newValue);

    expect(wrapper.emitted('update:showType')).toEqual([[newValue]]);
  });
});
