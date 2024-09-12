import { flushPromises, generateRenderer } from '@unit/utils/vue';

import { SERVICE_WEATHER_STATE_COUNTERS } from '@/constants';

import AlarmStateCounters from '@/components/widgets/service-weather/alarm-state-counters.vue';

describe('alarm-state-counters', () => {
  const counters = {
    depends: 1,
    all: 2,
    active: 3,
    state: {
      ok: 4,
      minor: 5,
      major: 6,
      critical: 7,
    },
    acked: 8,
    unacked: 9,
    acked_under_pbh: 10,
    under_pbh: 11,
  };
  const snapshotFactory = generateRenderer(AlarmStateCounters, {

    attachTo: document.body,
  });

  test('Renders `alarm-state-counters` with default props', async () => {
    const wrapper = snapshotFactory();

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchTooltipSnapshot();
  });

  test('Renders `alarm-state-counters` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        types: Object.values(SERVICE_WEATHER_STATE_COUNTERS),
        counters,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchTooltipSnapshot();
  });
});
