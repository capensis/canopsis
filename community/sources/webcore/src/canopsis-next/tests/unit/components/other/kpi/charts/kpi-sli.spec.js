import flushPromises from 'flush-promises';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { mockDateNow } from '@unit/utils/mock-hooks';

import { createMockedStoreModules } from '@unit/utils/store';
import { QUICK_RANGES, SAMPLINGS } from '@/constants';

import KpiSli from '@/components/other/kpi/charts/kpi-sli';

const localVue = createVueInstance();

const stubs = {
  'c-progress-overlay': true,
  'kpi-sli-filters': true,
  'kpi-sli-chart': true,
  'kpi-error-overlay': true,
};

const factory = (options = {}) => shallowMount(KpiSli, {
  localVue,
  stubs,
  parentComponent: {
    provide: {
      $system: {},
    },
  },

  ...options,
});

const snapshotFactory = (options = {}) => mount(KpiSli, {
  localVue,
  stubs,
  parentComponent: {
    provide: {
      $system: {},
    },
  },

  ...options,
});

describe('kpi-sli', () => {
  const nowTimestamp = 1386435600000;
  const nowUnix = nowTimestamp / 1000;

  mockDateNow(nowTimestamp);

  it('Metrics fetched after mount', async () => {
    const expectedDefaultParams = {
      /* now - 7d  */
      from: 1385830800,
      in_percents: true,
      sampling: SAMPLINGS.day,
      to: nowUnix,
      filter: null,
    };
    const fetchSliMetrics = jest.fn(() => ({
      data: [],
      meta: { min_date: expectedDefaultParams.from },
    }));

    factory({
      store: createMockedStoreModules([{
        name: 'metrics',
        actions: {
          fetchSliMetricsWithoutStore: fetchSliMetrics,
        },
      }]),
    });

    expect(fetchSliMetrics).toBeCalledTimes(1);
    expect(fetchSliMetrics).toBeCalledWith(
      expect.any(Object),
      { params: expectedDefaultParams },
      undefined,
    );
  });

  it('Metrics refreshed after change interval', async () => {
    const { start, stop } = QUICK_RANGES.last2Days;
    const expectedParamsAfterUpdate = {
      /* now - 7d  */
      from: 1385830800,
      in_percents: true,
      sampling: SAMPLINGS.day,
      to: nowUnix,
      filter: null,
    };
    const fetchSliMetrics = jest.fn(() => ({
      data: [],
      meta: { min_date: expectedParamsAfterUpdate.from },
    }));

    const wrapper = factory({
      store: createMockedStoreModules([{
        name: 'metrics',
        actions: {
          fetchSliMetricsWithoutStore: fetchSliMetrics,
        },
      }]),
    });

    const kpiSliFiltersElement = wrapper.find('kpi-sli-filters-stub');

    kpiSliFiltersElement.vm.$emit('input', {
      sampling: SAMPLINGS.day,
      filter: null,
      in_percents: true,
      interval: {
        from: start,
        to: stop,
      },
    });

    await flushPromises();

    expect(fetchSliMetrics).toBeCalledTimes(2);
    expect(fetchSliMetrics).toBeCalledWith(
      expect.any(Object),
      { params: expectedParamsAfterUpdate },
      undefined,
    );
  });

  it('Metrics doesn\'t refreshed if min date less than from', async () => {
    const fetchSliMetrics = jest.fn(() => ({
      data: [],
      meta: { min_date: 1385930800 },
    }));

    factory({
      store: createMockedStoreModules([{
        name: 'metrics',
        actions: {
          fetchSliMetricsWithoutStore: fetchSliMetrics,
        },
      }]),
    });

    fetchSliMetrics.mockReset();

    await flushPromises();

    expect(fetchSliMetrics).not.toHaveBeenCalled();
  });

  it('Renders `kpi-sli` without metrics', async () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([{
        name: 'metrics',
        actions: {
          fetchSliMetricsWithoutStore: jest.fn(() => ({
            data: [],
            meta: { min_date: 1385830800 },
          })),
        },
      }]),
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `kpi-sli` with fetching error', async () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([{
        name: 'metrics',
        actions: {
          fetchSliMetricsWithoutStore: jest.fn().mockRejectedValue(),
        },
      }]),
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });
});
