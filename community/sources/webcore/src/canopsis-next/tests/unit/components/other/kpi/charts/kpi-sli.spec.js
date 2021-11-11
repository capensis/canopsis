import flushPromises from 'flush-promises';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { stubDateNow } from '@unit/utils/stub-hooks';

import { createMockedStoreModule } from '@unit/utils/store';
import { QUICK_RANGES, SAMPLINGS } from '@/constants';

import KpiSli from '@/components/other/kpi/charts/kpi-sli';

const localVue = createVueInstance();

const stubs = {
  'c-quick-date-interval-field': true,
  'kpi-sli-chart': true,
};

const snapshotStubs = {
  'c-quick-date-interval-field': true,
  'kpi-sli-chart': true,
};

const factory = (options = {}) => shallowMount(KpiSli, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(KpiSli, {
  localVue,
  stubs: snapshotStubs,

  ...options,
});

describe('kpi-sli', () => {
  const nowTimestamp = 1386435600000;
  const nowUnix = nowTimestamp / 1000;

  stubDateNow(nowTimestamp);

  it('Metrics fetched after mount', async () => {
    const fetchSliMetrics = jest.fn(() => []);
    const expectedDefaultParams = {
      /* now - 30d  */
      from: 1383843600,
      in_percents: true,
      sampling: SAMPLINGS.day,
      to: nowUnix,
    };

    factory({
      store: createMockedStoreModule('metrics', {
        actions: {
          fetchSliMetricsWithoutStore: fetchSliMetrics,
        },
      }),
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
      /* now - 30d  */
      from: 1386262800,
      in_percents: true,
      sampling: SAMPLINGS.day,
      to: nowUnix,
    };
    const fetchSliMetrics = jest.fn(() => []);

    const wrapper = factory({
      store: createMockedStoreModule('metrics', {
        actions: {
          fetchSliMetricsWithoutStore: fetchSliMetrics,
        },
      }),
    });

    const quickIntervalField = wrapper.find('c-quick-date-interval-field-stub');

    quickIntervalField.vm.$emit('input', {
      from: start,
      to: stop,
    });

    await flushPromises();

    expect(fetchSliMetrics).toBeCalledTimes(2);
    expect(fetchSliMetrics).toBeCalledWith(
      expect.any(Object),
      { params: expectedParamsAfterUpdate },
      undefined,
    );
  });

  it('Renders `kpi-sli` without metrics', async () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModule('metrics', {
        actions: {
          fetchSliMetricsWithoutStore: jest.fn(() => []),
        },
      }),
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });
});
