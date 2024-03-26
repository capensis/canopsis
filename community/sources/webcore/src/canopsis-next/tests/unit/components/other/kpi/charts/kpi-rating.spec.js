import Faker from 'faker';

import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { mockDateNow } from '@unit/utils/mock-hooks';
import { createMockedStoreModules } from '@unit/utils/store';

import { ALARM_METRIC_PARAMETERS, QUICK_RANGES, USER_METRIC_PARAMETERS } from '@/constants';

import KpiRating from '@/components/other/kpi/charts/kpi-rating';

const stubs = {
  'c-progress-overlay': true,
  'kpi-rating-filters': true,
  'kpi-rating-chart': true,
  'kpi-error-overlay': true,
};

describe('kpi-rating', () => {
  const nowTimestamp = 1386435600000;
  const nowUnix = nowTimestamp / 1000;

  mockDateNow(nowTimestamp);

  const factory = generateShallowRenderer(KpiRating, {
    stubs,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });
  const snapshotFactory = generateRenderer(KpiRating, {
    stubs,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });

  it('Metrics doesn\'t fetched after mount without criteria', async () => {
    const fetchRatingMetrics = jest.fn(() => []);

    factory({
      store: createMockedStoreModules([{
        name: 'metrics',
        actions: {
          fetchRatingMetricsWithoutStore: fetchRatingMetrics,
        },
      }]),
    });

    await flushPromises();

    expect(fetchRatingMetrics).toBeCalledTimes(0);
  });

  it('Metrics fetched after after set criteria', async () => {
    const expectedDefaultParams = {
      /* now - 30d  */
      from: 1383843600,
      criteria: 1,
      metric: ALARM_METRIC_PARAMETERS.ticketActiveAlarms,
      limit: 10,
      to: nowUnix,
    };
    const fetchRatingMetrics = jest.fn(() => ({
      data: [],
      meta: {
        min_date: expectedDefaultParams.from,
      },
    }));

    const wrapper = factory({
      store: createMockedStoreModules([{
        name: 'metrics',
        actions: {
          fetchRatingMetricsWithoutStore: fetchRatingMetrics,
        },
      }]),
    });

    const kpiRatingFiltersElement = wrapper.find('kpi-rating-filters-stub');

    kpiRatingFiltersElement.triggerCustomEvent('input', {
      criteria: {
        id: expectedDefaultParams.criteria,
      },
      metric: expectedDefaultParams.metric,
      itemsPerPage: expectedDefaultParams.limit,
      interval: {
        from: expectedDefaultParams.from,
        to: expectedDefaultParams.to,
      },
    });

    await flushPromises();

    expect(fetchRatingMetrics).toBeCalledTimes(1);
    expect(fetchRatingMetrics).toBeCalledWith(
      expect.any(Object),
      {
        params: {
          ...expectedDefaultParams,

          from: 1383778800,
          to: 1386370800,
        },
      },
      undefined,
    );
  });

  it('Metrics doesn\'t refreshed if min date less than from', async () => {
    const expectedDefaultParams = {
      /* now - 30d  */
      from: 1383843600,
      criteria: 1,
      metric: ALARM_METRIC_PARAMETERS.ticketActiveAlarms,
      limit: 10,
      to: nowUnix,
    };
    const fetchRatingMetrics = jest.fn(() => ({
      data: [],
      meta: {
        min_date: 1383943600,
      },
    }));

    const wrapper = factory({
      store: createMockedStoreModules([{
        name: 'metrics',
        actions: {
          fetchRatingMetricsWithoutStore: fetchRatingMetrics,
        },
      }]),
    });

    const kpiRatingFiltersElement = wrapper.find('kpi-rating-filters-stub');

    kpiRatingFiltersElement.triggerCustomEvent('input', {
      criteria: {
        id: expectedDefaultParams.criteria,
      },
      metric: expectedDefaultParams.metric,
      itemsPerPage: expectedDefaultParams.itemsPerPage,
      interval: {
        from: expectedDefaultParams.from,
        to: expectedDefaultParams.to,
      },
    });

    await flushPromises();

    expect(fetchRatingMetrics).toHaveBeenCalled();

    fetchRatingMetrics.mockReset();

    await flushPromises();

    expect(fetchRatingMetrics).not.toHaveBeenCalled();
  });

  it('Metrics refreshed after change interval', async () => {
    const { start, stop } = QUICK_RANGES.last2Days;
    const expectedParamsAfterUpdate = {
      /* today - 1d  */
      from: 1386284400,
      criteria: 1,
      filter: Faker.datatype.string(),
      metric: ALARM_METRIC_PARAMETERS.ticketActiveAlarms,
      limit: 10,
      to: 1386370800,
    };
    const fetchRatingMetrics = jest.fn(() => ({
      data: [],
      meta: {
        min_date: expectedParamsAfterUpdate.from,
      },
    }));

    const wrapper = factory({
      store: createMockedStoreModules([{
        name: 'metrics',
        actions: {
          fetchRatingMetricsWithoutStore: fetchRatingMetrics,
        },
      }]),
    });

    const kpiRatingFiltersElement = wrapper.find('kpi-rating-filters-stub');

    kpiRatingFiltersElement.triggerCustomEvent('input', {
      criteria: {
        id: expectedParamsAfterUpdate.criteria,
      },
      filter: expectedParamsAfterUpdate.filter,
      metric: expectedParamsAfterUpdate.metric,
      limit: expectedParamsAfterUpdate.limit,
      interval: {
        from: start,
        to: stop,
      },
    });

    await flushPromises();

    expect(fetchRatingMetrics).toBeCalledTimes(1);
    expect(fetchRatingMetrics).toBeCalledWith(
      expect.any(Object),
      { params: expectedParamsAfterUpdate },
      undefined,
    );
  });

  it('Metrics refreshed without filter with total active time metric', async () => {
    const { start, stop } = QUICK_RANGES.last2Days;
    const expectedParamsAfterUpdate = {
      /* today - 1d  */
      from: 1386284400,
      criteria: 1,
      metric: USER_METRIC_PARAMETERS.totalUserActivity,
      limit: 10,
      to: 1386370800,
    };
    const fetchRatingMetrics = jest.fn(() => ({
      data: [],
      meta: {
        min_date: expectedParamsAfterUpdate.from,
      },
    }));

    const wrapper = factory({
      store: createMockedStoreModules([{
        name: 'metrics',
        actions: {
          fetchRatingMetricsWithoutStore: fetchRatingMetrics,
        },
      }]),
    });

    const kpiRatingFiltersElement = wrapper.find('kpi-rating-filters-stub');

    kpiRatingFiltersElement.triggerCustomEvent('input', {
      criteria: {
        id: expectedParamsAfterUpdate.criteria,
      },
      metric: expectedParamsAfterUpdate.metric,
      limit: expectedParamsAfterUpdate.limit,
      filter: Faker.datatype.string(),
      interval: {
        from: start,
        to: stop,
      },
    });

    await flushPromises();

    expect(fetchRatingMetrics).toBeCalledTimes(1);
    expect(fetchRatingMetrics).toBeCalledWith(
      expect.any(Object),
      { params: expectedParamsAfterUpdate },
      undefined,
    );
  });

  it('Renders `kpi-rating` without metrics', async () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([{
        name: 'metrics',
        actions: {
          fetchRatingMetricsWithoutStore: jest.fn(() => ({
            data: [],
            meta: {},
          })),
        },
      }]),
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `kpi-rating` with fetching error', async () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([{
        name: 'metrics',
        actions: {
          fetchRatingMetricsWithoutStore: jest.fn().mockRejectedValue(),
        },
      }]),
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });
});
