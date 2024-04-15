import { generateRenderer } from '@unit/utils/vue';

import { AVAILABILITY_DISPLAY_PARAMETERS, AVAILABILITY_SHOW_TYPE } from '@/constants';

import AvailabilityListColumnValue from '@/components/other/availability/partials/availability-list-column-value.vue';

describe('availability-list-column-value', () => {
  const snapshotFactory = generateRenderer(AvailabilityListColumnValue);

  test('Renders `availability-list-column-value` with uptime and percent type props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        availability: {
          uptime_share: 72.2,
        },
        showType: AVAILABILITY_SHOW_TYPE.percent,
        displayParameter: AVAILABILITY_DISPLAY_PARAMETERS.uptime,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `availability-list-column-value` with uptime and duration type props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        availability: {
          uptime_duration: 134567,
        },
        showType: AVAILABILITY_SHOW_TYPE.duration,
        displayParameter: AVAILABILITY_DISPLAY_PARAMETERS.uptime,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `availability-list-column-value` with downtime and percent type props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        availability: {
          downtime_share: 27.8,
        },
        showType: AVAILABILITY_SHOW_TYPE.percent,
        displayParameter: AVAILABILITY_DISPLAY_PARAMETERS.downtime,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `availability-list-column-value` with downtime and duration type props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        availability: {
          downtime_duration: 23456,
        },
        showType: AVAILABILITY_SHOW_TYPE.duration,
        displayParameter: AVAILABILITY_DISPLAY_PARAMETERS.downtime,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `availability-list-column-value` with trend and duration type', () => {
    const wrapper = snapshotFactory({
      propsData: {
        availability: {
          downtime_share: 21,
        },
        showType: AVAILABILITY_SHOW_TYPE.percent,
        displayParameter: AVAILABILITY_DISPLAY_PARAMETERS.downtime,
        showTrend: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `availability-list-column-value` with trend up and percent type', () => {
    const wrapper = snapshotFactory({
      propsData: {
        availability: {
          downtime_share: 55,
          downtime_share_history: 54,
        },
        showType: AVAILABILITY_SHOW_TYPE.percent,
        displayParameter: AVAILABILITY_DISPLAY_PARAMETERS.downtime,
        showTrend: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `availability-list-column-value` with trend down and percent type', () => {
    const wrapper = snapshotFactory({
      propsData: {
        availability: {
          downtime_share: 44,
          downtime_share_history: 45,
        },
        showType: AVAILABILITY_SHOW_TYPE.percent,
        displayParameter: AVAILABILITY_DISPLAY_PARAMETERS.downtime,
        showTrend: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `availability-list-column-value` without trend changes and percent type', () => {
    const wrapper = snapshotFactory({
      propsData: {
        availability: {
          downtime_share: 23,
          downtime_share_history: 23,
        },
        showType: AVAILABILITY_SHOW_TYPE.percent,
        displayParameter: AVAILABILITY_DISPLAY_PARAMETERS.downtime,
        showTrend: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
