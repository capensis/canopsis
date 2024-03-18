import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { AVAILABILITY_SHOW_TYPE, QUICK_RANGES } from '@/constants';

import AvailabilityBar from '@/components/other/availability/partials/availability-bar.vue';

const stubs = {
  'availability-filters': true,
  'availability-bar-chart': true,
};

const selectAvailabilityBarFilters = wrapper => wrapper.find('availability-filters-stub');
const selectAvailabilityBarChart = wrapper => wrapper.find('availability-bar-chart-stub');

describe('availability-bar', () => {
  const factory = generateShallowRenderer(AvailabilityBar, {
    stubs,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });
  const snapshotFactory = generateRenderer(AvailabilityBar, {
    stubs,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });

  test('Update query emitted after trigger filters', () => {
    const wrapper = factory({
      propsData: {
        availability: {
          uptime: Faker.datatype.number(),
          downtime: Faker.datatype.number(),
          inactive_time: Faker.datatype.number(),
        },
        query: {
          interval: {
            from: QUICK_RANGES.today.start,
            to: QUICK_RANGES.today.stop,
          },
        },
      },
    });

    const newInterval = QUICK_RANGES.last3Months;

    selectAvailabilityBarFilters(wrapper).triggerCustomEvent('update:interval', newInterval);

    expect(wrapper).toEmit('update:interval', newInterval);
  });

  test('Show type changed after trigger filters', async () => {
    const wrapper = factory({
      propsData: {
        availability: {
          uptime_duration: Faker.datatype.number(),
          downtime_duration: Faker.datatype.number(),
          inactive_time: Faker.datatype.number(),
        },
        query: {
          interval: {
            from: QUICK_RANGES.today.start,
            to: QUICK_RANGES.today.stop,
          },
        },
      },
    });

    expect(
      +selectAvailabilityBarChart(wrapper).attributes('showtype'),
    ).toEqual(AVAILABILITY_SHOW_TYPE.percent);

    await selectAvailabilityBarFilters(wrapper).triggerCustomEvent('update:show-type', AVAILABILITY_SHOW_TYPE.duration);

    expect(
      +selectAvailabilityBarChart(wrapper).attributes('showtype'),
    ).toEqual(AVAILABILITY_SHOW_TYPE.duration);
  });

  test('Renders `availability-bar` with required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        availability: {
          uptime_duration: 20000,
          uptime_share: 66.67,
          downtime: 10000,
          downtime_share: 33.33,
          inactive_time: 1000,
        },
        query: {
          from: QUICK_RANGES.yesterday.start,
          to: QUICK_RANGES.yesterday.stop,
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `availability-bar` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        availability: {
          uptime_duration: 30000,
          uptime_share: 60,
          downtime_duration: 20000,
          downtime_share: 40,
          inactive_time: 2000,
        },
        query: {
          from: QUICK_RANGES.last6Months.start,
          to: QUICK_RANGES.last6Months.stop,
        },
        minDate: 1709038503904,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
