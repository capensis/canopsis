import flushPromises from 'flush-promises';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { mockDateNow } from '@unit/utils/mock-hooks';

import { createMockedStoreModules } from '@unit/utils/store';
import { ALARM_METRIC_PARAMETERS, QUICK_RANGES, SAMPLINGS } from '@/constants';

import KpiAlarms from '@/components/other/kpi/charts/kpi-alarms';

const stubs = {
  'c-progress-overlay': true,
  'kpi-alarms-filters': true,
  'kpi-alarms-chart': true,
  'kpi-error-overlay': true,
};

describe('kpi-alarms', () => {
  const nowTimestamp = 1386435600000;

  mockDateNow(nowTimestamp);

  const factory = generateShallowRenderer(KpiAlarms, {
    stubs,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });
  const snapshotFactory = generateRenderer(KpiAlarms, {
    stubs,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });

  test('Metrics fetched after mount', async () => {
    const expectedDefaultParams = {
      /* now - 7d  */
      from: 1385852400,
      parameters: [{ metric: ALARM_METRIC_PARAMETERS.createdAlarms }],
      sampling: SAMPLINGS.day,
      filter: null,
      to: 1386370800,
    };
    const fetchAlarmsMetrics = jest.fn(() => ({
      data: [],
      meta: { min_date: expectedDefaultParams.from },
    }));

    factory({
      store: createMockedStoreModules([{
        name: 'metrics',
        actions: {
          fetchAlarmsMetricsWithoutStore: fetchAlarmsMetrics,
        },
      }]),
    });

    expect(fetchAlarmsMetrics).toBeCalledTimes(1);
    expect(fetchAlarmsMetrics).toBeCalledWith(
      expect.any(Object),
      { params: expectedDefaultParams },
      undefined,
    );
  });

  test('Metrics refreshed after change interval', async () => {
    const { start, stop } = QUICK_RANGES.last2Days;
    const expectedParamsAfterUpdate = {
      /* now - 2d  */
      from: 1385852400,
      parameters: [{ metric: ALARM_METRIC_PARAMETERS.createdAlarms }],
      sampling: SAMPLINGS.day,
      filter: null,
      to: 1386370800,
    };
    const fetchAlarmsMetrics = jest.fn(() => ({
      data: [],
      meta: { min_date: expectedParamsAfterUpdate.from },
    }));

    const wrapper = factory({
      store: createMockedStoreModules([{
        name: 'metrics',
        actions: {
          fetchAlarmsMetricsWithoutStore: fetchAlarmsMetrics,
        },
      }]),
    });

    const kpiSliFiltersElement = wrapper.find('kpi-alarms-filters-stub');

    kpiSliFiltersElement.triggerCustomEvent('input', {
      parameters: [ALARM_METRIC_PARAMETERS.createdAlarms],
      sampling: SAMPLINGS.day,
      filter: null,
      interval: {
        from: start,
        to: stop,
      },
    });

    await flushPromises();

    expect(fetchAlarmsMetrics).toBeCalledTimes(2);
    expect(fetchAlarmsMetrics).toBeCalledWith(
      expect.any(Object),
      { params: expectedParamsAfterUpdate },
      undefined,
    );
  });

  test('Metrics doesn\'t refreshed if min date less than from', async () => {
    const fetchAlarmsMetrics = jest.fn(() => ({
      data: [],
      meta: {
        min_date: 1385930800,
      },
    }));

    factory({
      store: createMockedStoreModules([{
        name: 'metrics',
        actions: {
          fetchAlarmsMetricsWithoutStore: fetchAlarmsMetrics,
        },
      }]),
    });

    fetchAlarmsMetrics.mockReset();

    await flushPromises();

    expect(fetchAlarmsMetrics).not.toHaveBeenCalled();
  });

  test('Renders `kpi-alarms` without metrics', async () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([{
        name: 'metrics',
        actions: {
          fetchAlarmsMetricsWithoutStore: jest.fn(() => ({
            data: [],
            meta: { min_date: 1385830800 },
          })),
        },
      }]),
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `kpi-alarms` with fetching error', async () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([{
        name: 'metrics',
        actions: {
          fetchAlarmsMetricsWithoutStore: jest.fn().mockRejectedValue(),
        },
      }]),
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });
});
