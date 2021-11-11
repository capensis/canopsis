import flushPromises from 'flush-promises';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { stubDateNow } from '@unit/utils/stub-hooks';

import { createMockedStoreModule } from '@unit/utils/store';
import { ALARM_METRIC_PARAMETERS, QUICK_RANGES, SAMPLINGS } from '@/constants';

import KpiAlarms from '@/components/other/kpi/charts/kpi-alarms';

const localVue = createVueInstance();

const stubs = {
  'c-quick-date-interval-field': true,
  'kpi-alarms-chart': true,
};

const snapshotStubs = {
  'c-quick-date-interval-field': true,
  'kpi-alarms-chart': true,
};

const factory = (options = {}) => shallowMount(KpiAlarms, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(KpiAlarms, {
  localVue,
  stubs: snapshotStubs,

  ...options,
});

describe('kpi-alarms', () => {
  const nowTimestamp = 1386435600000;
  const nowUnix = nowTimestamp / 1000;

  stubDateNow(nowTimestamp);

  it('Metrics fetched after mount', async () => {
    const fetchAlarmsMetrics = jest.fn(() => []);
    const expectedDefaultParams = {
      /* now - 30d  */
      from: 1383843600,
      parameters: [ALARM_METRIC_PARAMETERS.totalAlarms],
      sampling: SAMPLINGS.day,
      to: nowUnix,
    };

    factory({
      store: createMockedStoreModule('metrics', {
        actions: {
          fetchAlarmsMetricsWithoutStore: fetchAlarmsMetrics,
        },
      }),
    });

    expect(fetchAlarmsMetrics).toBeCalledTimes(1);
    expect(fetchAlarmsMetrics).toBeCalledWith(
      expect.any(Object),
      { params: expectedDefaultParams },
      undefined,
    );
  });

  it('Metrics refreshed after change interval', async () => {
    const { start, stop } = QUICK_RANGES.last2Days;
    const expectedParamsAfterUpdate = {
      /* now - 30d  */
      from: 1386262800,
      parameters: [ALARM_METRIC_PARAMETERS.totalAlarms],
      sampling: SAMPLINGS.day,
      to: nowUnix,
    };
    const fetchAlarmsMetrics = jest.fn(() => []);

    const wrapper = factory({
      store: createMockedStoreModule('metrics', {
        actions: {
          fetchAlarmsMetricsWithoutStore: fetchAlarmsMetrics,
        },
      }),
    });

    const quickIntervalField = wrapper.find('c-quick-date-interval-field-stub');

    quickIntervalField.vm.$emit('input', {
      from: start,
      to: stop,
    });

    await flushPromises();

    expect(fetchAlarmsMetrics).toBeCalledTimes(2);
    expect(fetchAlarmsMetrics).toBeCalledWith(
      expect.any(Object),
      { params: expectedParamsAfterUpdate },
      undefined,
    );
  });

  it('Renders `kpi-alarms` without metrics', async () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModule('metrics', {
        actions: {
          fetchAlarmsMetricsWithoutStore: jest.fn(() => []),
        },
      }),
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });
});
