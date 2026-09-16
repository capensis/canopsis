import { randomArrayItem } from '@unit/utils/array';
import { generateShallowRenderer } from '@unit/utils/vue';

import { AVAILABILITY_SHOW_TYPE, SAMPLINGS } from '@/constants';

import AvailabilityHistory from '@/components/other/availability/partials/availability-history.vue';

const stubs = {
  'availability-history-filters': true,
  'availability-line-chart': true,
};

const selectAvailabilityHistoryFilters = wrapper =>
  wrapper.find('availability-history-filters-stub');
const selectAvailabilityLineChart = wrapper => wrapper.find('availability-line-chart-stub');

describe('availability-history', () => {
  const factory = generateShallowRenderer(AvailabilityHistory, {
    stubs,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });

  test('update:sampling emitted after trigger filters', () => {
    const wrapper = factory({
      propsData: {
        query: {},
      },
    });

    const newSampling = randomArrayItem(Object.values(SAMPLINGS));

    selectAvailabilityHistoryFilters(wrapper).triggerCustomEvent('update:sampling', newSampling);

    expect(wrapper).toEmit('update:sampling', newSampling);
  });

  test('Show type changed after trigger filters', async () => {
    const wrapper = factory({
      propsData: {
        query: {},
      },
    });

    const newShowType = randomArrayItem(Object.values(AVAILABILITY_SHOW_TYPE));

    await selectAvailabilityHistoryFilters(wrapper).triggerCustomEvent(
      'update:show-type',
      newShowType,
    );

    expect(+selectAvailabilityHistoryFilters(wrapper).attributes('showtype')).toEqual(newShowType);
    expect(+selectAvailabilityLineChart(wrapper).attributes('showtype')).toEqual(newShowType);
  });
});
