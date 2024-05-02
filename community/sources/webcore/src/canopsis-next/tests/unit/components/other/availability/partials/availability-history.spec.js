import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { randomArrayItem } from '@unit/utils/array';

import { AVAILABILITY_DISPLAY_PARAMETERS, AVAILABILITY_SHOW_TYPE, SAMPLINGS } from '@/constants';

import AvailabilityHistory from '@/components/other/availability/partials/availability-history.vue';

const stubs = {
  'availability-history-filters': true,
  'availability-line-chart': true,
};

const selectAvailabilityHistoryFilters = wrapper => wrapper.find('availability-history-filters-stub');
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
  const snapshotFactory = generateRenderer(AvailabilityHistory, {
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

    await selectAvailabilityHistoryFilters(wrapper).triggerCustomEvent('update:show-type', newShowType);

    expect(
      +selectAvailabilityHistoryFilters(wrapper).attributes('showtype'),
    ).toEqual(newShowType);
    expect(
      +selectAvailabilityLineChart(wrapper).attributes('showtype'),
    ).toEqual(newShowType);
  });

  test('Renders `availability-history` with required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        query: {
          sampling: SAMPLINGS.hour,
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `availability-history` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        availabilities: [
          {
            uptime_duration: 20000,
            downtime_duration: 10000,
          },
          {
            uptime_duration: 30000,
            downtime_duration: 10000,
          },
        ],
        query: {
          sampling: SAMPLINGS.hour,
        },
        minDate: 1709038503904,
        defaultShowType: AVAILABILITY_SHOW_TYPE.duration,
        displayParameter: AVAILABILITY_DISPLAY_PARAMETERS.downtime,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
