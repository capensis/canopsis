import { randomArrayItem } from '@unit/utils/array';
import { generateShallowRenderer } from '@unit/utils/vue';

import { AVAILABILITY_SHOW_TYPE, SAMPLINGS } from '@/constants';

import AvailabilityHistory from '@/components/other/availability/partials/availability-history.vue';

const stubs = {
  'availability-history-filters': true,
  'availability-line-chart': true,
};

describe('availability-history', () => {
  const factory = generateShallowRenderer(AvailabilityHistory, {
    stubs,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });

  test('update:sampling emitted after filters change', () => {
    const wrapper = factory({
      propsData: {
        query: {},
      },
    });

    const newSampling = randomArrayItem(Object.values(SAMPLINGS));

    wrapper.vm.$emit('update:sampling', newSampling);

    expect(wrapper.emitted('update:sampling')).toEqual([[newSampling]]);
  });

  test('Show type changed after trigger filters', async () => {
    const wrapper = factory({
      propsData: {
        query: {},
      },
    });

    const newShowType = randomArrayItem(Object.values(AVAILABILITY_SHOW_TYPE));

    wrapper.vm.showType = newShowType;

    expect(wrapper.vm.showType).toBe(newShowType);
  });
});
